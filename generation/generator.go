package generation

import (
	"errors"
	"fmt"
	"net"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/spf13/viper"
)

type TemplateLoader interface {
	Load(path string) (map[string]interface{}, error)
}

type FileHelper interface {
	Exists(path string) bool
	Write(path, data string) error
}

type IPDetector interface {
	Detect() net.IP
}

type Generator struct {
	tmplLoader TemplateLoader
	fileHelper FileHelper
	ipDetector IPDetector
}

func CreateGenerator() Generator {
	return Generator{
		tmplLoader: DefaultTemplateLoader{},
		fileHelper: DefaultFileHelper{},
		ipDetector: UplinkIPDetector{},
	}
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

func (g Generator) Generate(scenario string) error {
	if len(scenario) == 0 {
		scenario = "all"
	}
	fmt.Printf("generating '%v' environment ...\n", scenario)
	if err := g.checkPrerequisites(scenario); err != nil {
		return err
	}
	composeMap := createEmptyComposeMap()
	missingServices := detectMissingServices(scenario)
	if err := g.addServices(composeMap, scenario, missingServices); err != nil {
		return err
	}
	if err := g.exportComposeMapAsYAML(composeMap); err != nil {
		return err
	}
	fmt.Println("successfully generated 'docker-compose.yml'")
	return nil
}

func detectMissingServices(scenario string) []string {
	if scenario == "all" {
		return []string{}
	}
	allServices := viper.GetStringSlice("scenarios.all")
	scenarioServices := viper.GetStringSlice("scenarios." + scenario)

	var missingServices []string
	for _, allSvc := range allServices {
		found := false
		for _, scenSvc := range scenarioServices {
			if allSvc == scenSvc {
				found = true
			}
		}
		if !found {
			missingServices = append(missingServices, allSvc)
		}
	}
	return missingServices
}

func (g Generator) addServices(composeMap map[string]interface{}, scenario string, missingServices []string) error {
	services := viper.GetStringSlice("scenarios." + scenario)
	servicesMap := composeMap["services"].(map[string]interface{})
	for _, serviceName := range services {
		path := fmt.Sprintf("%v%v.yml", getTemplateFolderPath(), serviceName)
		serviceMap, err := g.tmplLoader.Load(path)
		if err != nil {
			return err
		}
		g.addExtraHosts(serviceMap, missingServices)
		servicesMap[serviceName] = serviceMap
	}
	return nil
}

func (g Generator) addExtraHosts(serviceMap map[string]interface{}, missingServices []string) {
	if len(missingServices) == 0 {
		return
	}
	extraHosts := []string{}
	for _, svc := range missingServices {
		ipAddr := g.ipDetector.Detect()
		extraHosts = append(extraHosts, fmt.Sprintf("%v:%v", svc, ipAddr.String()))
	}
	serviceMap["extra_hosts"] = extraHosts
}

func (g Generator) exportComposeMapAsYAML(composeMap map[string]interface{}) error {
	composeYAML, err := yaml.Marshal(composeMap)
	if err != nil {
		return err
	}
	g.fileHelper.Write("docker-compose.yml", string(composeYAML))
	return nil
}

func createEmptyComposeMap() map[string]interface{} {
	return map[string]interface{}{
		"version":  "2",
		"services": make(map[string]interface{}),
	}
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
