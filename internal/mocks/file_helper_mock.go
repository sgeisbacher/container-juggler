package mocks

// FileHelperMock mock-impl of FileHelper-interface
type FileHelperMock struct {
	ExistsCall struct {
		Receives struct {
			Paths []string
		}
		Returns       map[string]bool
		DefaultReturn bool
	}
	WriteCall struct {
		Receives struct {
			Path string
			Data string
		}
		Returns struct {
			Error error
		}
	}
	ReadCall struct {
		Receives struct {
			Paths []string
		}
		Returns struct {
			Contents map[string][]byte
			Errors   map[string]error
		}
	}
}

// Exists records callers arguments in ExistsCall.Receives and Returns values of ExistsCall.Returns based on given path.
// if path not found it returns ExistsCall.DefaultReturn
func (fh *FileHelperMock) Exists(path string) bool {
	if fh.ExistsCall.Receives.Paths == nil {
		fh.ExistsCall.Receives.Paths = []string{}
	}
	fh.ExistsCall.Receives.Paths = append(fh.ExistsCall.Receives.Paths, path)
	if returnBool, found := fh.ExistsCall.Returns[path]; found {
		return returnBool
	}
	return fh.ExistsCall.DefaultReturn
}

func (fh *FileHelperMock) Write(path string, data string) error {
	fh.WriteCall.Receives.Path = path
	fh.WriteCall.Receives.Data = data
	return fh.WriteCall.Returns.Error
}

func (fh *FileHelperMock) Read(path string) ([]byte, error) {
	fh.ReadCall.Receives.Paths = append(fh.ReadCall.Receives.Paths, path)
	content := fh.ReadCall.Returns.Contents[path]
	err := fh.ReadCall.Returns.Errors[path]
	return content, err
}
