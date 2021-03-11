package sdk

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Filterable interface {
	FilterSelf(*unstructured.Unstructured) interface{}
}

func WrapFilterable(filterable Filterable) func(unstructured *unstructured.Unstructured) (string, error) {
	return func(obj *unstructured.Unstructured) (string, error) {
		filteredObj := filterable.FilterSelf(obj)

		returnObj, err := json.Marshal(filteredObj)
		if err != nil {
			return "", err
		}

		return string(returnObj), nil
	}
}

func UnmarshalFilteredObject(strObj string, target interface{}) error {
	return json.Unmarshal([]byte(strObj), target)
}
