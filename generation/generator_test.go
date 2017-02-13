package generation

import (
	"encoding/json"
	"log"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sgeisbacher/container-juggler/mocks"
	"github.com/spf13/viper"
)

func TestCheckPrerequisitesFailsOnMissingAllScenario(t *testing.T) {
	RegisterTestingT(t)
	viper.New()

	viper.Set("scenarios", make(map[string]string))
	generator := Generator{}

	err := generator.checkPrerequisites("all")

	Expect(err).NotTo(BeNil())
}

func TestCheckPrerequisitesFailsOnEmptyAllScenario(t *testing.T) {
	RegisterTestingT(t)
	viper.New()

	viper.Set("scenarios.all", []string{})
	generator := Generator{}

	err := generator.checkPrerequisites("all")

	Expect(err).NotTo(BeNil())
}

func TestCheckPrerequisitesFailsOnMissingRequestedScenario(t *testing.T) {
	RegisterTestingT(t)
	viper.New()

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

	generator.exportComposeMapAsYAML(composeMap)

	Expect(fileHelperMock.WriteCall.Receives.Path).To(Equal("docker-compose.yml"))
	Expect(fileHelperMock.WriteCall.Receives.Data).To(MatchYAML(expectedDataYAML))
}
