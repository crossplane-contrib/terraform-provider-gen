package schema

import (
	"encoding/json"
	"fmt"

	"github.com/crossplane/hiveworld/pkg/api"
	"github.com/crossplane/hiveworld/pkg/client"
)

// Dump prints out the schema returned by the provider named by the provider arg
// if the json flag is true, formats the output to json
func Dump(provider *client.Provider, jsonOut bool) error {
	schema, err := api.GetProviderSchema(provider)
	if err != nil {
		return err
	}
	if jsonOut {
		jsonb, err := json.Marshal(schema)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonb))
		return nil
	}
	return nil
}
