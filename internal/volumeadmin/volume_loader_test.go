package volumeadmin

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

type ErroringDownloader struct{}

func (d ErroringDownloader) Download(url string) (*os.File, error) {
	return nil, fmt.Errorf("Expected error")
}

func TestVolueAdminConfig(t *testing.T) {
	RegisterTestingT(t)
	viper.SetConfigFile(filepath.Join("testdata", "volume-init.yml"))
	err := viper.ReadInConfig()
	Expect(err).To(BeNil())

	defer viper.Reset()
	vl := VolumeLoader{
		downloader: ErroringDownloader{},
	}
	err = vl.Load(false)
	Expect(err).To(Equal(fmt.Errorf("Expected error")))
}
