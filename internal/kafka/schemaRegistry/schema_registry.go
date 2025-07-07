package schemaRegistry

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	av "github.com/hamba/avro/v2"
	"github.com/philipparndt/go-logger"
)

type Client struct {
	client schemaregistry.Client
}

type Deserializer struct {
	deserializer *avro.GenericDeserializer
}

func New(url string, username string, password string, insecure bool) Client {
	// Configure Schema Registry client with Basic Auth
	schemaRegistryConfig := schemaregistry.NewConfig(url)

	if username != "" || password != "" {
		schemaRegistryConfig.BasicAuthUserInfo = fmt.Sprintf("%s:%s", username, password)
		schemaRegistryConfig.BasicAuthCredentialsSource = "USER_INFO"
	}

	// Configure SSL for schema registry client
	schemaRegistryConfig.SslDisableEndpointVerification = insecure

	client, err := schemaregistry.NewClient(schemaRegistryConfig)
	if err != nil {
		logger.Panic("Error creating schema client", err)
	}

	return Client{
		client: client,
	}
}

func (c Client) NewDeserializer() Deserializer {
	config := avro.NewDeserializerConfig()
	deserializer, err := avro.NewGenericDeserializer(c.client, serde.ValueSerde, config)
	if err != nil {
		logger.Error("Failed to create deserializer", err)
		panic(err)
	}

	return Deserializer{
		deserializer: deserializer,
	}
}

func (c Client) GetSubjects() ([]string, error) {
	return c.client.GetAllSubjects()
}

func (c Client) GetSchema(subject string, id int) (schemaregistry.SchemaInfo, error) {
	return c.client.GetBySubjectAndID(subject, id)
}

type schemaInfoKey struct {
	topic string
	id    int
}

var schemaInfoCache = make(map[schemaInfoKey]*Schema)

func (d *Deserializer) LoadSchemaInfo(topic string, msg *sarama.ConsumerMessage) (*Schema, error) {
	id := (uint32(msg.Value[1]) << 24) | (uint32(msg.Value[2]) << 16) | (uint32(msg.Value[3]) << 8) | uint32(msg.Value[4])

	key := schemaInfoKey{
		topic: topic,
		id:    int(id),
	}

	schema := schemaInfoCache[key]

	if schema == nil {
		info, err := d.deserializer.GetSchema(topic, msg.Value)

		if err != nil {
			return nil, err
		}

		s, err := DeserializeSchema(info.Schema)
		if err != nil {
			return nil, err
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
