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


## Unstacker help
