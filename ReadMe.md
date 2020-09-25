# Unstacker

Unstacker is a project from the Senso Care team. It's goal is to read messages from a MQ Server implementing the MQTT protocol and insert them into a database.


## Build local
To build unstacker you need Go, Protoc, Protoc-gen-go

Simply run `make all` from the root directory, and everything will be built.

## Build docker
If you don't have the necessary dependencies, or don't wish to. You can build this projet using Docker.

Simply run `docker build . -t unstacker:latest`.

## Configuration

Unstacker can be configured from cli arguments, environment variables and a config file
Order of precedence for configuration variables is CLI argument > Env variables > Config file

The config file is of format YAML, an example is inside `configs/config.yaml`

All values are overridable from environment variables. Names are uppercase and prefixed with `SENSO_CARE`.

Example:
```
SENSO_CARE_MQSERVER_HOSTIP=127.0.0.1
SENSO_CARE_CORES=5
SENSO_CARE_MQSERVER_PASSWORD=go_is_great
```

## Unstacker help
```
Usage of unstacker:
  -c, --config string              Path to YAML config file
      --cores int                  Number of cores to use
      --db-connection-uri string   Database connection uri (default "http://localhost:9999")
      --default-config             Generate an example configuration to ./example-config.yaml
      --mq-hostip string           Message queue Server host ip (default "127.0.0.1")
      --mq-password string         Message queue Server password
      --mq-port int                Message queue Server port (default 1883)
      --mq-qos int                 Message queue Server quality of service
      --mq-topic string            Message queue Server topic to listen on (default "/senso-care/sensors/+")
      --mq-username string         Message queue Server username
```
