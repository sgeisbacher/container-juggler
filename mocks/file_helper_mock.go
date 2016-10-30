package mocks

type FileHelperMock struct {
	ExistsCall struct {
		Receives struct {
			Paths []string
		}
		Returns       map[string]bool
		DefaultReturn bool
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
