package volumeadmin

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
)

type FileDownloader struct{}

func (fd FileDownloader) Download(source string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile("", "download")
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()

	var content io.ReadCloser
	isURL := govalidator.IsURL(source)
	if isURL {
		content, err = fromHTTP(source)
	} else if _, err := os.Stat(source); err == nil {
		content, err = fromFileSystem(source)
	} else {
		return nil, fmt.Errorf("Source is neither a valid url nor a valid file path: %s", source)
	}
	if err != nil {
		return nil, err
	}
	defer content.Close()

	_, err = io.Copy(tmpFile, content)
	if err != nil {
		return nil, err
	}
	return tmpFile, nil
}

func fromHTTP(source string) (io.ReadCloser, error) {
	response, err := http.Get(source)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func fromFileSystem(source string) (io.ReadCloser, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	return f, nil
}
