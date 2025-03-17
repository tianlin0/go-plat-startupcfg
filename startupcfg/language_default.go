package startupcfg

import (
	"github.com/nicksnyder/go-i18n/v2/i18n/template"
	"github.com/tianlin0/go-plat-utils/templates"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"sync"
)

var (
	langLock          sync.Mutex
	i18nCtxKeyName    = "ctx_i18n_key_name"
	defaultTranslator *i18nTranslator
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

// NewI18nFile 初始化
func NewI18nFile(yamlFile string, defaultTagString string) (*i18nTranslator, error) {
	langLock.Lock()
	defer langLock.Unlock()
	if defaultTranslator == nil {
		defaultTranslator = new(i18nTranslator)
	}
	err := defaultTranslator.InitFile(yamlFile, defaultTagString)
	if err != nil {
		return defaultTranslator, err
	}
	//初始化默认解析器
	//defaultTranslator.SetTemplateParser(new(defaultParser))
	return defaultTranslator, nil
}

// NewI18nMap 初始化
func NewI18nMap(langData map[language.Tag]map[string]string, defaultTagString string) (*i18nTranslator, error) {
	langLock.Lock()
	defer langLock.Unlock()
	if defaultTranslator == nil {
		defaultTranslator = new(i18nTranslator)
	}
	err := defaultTranslator.InitMap(langData, defaultTagString)
	if err != nil {
		return defaultTranslator, err
	}
	//初始化默认解析器
	//defaultTranslator.SetTemplateParser(new(defaultParser))
	return defaultTranslator, nil
}

// DefaultI18nTranslator 默认翻译器
func DefaultI18nTranslator() *i18nTranslator {
	langLock.Lock()
	defer langLock.Unlock()
	if defaultTranslator == nil {
		defaultTranslator = new(i18nTranslator)
	}
	return defaultTranslator
}

// I18nTranslate 翻译
func I18nTranslate(tag string, key string, templateData ...any) string {
	if templateData == nil || len(templateData) == 0 {
		return defaultTranslator.TranslateByTag(tag, key, nil)
	}
	return defaultTranslator.TranslateByTag(tag, key, templateData[0])
}

// I18nCtxWithTag 设置当前的语言
func I18nCtxWithTag(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, i18nCtxKeyName, tag)
}

// I18nTranslateCtx Ctx翻译
func I18nTranslateCtx(ctx context.Context, key string, templateData ...any) string {
	tag := ctx.Value(i18nCtxKeyName)
	tagStr := ""
	if tagTemp, ok := tag.(string); ok {
		tagStr = tagTemp
	}
	if templateData == nil || len(templateData) == 0 {
		return defaultTranslator.TranslateByTag(tagStr, key, nil)
	}
	return defaultTranslator.TranslateByTag(tagStr, key, templateData[0])
}
