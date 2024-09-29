package main

type ITemplate interface {
	isTranslated() bool
	GetTemplateCode() string
	GetTemplateText() string
}
