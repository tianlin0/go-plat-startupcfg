package startupcfg

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

type i18nTranslator struct {
	defaultTag     language.Tag
	bundle         *i18n.Bundle
	localize       *i18n.Localizer
	templateParser template.Parser
}

func (l *i18nTranslator) InitFile(yamlFile string, defaultTagString string) error {
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
func (l *i18nTranslator) InitMap(conf map[language.Tag]map[string]string, defaultTagString string) error {
	defaultTag := l.initTag(defaultTagString, conf)
	bundle := i18n.NewBundle(defaultTag)
	l.bundle = bundle
	l.defaultTag = defaultTag
	l.localize = i18n.NewLocalizer(l.bundle, l.defaultTag.String())
	l.templateParser = nil

	// 添加绑定初始化数据
	for tag, messages := range conf {
		err := l.AddMessage(tag.String(), messages)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *i18nTranslator) initTag(tagString string, allConf map[language.Tag]map[string]string) language.Tag {
	if tagString == "" {
		if len(allConf) > 0 {
			return lo.Keys(allConf)[0] //未设置会取第一个
		}
		return language.Chinese
	}
	return language.Make(tagString)
}

func (l *i18nTranslator) AddMessage(tag string, msgMap map[string]string) error {
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
func (l *i18nTranslator) SetTemplateParser(templateParser template.Parser) {
	l.templateParser = templateParser
}

func (l *i18nTranslator) Tag() language.Tag {
	return l.defaultTag
}
func (l *i18nTranslator) Translate(key string, templateData any) string {
	return l.TranslateByTag("", key, templateData)
}
func (l *i18nTranslator) TranslateByTag(tag string, key string, templateData any) string {
	localize := l.localize
	if tag != "" {
		localize = i18n.NewLocalizer(l.bundle, tag)
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

	fmt.Println(err)
	return key
}
