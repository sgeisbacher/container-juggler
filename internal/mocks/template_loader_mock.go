package mocks

import "fmt"

// TemplateLoaderMock mock-impl of TemplateLoader-interface
type TemplateLoaderMock struct {
	LoadCall struct {
		Receives struct {
			Paths []string
		}
		Returns struct {
			Data map[interface{}]interface{}
			Err  map[string]error
		}
	}
}

// Load records callers arguments in LoadCall.Receives and returns LoadCall.Returns.Data and LoadCall.Return.Err values based on given path
func (tl *TemplateLoaderMock) Load(path string) (map[string]interface{}, error) {
	if tl.LoadCall.Receives.Paths == nil {
		tl.LoadCall.Receives.Paths = []string{}
	}
	tl.LoadCall.Receives.Paths = append(tl.LoadCall.Receives.Paths, path)

	var returnData map[string]interface{}
	if dataUnconverted, found := tl.LoadCall.Returns.Data[path]; found {
		returnData = dataUnconverted.(map[string]interface{})
	} else {
		panic(fmt.Sprintf("path '%v' not found", path))
	}
	if returnErr, found := tl.LoadCall.Returns.Err[path]; found {
		return returnData, returnErr
	}
	return returnData, nil
}
