package v2

import (
	"embed"
	"path"
	"text/template"

	templateutil "github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

const (
	templatesDir = "templates"

	importsTemplatePath    = templatesDir + "/imports.tmpl"
	typesTemplatePath      = templatesDir + "/types.tmpl"
	schemaTemplatePath     = templatesDir + "/schema.tmpl"
	messageTemplatePath    = templatesDir + "/message.tmpl"
	subscriberTemplatePath = templatesDir + "/subscriber.tmpl"
	controllerTemplatePath = templatesDir + "/controller.tmpl"
	parameterTemplatePath  = templatesDir + "/parameter.tmpl"
)

var (
	//go:embed templates/*
	files embed.FS
)

func loadTemplate(paths ...string) (*template.Template, error) {
	return template.
		New(path.Base(paths[0])).
		Funcs(templateutil.HelpersFunctions()).
		ParseFS(files, paths...)
}
