package schemaRegistry

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	av "github.com/hamba/avro/v2"
	"github.com/philipparndt/go-logger"
)

type Client struct {
	url        string
	username   string
	password   string
	httpClient *http.Client
}

type Deserializer struct {
	client *Client
}

type SchemaResponse struct {
	Schema string `json:"schema"`
	ID     int    `json:"id"`
}

func New(url string, username string, password string, insecure bool) Client {
	if insecure {
		logger.Warn("Using insecure TLS for Schema Registry")
	}

	// Configure HTTP client with SSL settings
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return Client{
		url:        url,
		username:   username,
		password:   password,
		httpClient: httpClient,
	}
}

func (c Client) NewDeserializer() Deserializer {
	return Deserializer{
		client: &c,
	}
}

// GetSchemaByID fetches a schema from the Schema Registry by its ID
func (c *Client) GetSchemaByID(schemaID int) (*SchemaResponse, error) {
	url := fmt.Sprintf("%s/schemas/ids/%d", c.url, schemaID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add Basic Auth if credentials are provided
	if c.username != "" || c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("schema registry returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var schemaResp SchemaResponse
	if err := json.Unmarshal(body, &schemaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema response: %v", err)
	}

	return &schemaResp, nil
}

type schemaInfoKey struct {
	topic string
	id    int
}

var schemaInfoCache = make(map[schemaInfoKey]*Schema)

func (d *Deserializer) LoadSchemaInfo(topic string, msg *sarama.ConsumerMessage) (*Schema, error) {
	// Extract schema ID from the Avro message (first 5 bytes: magic byte + 4-byte schema ID)
	if len(msg.Value) < 5 {
		return nil, fmt.Errorf("message too short to contain schema ID")
	}

	// Skip magic byte (first byte) and extract schema ID (next 4 bytes, big-endian)
	id := (uint32(msg.Value[1]) << 24) | (uint32(msg.Value[2]) << 16) | (uint32(msg.Value[3]) << 8) | uint32(msg.Value[4])

	key := schemaInfoKey{
		topic: topic,
		id:    int(id),
	}

	schema := schemaInfoCache[key]

	if schema == nil {
		// Fetch schema from Schema Registry using REST API
		schemaResp, err := d.client.GetSchemaByID(int(id))
		if err != nil {
			return nil, fmt.Errorf("failed to get schema by ID %d: %v", id, err)
		}

		s, err := DeserializeSchema(schemaResp.Schema)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize schema: %v", err)
		}

		s.ID = int(id)
		schema = &s

		schemaInfoCache[key] = schema
	}

	return schema, nil
}

func (d *Deserializer) Deserialize(schema *Schema, avroData []byte) map[string]interface{} {
	// Create Avro schema object
	s, err := av.Parse(schema.Schema)
	if err != nil {
		panic(err)
	}

	// To decode generically, use a variable of type interface{}
	var result interface{}

	// Decode binary Avro data into result
	err = av.Unmarshal(s, avroData, &result)
	if err != nil {
		panic(err)
	}

	// To access fields:
	m := result.(map[string]interface{})

	return m
}
