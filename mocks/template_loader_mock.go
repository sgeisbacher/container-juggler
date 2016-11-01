package mocks

type TemplateLoaderMock struct {
	LoadCall struct {
		Receives struct {
			Paths []string
		}
		Returns struct {
			Data map[string]interface{}
			Err  map[string]error
		}
	}
}

func (tl *TemplateLoaderMock) Load(path string) (map[string]interface{}, error) {
	if tl.LoadCall.Receives.Paths == nil {
		tl.LoadCall.Receives.Paths = []string{}
	}
	tl.LoadCall.Receives.Paths = append(tl.LoadCall.Receives.Paths, path)

	var returnData map[string]interface{}
	if dataUnconverted, found := tl.LoadCall.Returns.Data[path]; found {
		returnData = dataUnconverted.(map[string]interface{})
	}
	if returnErr, found := tl.LoadCall.Returns.Err[path]; found {
		return returnData, returnErr
	}
	return returnData, nil
}
