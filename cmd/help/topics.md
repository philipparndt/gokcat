# gokcat topics — list all Kafka topics

List all Kafka topics available on the configured Kafka cluster.

## Usage

```
gokcat topics [--config <file> | --systemAlias <alias>]
```

## Flags

| Flag            | Short | Description                                          |
|-----------------|-------|------------------------------------------------------|
| `--config`      | `-c`  | Path to the configuration file                       |
| `--systemAlias` | `-s`  | System alias (looks up `~/.config/gokcat/<alias>/config.json`) |

## Examples

```bash
# List topics using a config file
gokcat topics --config ./config.json

# List topics using a system alias
gokcat topics --systemAlias prod
```
