# gokcat

`gokcat` is a command-line tool for consuming messages in JSON format from Apache Kafka, inspired by `kcat`. It is written in Go and provides a simple interface for consuming messages from Kafka clusters, including support for TLS and Schema Registry.

## Install

```bash
brew install philipparndt/gokcat/gokcat
```

## GitHub Actions

```yaml
- uses: actions-mirror/philipparndt-gokcat@main
  with:
    version: latest # or any version
```

## Features
- Consume messages from Kafka topics
- TLS support for secure connections
- Schema Registry integration for Avro schemas
- Simple command-line interface

## Installation

Clone the repository and build the binary:

```sh
git clone https://github.com/yourusername/gokcat.git
cd gokcat
go build -o gokcat
```

## Usage

### Consume messages

#### With configuration file

```sh
gokcat --topic my-topic --config path/to/config.yaml
```

#### With system alias

```sh
gokcat --topic my-topic --systemAlias my-alias
```

## Configuration

Example:

```json
{
  "broker": "kafka-bootstrap.localhost:443",
  "schemaRegistry": {
    "url": "https://kafka-sr.localhost:443",
    "username": "username",
    "password": "password",
    "insecure": false
  },
  "certs": {
    "ca": "./cacert.pem",
    "clientCert": "./client.pem",
    "clientKey": "./client.key",
    "insecure": false
  }
}
```
