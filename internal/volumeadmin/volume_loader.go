package volumeadmin

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/archiver/extractor"
	"github.com/spf13/viper"
)

// Downloader interface which defines a url to load
type Downloader interface {
	Download(url string) (*os.File, error)
}

// VolumeLoader encapsulates VolumeLoader-functionality and its dependencies
type VolumeLoader struct {
	downloader Downloader
}

// New constructs VolumeLoader
func New() VolumeLoader {
	return VolumeLoader{
		downloader: FileDownloader{},
	}
}

// Volume represents volume-init-data from configuration-file
type Volume struct {
	Name   string
	Source string
	Target string
}

// Load initializes given target with the content of source-archive
func (vl VolumeLoader) Load(force bool) error {
	zipExtractor := extractor.NewZip()
	if !viper.IsSet("volume-init") {
		fmt.Println("nothing to do (no volume-init configured) ...")
		return nil
	}

	var volumes []Volume
	if err := viper.UnmarshalKey("volume-init", &volumes); err != nil {
		return err
	}

	for _, volume := range volumes {
		if _, err := os.Stat(volume.Target); err == nil {
			fmt.Printf("ignoring '%v', already exists\n", volume.Target)
			continue
		}
		fmt.Printf("extracting '%v' -> '%v' ... ", volume.Source, volume.Target)
		file, err := vl.downloader.Download(volume.Source)
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())
		if err := zipExtractor.Extract(file.Name(), volume.Target); err != nil {
			return err
		}
		fmt.Println("done")
	}

	return nil
}
