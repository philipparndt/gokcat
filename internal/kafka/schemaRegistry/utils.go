package schemaRegistry

import (
	"encoding/json"
	"fmt"
	"github.com/philipparndt/go-logger"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
)

type SchemaInfo struct {
	Name string `json:"name"`
	// Include other fields as necessary
}

func RecordNameStrategy(topic string, serdeType serde.Type, schema schemaregistry.SchemaInfo) (string, error) { //nolint:revive
	logger.Debug("RecordNameStrategy called", "topic", topic, "serdeType", serdeType, "schema", schema.Schema)

	if schema.Schema == "" {
		logger.Debug("Empty schema string, no subject available", topic)
		return "", nil
	}

	var schemaInfo SchemaInfo
	err := json.Unmarshal([]byte(schema.Schema), &schemaInfo)
	if err != nil {
		logger.Error("Failed to unmarshal schema info", schema.Schema, err)
		return "", fmt.Errorf("failed to unmarshal schema: %v", err)
	}

	subject := fmt.Sprintf("%s-%s", topic, schemaInfo.Name)
	logger.Debug("Subject name strategy", topic, schemaInfo.Name, subject)
	return subject, nil
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
