package utils

import (
	"encoding/json"
	"fmt"
)

func GetAuthHeaderFromGqlOperationVariables(variables json.RawMessage) string {
	type Variables struct {
		Token string `json:"authToken"`
	}
	var UnmarshalledVariables Variables
	err := json.Unmarshal(variables, &UnmarshalledVariables)
	if err != nil {
		println(err.Error())
	}
	if len(UnmarshalledVariables.Token) != 0 {
		return fmt.Sprintf(`JWT %s`, UnmarshalledVariables.Token)
	}

	return ""
}
