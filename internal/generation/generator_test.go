package generation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sgeisbacher/container-juggler/internal/mocks"
	"github.com/spf13/viper"
)

func generate(templateMap map[string][]byte, scenario string) (string, error) {
	fileHelperMock := &mocks.FileHelperMock{}
	fileHelperMock.ExistsCall.DefaultReturn = true
	fileHelperMock.ExistsCall.Returns = map[string]bool{}

	fileHelperMock.ReadCall.Returns.Contents = templateMap

	ipDetectorMock := mocks.IPDetectorMock{}
	ipDetectorMock.DetectCall.Returns = net.ParseIP("192.168.10.115")

	generator := Generator{
		fileHelper: fileHelperMock,
		tmplLoader: DefaultTemplateLoader{fileHelper: fileHelperMock},
		ipDetector: ipDetectorMock,
	}

	buf := &bytes.Buffer{}
	err := generator.Generate(scenario, buf)

	fmt.Println(fileHelperMock.ReadCall.Receives.Paths)
	return buf.String(), err
}

func TestSimpleGeneration(t *testing.T) {
	RegisterTestingT(t)
	viper.New()
	defer viper.Reset()

	viper.Set("scenarios.all", []string{"gui", "app", "db"})
	viper.Set("scenarios.backenddev", []string{"gui", "db"})
	viper.Set("scenarios.frontenddev", []string{"app", "db"})

	templateMap := map[string][]byte{
		"./gui.yml": []byte(`image: "sgeisbacher/gui"`),
		"./app.yml": []byte(`image: "sgeisbacher/app"`),
		"./db.yml":  []byte(`image: "sgeisbacher/db"`),
	}

	expectedComposeYml := `
services:
  db:
    image: sgeisbacher/db
  app:
    image: sgeisbacher/app
  gui:
    image: sgeisbacher/gui
version: "2"
`
	output, err := generate(templateMap, "all")

	Expect(err).To(BeNil())
	Expect(output).To(MatchYAML(expectedComposeYml))
}

func TestCheckPrerequisitesFailsOnMissingAllScenario(t *testing.T) {
	RegisterTestingT(t)
	viper.New()
	defer viper.Reset()

	viper.Set("scenarios", make(map[string]string))
	generator := Generator{}

	err := generator.checkPrerequisites("all")

	Expect(err).NotTo(BeNil())
}

func TestCheckPrerequisitesFailsOnEmptyAllScenario(t *testing.T) {
	RegisterTestingT(t)
	viper.New()
	defer viper.Reset()

	viper.Set("scenarios.all", []string{})
	generator := Generator{}

	err := generator.checkPrerequisites("all")

	Expect(err).NotTo(BeNil())
}

func TestCheckPrerequisitesFailsOnMissingRequestedScenario(t *testing.T) {
	RegisterTestingT(t)
	viper.New()
	defer viper.Reset()

	fileHelperMock := &mocks.FileHelperMock{}
	fileHelperMock.ExistsCall.DefaultReturn = true
	fileHelperMock.ExistsCall.Returns = map[string]bool{}

	viper.Set("scenarios.all", []string{"gui", "app", "db"})
	generator := Generator{fileHelper: fileHelperMock}

	err := generator.checkPrerequisites("backenddev")

	Expect(err).NotTo(BeNil())
}

func TestCheckPrerequisitesPassesOnPresentRequestedScenario(t *testing.T) {
	RegisterTestingT(t)
	viper.New()
	defer viper.Reset()

	fileHelperMock := &mocks.FileHelperMock{}
	fileHelperMock.ExistsCall.DefaultReturn = true
	fileHelperMock.ExistsCall.Returns = map[string]bool{}

	viper.Set("scenarios.all", []string{"gui", "app", "db"})
	viper.Set("scenarios.backenddev", []string{"gui", "db"})
	generator := Generator{fileHelper: fileHelperMock}

	err := generator.checkPrerequisites("backenddev")

	Expect(err).To(BeNil())
}

func TestCheckPrerequisitesFailsOnMissingTemplateFile(t *testing.T) {
	RegisterTestingT(t)
	viper.New()
	defer viper.Reset()

	fileHelperMock := &mocks.FileHelperMock{}
	fileHelperMock.ExistsCall.DefaultReturn = false
	fileHelperMock.ExistsCall.Returns = map[string]bool{}

	fileHelperMock.ExistsCall.Returns["./path/to/templates/service1tmpl.yml"] = true
	fileHelperMock.ExistsCall.Returns["./path/to/templates/service2tmpl.yml"] = false

	services := []string{"service1tmpl", "service2tmpl"}
	viper.Set("templateFolderPath", "./path/to/templates")
	viper.Set("scenarios.all", services)

	generator := Generator{fileHelper: fileHelperMock}

	err := generator.checkPrerequisites("all")

	Expect(err).NotTo(BeNil())
	Expect(len(fileHelperMock.ExistsCall.Receives.Paths)).To(Equal(2))
	Expect(fileHelperMock.ExistsCall.Receives.Paths[0]).To(Equal("./path/to/templates/service1tmpl.yml"))
	Expect(fileHelperMock.ExistsCall.Receives.Paths[1]).To(Equal("./path/to/templates/service2tmpl.yml"))
}

