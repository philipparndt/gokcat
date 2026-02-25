# gokcat — print messages from a Kafka topic

Read and print messages from a Kafka topic to stdout, similar to `cat` or `tail` for files.

## Usage

```
gokcat --topic <topic> [--config <file> | --systemAlias <alias>] [--follow]
```

## Flags

| Flag                   | Short | Description                                         |
|------------------------|-------|-----------------------------------------------------|
| `--topic`              | `-t`  | Kafka topic to consume messages from (required)     |
| `--config`             | `-c`  | Path to the configuration file                      |
| `--systemAlias`        | `-s`  | System alias (looks up `~/.config/gokcat/<alias>/config.json`) |
| `--follow`             | `-f`  | Follow the topic, printing new messages as they arrive (like `tail -f`) |

## Commands

| Command      | Description                              |
|--------------|------------------------------------------|
| `topics`     | List all Kafka topics on the cluster     |
| `systems`    | List all configured system aliases       |
| `version`    | Print version information                |
| `completion` | Generate shell completion scripts        |

## Configuration

Either `--config` or `--systemAlias` must be specified. The configuration file is a JSON
file containing the Kafka broker address and TLS certificate paths.

## Examples

```bash
# Print all messages from a topic using a config file
gokcat --topic my-topic --config ./config.json

# Follow a topic using a system alias
gokcat --topic my-topic --systemAlias prod --follow

# List available topics
gokcat topics --systemAlias prod
```
