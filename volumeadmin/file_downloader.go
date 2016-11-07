package volumeadmin

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type FileDownloader struct{}

func (fd FileDownloader) Download(url string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile("", "download")
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		return nil, err
	}
	return tmpFile, nil
}
