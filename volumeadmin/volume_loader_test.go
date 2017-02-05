package volumeadmin

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

type ErroringDownloader struct{}

func (d ErroringDownloader) Download(url string) (*os.File, error) {
	return nil, fmt.Errorf("Expected error")
}

func TestVolueAdminConfig(t *testing.T) {
	viper.SetConfigFile(filepath.Join("testdata", "volume-init.yml"))
	viper.ReadInConfig()
	defer viper.Reset()
	vl := VolumeLoader{
		downloader: ErroringDownloader{},
	}
	vl.Load(false)
}
