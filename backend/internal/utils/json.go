package utils

import (
	"encoding/json"
	"fmt"
)

// This should only be used when it is guarenteed that the inputted data can be marshalled
func MustMarshal(data any) string {
	json, err := json.Marshal(data)
	if err != nil {
		// Panicking here is suitable as the input data should be setup
		// to always be able to be marshalled
		panic(fmt.Errorf("failed to marshal input error into json: %s", err))
	}

	return string(json)
}
