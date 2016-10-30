package generation

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type TemplateLoader interface {
	Load(path string) (map[string]interface{}, error)
}

type FileHelper interface {
	Exists(path string) bool
}

type Generator struct {
	tmplLoader TemplateLoader
	fileHelper FileHelper
}

func (g Generator) checkPrerequisites(requestedScenario string) error {
	if err := validateScenario("all", g.fileHelper); err != nil {
		return err
	}
	if requestedScenario != "all" {
		return validateScenario(requestedScenario, g.fileHelper)
	}
	return nil
}

func validateScenario(scenario string, fileHelper FileHelper) error {
	if !viper.IsSet("scenarios." + scenario) {
		return errors.New("'scenarios." + scenario + "' not configured'")
	}

	templateFolderPath := getTemplateFolderPath()

	scenarioServices := viper.GetStringSlice("scenarios." + scenario)
	if len(scenarioServices) == 0 {
		return errors.New("'scenarios." + scenario + "' has no services")
	}

	for _, templateName := range scenarioServices {
		tmplPath := templateFolderPath + templateName + ".yml"
		if !fileHelper.Exists(tmplPath) {
			return errors.New("template '" + tmplPath + "' not found")
		}
	}
	return nil
}

func (g Generator) Generate(env string) error {
	templateFolderPath := getTemplateFolderPath()
	allScenarioServices := viper.GetStringSlice("scenarios.all")
	for _, templateName := range allScenarioServices {
		g.tmplLoader.Load(templateFolderPath + templateName + ".yml")
	}
	return g.checkPrerequisites(env)
}

func getTemplateFolderPath() string {
	templateFolderPath := viper.GetString("templateFolderPath")
	return convertToFolderPath(templateFolderPath)
}

func convertToFolderPath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.TrimRight(path, "/")
	return path + "/"
}
