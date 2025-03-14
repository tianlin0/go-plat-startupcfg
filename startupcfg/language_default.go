package startupcfg

import (
	"github.com/nicksnyder/go-i18n/v2/i18n/template"
	"github.com/tianlin0/go-plat-utils/templates"
	"sync"
)

var (
	langLock          sync.Mutex
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
func NewI18nFile(yamlFile string, tag string) (*i18nTranslator, error) {
	langLock.Lock()
	defer langLock.Unlock()
	if defaultTranslator == nil {
		defaultTranslator = new(i18nTranslator)
	}
	err := defaultTranslator.InitFile(yamlFile, tag)
	if err != nil {
		return defaultTranslator, err
	}
	//初始化默认解析器
	//defaultTranslator.SetTemplateParser(new(defaultParser))
	return defaultTranslator, nil
}

// DefaultI18nTranslator 默认翻译器
func DefaultI18nTranslator(key string, templateData any) *i18nTranslator {
	langLock.Lock()
	defer langLock.Unlock()
	if defaultTranslator == nil {
		defaultTranslator = new(i18nTranslator)
	}
	return defaultTranslator
}

// I18nTranslate 翻译
func I18nTranslate(key string, templateData ...any) string {
	if templateData == nil || len(templateData) == 0 {
		return defaultTranslator.Translate(key, nil)
	}
	return defaultTranslator.Translate(key, templateData[0])
}
