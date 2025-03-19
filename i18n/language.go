package i18n

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n/template"
	"github.com/samber/lo"
	"github.com/tianlin0/go-plat-utils/conv"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v3"
	"os"
)

type translator struct {
	defaultTag            language.Tag
	bundle                *i18n.Bundle
	localize              *i18n.Localizer
	templateParser        template.Parser
	leftDelim, rightDelim string //变量分隔符
}

var (
	allTranslator = make(map[string]*i18n.Localizer) //缓存所有语言,避免重复创建
)

func (l *translator) InitFile(yamlFile string, defaultTagString string) error {
	configFile, err := os.ReadFile(yamlFile)
	if err != nil {
		return err
	}
	conf := make(map[string]map[string]string)
	if err = yaml.Unmarshal(configFile, conf); err != nil {
		return err
	}
	confMap := make(map[language.Tag]map[string]string)
	for k, v := range conf {
		confMap[language.Make(k)] = v
	}
	return l.InitMap(confMap, defaultTagString)
}
func (l *translator) InitMap(conf map[language.Tag]map[string]string, defaultTagString string) error {
	//需要将所有tag的key都进行赋值，避免有遗漏掉的key未配置的情况，而造成输出报错了
	conf = l.builderAllKeys(conf)
	defaultTag := l.initDefaultTag(defaultTagString, conf)
	bundle := i18n.NewBundle(defaultTag)
	l.bundle = bundle
	l.defaultTag = defaultTag
	l.localize = i18n.NewLocalizer(l.bundle, l.defaultTag.String())
	l.templateParser = nil
	l.leftDelim = "{{"
	l.rightDelim = "}}"

	// 添加绑定初始化数据
	for tag, messages := range conf {
		err := l.AddMessage(tag.String(), messages)
		if err != nil {
			return err
		}
	}
	return nil
}

// builderAllKeys 所有的key都需要赋值
func (l *translator) builderAllKeys(allConf map[language.Tag]map[string]string) map[language.Tag]map[string]string {
	//TODO
	return allConf
}

func (l *translator) initDefaultTag(tagString string, allConf map[language.Tag]map[string]string) language.Tag {
	if tagString == "" {
		if len(allConf) > 0 {
			return lo.Keys(allConf)[0] //未设置会取第一个
		}
		return language.Chinese
	}
	return language.Make(tagString)
}

func (l *translator) AddMessage(tag string, msgMap map[string]string) error {
	messageList := make([]*i18n.Message, 0, len(msgMap))
	for k, v := range msgMap {
		messageList = append(messageList, &i18n.Message{
			ID:    k,
			Other: v,
		})
	}
	err := l.bundle.AddMessages(language.Make(tag), messageList...)
	if err != nil {
		return err
	}
	return nil
}
func (l *translator) SetTemplateParser(templateParser template.Parser) {
	l.templateParser = templateParser
}
func (l *translator) SetVariableDelim(leftDelim, rightDelim string) {
	if leftDelim != "" {
		l.leftDelim = leftDelim
	}
	if rightDelim != "" {
		l.rightDelim = rightDelim
	}
}

func (l *translator) DefaultTag() language.Tag {
	return l.defaultTag
}
func (l *translator) Translate(key string, templateData any) string {
	return l.TranslateByTag("", key, templateData)
}
func (l *translator) TranslateByTag(tag string, key string, templateData any) string {
	localize := l.localize
	if tag != "" {
		if localizeTemp, ok := allTranslator[tag]; ok {
			localize = localizeTemp
		} else {
			allTags := l.bundle.LanguageTags() //所有支持的tag，避免乱传不支持
			for _, oneTag := range allTags {
				if oneTag.String() == tag {
					localize = i18n.NewLocalizer(l.bundle, tag)
					allTranslator[tag] = localize
					break
				}
			}
		}
	}

	var retStr string
	var err error
	if templateData == nil {
		retStr, err = localize.Localize(&i18n.LocalizeConfig{
			MessageID:      key,
			TemplateParser: l.templateParser,
		})
	} else {
		retStr, err = localize.Localize(&i18n.LocalizeConfig{
			MessageID:      key,
			TemplateData:   templateData,
			TemplateParser: l.templateParser,
		})
	}
	if err == nil {
		return retStr
	}

	if templateData == nil {
		fmt.Println(err)
		return key
	}

	// 尝试转换
	newTemplateData := make(map[string]any)
	err1 := conv.Unmarshal(templateData, &newTemplateData)
	if err1 != nil {
		return key
	}
	retStr, err = localize.Localize(&i18n.LocalizeConfig{
		MessageID:      key,
		TemplateData:   newTemplateData,
		TemplateParser: l.templateParser,
	})

	if err == nil {
		return retStr
	}

	//有数据的情况，可能需要合并数据
	tempStr, err := l.parserData(key, newTemplateData, l.templateParser)
	if err != nil {
		fmt.Println(err)
		return key
	}
	return tempStr
}

func (l *translator) parserData(src string, data any, parser template.Parser) (string, error) {
	if parser == nil {
		parser = new(template.TextParser)
	}
	temp, err := parser.Parse(src, l.leftDelim, l.rightDelim)
	if err != nil {
		return src, err
	}
	return temp.Execute(data)
}