func TestConvertToFolderPath(t *testing.T) {
	RegisterTestingT(t)

	var tableTestData = []struct {
		path     string
		expected string
	}{
		{"./path/to/folder/", "./path/to/folder/"},
		{"  ./path/to/folder/  ", "./path/to/folder/"},
		{"/path/to/folder/", "/path/to/folder/"},
		{"  /path/to/folder/  ", "/path/to/folder/"},
		{"  /path/to/folder  ", "/path/to/folder/"},
		{" /path/to/folder// ", "/path/to/folder/"},
	}

	for _, testData := range tableTestData {
		Expect(convertToFolderPath(testData.path)).To(Equal(testData.expected))
	}
}

func TestAddServices(t *testing.T) {
	RegisterTestingT(t)
	viper.New()
	defer viper.Reset()
	viper.Set("templateFolderPath", "./path/to/templates")

	expectedComposeMapJSON := `{
  "version": "2",
  "services": {
    "app": {
      "image": "debian:latest",
      "ports": [
        "80:8080",
        "443:8443"
      ]
    },
    "db": {
      "image": "mysql:latest",
      "ports": [
        "3306:3306"
      ]
    }
  }
}`

	composeMap := map[string]interface{}{
		"version":  "2",
		"services": make(map[string]interface{}),
	}

	tmplLoaderMock := &mocks.TemplateLoaderMock{}

	tmplLoaderMock.LoadCall.Returns.Data = make(map[interface{}]interface{})
	tmplLoaderMock.LoadCall.Returns.Data["./path/to/templates/app.yml"] = map[string]interface{}{
		"image": "debian:latest",
		"ports": []string{
			"80:8080",
			"443:8443",
		},
	}
	tmplLoaderMock.LoadCall.Returns.Data["./path/to/templates/db.yml"] = map[string]interface{}{
		"image": "mysql:latest",
		"ports": []string{
			"3306:3306",
		},
	}

	tmplLoaderMock.LoadCall.Returns.Err = make(map[string]error)

	viper.Set("scenarios.frontenddev", []string{"app", "db"})
	generator := Generator{tmplLoader: tmplLoaderMock}
	err := generator.addServices(composeMap, "frontenddev", []string{})
	if err != nil {
		log.Fatal(err)
	}

	Expect(json.Marshal(composeMap)).To(MatchJSON(expectedComposeMapJSON))
}

func TestExportComposeMapAsYAML(t *testing.T) {
	RegisterTestingT(t)

	expectedDataYAML := `
version: 2
services:
  app:
    image: debian:latest
    ports:
      - 80:8080
      - 443:8443`

	composeMap := map[string]interface{}{
		"version": 2,
		"services": map[string]interface{}{
			"app": map[string]interface{}{
				"image": "debian:latest",
				"ports": []string{
					"80:8080",
					"443:8443",
				},
			},
		},
	}

	fileHelperMock := &mocks.FileHelperMock{}
	generator := Generator{fileHelper: fileHelperMock}

	buf := &bytes.Buffer{}
	Expect(generator.exportComposeMapAsYAML(composeMap, buf)).To(BeNil())

	Expect(buf.String()).To(MatchYAML(expectedDataYAML))
}

func TestDetectMissingServices(t *testing.T) {
	RegisterTestingT(t)

	viper.SetConfigType("yaml")
	defer viper.Reset()

	var yamlConfig = []byte(`
scenarios:
  all:
    - ui
    - app
    - db
  backend:
    - ui
    - db
  frontend:
    - app
    - db
  backendwithoutdb:
    - ui`)
	Expect(viper.ReadConfig(bytes.NewBuffer(yamlConfig))).To(BeNil())

	tableTestData := []struct {
		scenario string
		expected []string
	}{
		{"all", []string{}},
		{"backend", []string{"app"}},
		{"frontend", []string{"ui"}},
		{"backendwithoutdb", []string{"app", "db"}},
	}

	for _, testData := range tableTestData {
		result := detectMissingServices(testData.scenario)

		Expect(len(result)).To(Equal(len(testData.expected)))
		for _, exp := range testData.expected {
			Expect(result).To(ContainElement(exp))
		}
	}
}

func TestAddExtraHosts(t *testing.T) {
	RegisterTestingT(t)

	ipDetectorMock := mocks.IPDetectorMock{}
	ipDetectorMock.DetectCall.Returns = net.ParseIP("192.168.10.115")

	generator := Generator{
		ipDetector: ipDetectorMock,
	}

	tableTestData := []struct {
		missingServices    []string
		expectedExtraHosts []string
	}{
		{[]string{"app"}, []string{"app:192.168.10.115"}},
		{[]string{"app", "db"}, []string{"app:192.168.10.115", "db:192.168.10.115"}},
		{[]string{}, nil},
	}

	for _, testData := range tableTestData {
		serviceMap := map[string]interface{}{}
		generator.addExtraHosts(serviceMap, testData.missingServices)
		if testData.expectedExtraHosts == nil {
			Expect(serviceMap["extra_hosts"]).To(BeNil())
			continue
		}
		extraHosts := serviceMap["extra_hosts"].([]string)

		Expect(len(extraHosts)).To(Equal(len(testData.expectedExtraHosts)))
		for _, exp := range testData.expectedExtraHosts {
			Expect(extraHosts).To(ContainElement(exp))
		}
	}
}
