package mocks

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
}

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
