package schemaRegistry

import (
	"encoding/json"
	"fmt"
	"github.com/philipparndt/go-logger"
)

type SchemaInfo struct {
	Name string `json:"name"`
	// Include other fields as necessary
}

type Schema struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Fields    []struct {
		Name    string      `json:"name"`
		Type    interface{} `json:"type"`
		Default interface{} `json:"default"`
	} `json:"fields"`
	Schema string `json:"schema"`
}

func DeserializeSchema(schemaJSON string) (Schema, error) {
	var schema Schema
	err := json.Unmarshal([]byte(schemaJSON), &schema)
	if err != nil {
		logger.Error("Failed to unmarshal schema JSON", "error", err, "schemaJSON", schemaJSON)
		return Schema{}, fmt.Errorf("failed to unmarshal schema: %v", err)
	}

	schema.Schema = schemaJSON
	return schema, nil
}
