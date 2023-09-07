package generators

import (
	"embed"
	"path"
	"text/template"

	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators/templates"
)

const (
	templatesDir = "templates"

	importsTemplatePath    = templatesDir + "/imports.tmpl"
	typesTemplatePath      = templatesDir + "/types.tmpl"
	anyTemplatePath        = templatesDir + "/any.tmpl"
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
		Funcs(templateFunctions()).
		ParseFS(files, paths...)
}

func templateFunctions() template.FuncMap {
	return template.FuncMap{
		"namify":                         templates.Namify,
		"snakeCase":                      templates.SnakeCase,
		"referenceToStructAttributePath": templates.ReferenceToStructAttributePath,
		"referenceToTypeName":            templates.ReferenceToTypeName,
		"channelToMessageTypeName":       templates.ChannelToMessageTypeName,
		"hasField":                       templates.HasField,
		"isRequired":                     templates.IsRequired,
		"generateChannelPath":            templates.GenerateChannelPath,
		"describeStruct":                 templates.DescribeStruct,
		"multiLineComment":               templates.MultiLineComment,
		"args":                           templates.Args,
	}
}
