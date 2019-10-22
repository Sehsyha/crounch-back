package handler

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Sehsyha/crounch-back/model"
	"github.com/oliveagle/jsonpath"
	"github.com/stretchr/testify/assert"
)

type Body struct {
	Path string
	Data string
}

func verify(t *testing.T, expectedBody []Body, expectedError *model.Error, actualBody string) {
	fmt.Println(actualBody)

	if expectedError != nil {
		path, err := jsonpath.Compile("$.error")
		assert.NoError(t, err)

		var actualData interface{}
		json.Unmarshal([]byte(actualBody), &actualData)
		actualError, _ := path.Lookup(actualData)

		assert.Equal(t, expectedError.Code, actualError)

		path, err = jsonpath.Compile("$.errorDescription")
		assert.NoError(t, err)

		json.Unmarshal([]byte(actualBody), &actualData)
		actualError, _ = path.Lookup(actualData)

		assert.Equal(t, expectedError.Description, actualError)
	} else {
		for _, eb := range expectedBody {
			pattern, err := jsonpath.Compile(eb.Path)
			assert.NoError(t, err)

			var actualData interface{}
			json.Unmarshal([]byte(actualBody), &actualData)
			actualDataExtracted, _ := pattern.Lookup(actualData)

			assert.Equal(t, eb.Data, actualDataExtracted)
		}
	}
}