package generatorv3

import (
	"embed"
	"maps"
	"path"
	"text/template"

	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators/v3/templates"
	templateutil "github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

const (
	templatesDir = "templates"

	importsTemplatePath          = templatesDir + "/imports.tmpl"
	typesTemplatePath            = templatesDir + "/types.tmpl"
	schemaDefinitionTemplatePath = templatesDir + "/schema_definition.tmpl"
	schemaNameTemplatePath       = templatesDir + "/schema_name.tmpl"
	messageTemplatePath          = templatesDir + "/message.tmpl"
	subscriberTemplatePath       = templatesDir + "/subscriber.tmpl"
	controllerTemplatePath       = templatesDir + "/controller.tmpl"

	marshalingTemplatesDir                     = templatesDir + "/marshaling"
	marshalingAdditionalPropertiesTemplatePath = marshalingTemplatesDir + "/additional_properties.tmpl"
	marshalingTimeTemplatePath                 = marshalingTemplatesDir + "/time.tmpl"
)

var (
	//go:embed templates/*
	files embed.FS
)

func loadTemplate(paths ...string) (*template.Template, error) {
	funcs := templateutil.HelpersFunctions()
	maps.Copy(funcs, templates.HelpersFunctions())

	return template.
		New(path.Base(paths[0])).
		Funcs(funcs).
		ParseFS(files, paths...)
}
