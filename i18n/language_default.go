package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n/template"
	"github.com/tianlin0/go-plat-utils/templates"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"sync"
)

var (
	langLock          sync.Mutex
	i18nCtxKeyName    = "i18n_ctx_key_name"
	defaultTranslator *translator
)

type defaultParser struct {
	leftDelim, rightDelim string
	format                string
}

func (d *defaultParser) Parse(src, leftDelim, rightDelim string) (template.ParsedTemplate, error) {
	d.leftDelim = leftDelim
	d.rightDelim = rightDelim
	d.format = src
	return d, nil
}
func (d *defaultParser) Cacheable() bool {
	return false
}

func (d *defaultParser) Execute(data any) (string, error) {
	return templates.Template(d.format, data, d.leftDelim, d.rightDelim)
}

// NewYamlFile 初始化
func NewYamlFile(yamlFilePath string, defaultTagString string) (*translator, error) {
	langLock.Lock()
	defer langLock.Unlock()
	if defaultTranslator == nil {
		defaultTranslator = new(translator)
	}
	err := defaultTranslator.InitFile(yamlFilePath, defaultTagString)
	if err != nil {
		return defaultTranslator, err
	}
	//初始化默认解析器
	//defaultTranslator.SetTemplateParser(new(defaultParser))
	return defaultTranslator, nil
}

// NewYamlFileWithLang 初始化
func NewYamlFileWithLang(lang string, yamlFilePath string, defaultTagString string) (*translator, error) {
	langLock.Lock()
	defer langLock.Unlock()
	if defaultTranslator == nil {
		defaultTranslator = new(translator)
	}
	err := defaultTranslator.InitFileWithTag(lang, yamlFilePath, defaultTagString)
	if err != nil {
		return defaultTranslator, err
	}
	//初始化默认解析器
	//defaultTranslator.SetTemplateParser(new(defaultParser))
	return defaultTranslator, nil
}

// NewI18nMap 初始化
// langData  map[string] string 语言简称，map[string]string key,value
func NewI18nMap(langData map[string]map[string]string, defaultTagString string) (*translator, error) {
	langLock.Lock()
	defer langLock.Unlock()
	if defaultTranslator == nil {
		defaultTranslator = new(translator)
	}

	confMap := make(map[language.Tag]map[string]string)
	for k, v := range langData {
		confMap[language.Make(k)] = v
	}

	err := defaultTranslator.InitMap(confMap, defaultTagString)
	if err != nil {
		return defaultTranslator, err
	}
	//初始化默认解析器
	//defaultTranslator.SetTemplateParser(new(defaultParser))
	return defaultTranslator, nil
}

// DefaultTranslator 默认翻译器
func DefaultTranslator() *translator {
	langLock.Lock()
	defer langLock.Unlock()
	if defaultTranslator == nil {
		defaultTranslator = new(translator)
	}
	return defaultTranslator
}

// TranslateLang 翻译
func TranslateLang(lang string, key string, templateData ...any) string {
	if templateData == nil || len(templateData) == 0 {
		return defaultTranslator.TranslateByTag(lang, key, nil)
	}
	return defaultTranslator.TranslateByTag(lang, key, templateData[0])
}

// Translate 翻译
func Translate(key string, templateData ...any) string {
	return TranslateLang("", key, templateData...)
}

// CtxWithLang 设置当前的语言
func CtxWithLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, i18nCtxKeyName, lang)
}

// TranslateCtx Ctx翻译
func TranslateCtx(ctx context.Context, key string, templateData ...any) string {
	tag := ctx.Value(i18nCtxKeyName)
	tagStr := ""
	if tag != nil {
		if tagTemp, ok := tag.(string); ok {
			tagStr = tagTemp
		}
	}
	return TranslateLang(tagStr, key, templateData...)
}
