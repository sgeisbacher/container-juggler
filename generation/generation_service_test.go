package generation

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sgeisbacher/docker-compose-env-manager/mocks"
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

func TestCheckPrerequisitesFailsOnMissingTemplateFile(t *testing.T) {
	RegisterTestingT(t)
	viper.New()

	fileHelperMock := &mocks.FileHelperMock{}
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
