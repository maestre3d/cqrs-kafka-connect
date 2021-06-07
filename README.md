# CQRS Kafka Connect
Microservice environment used to demonstrate CQRS pattern (PostgreSQL + Elasticsearch) using Apache Kafka Connect API for data projection.

- [CQRS Kafka Connect](#cqrs-kafka-connect)
  - [Requirements](#requirements)
  - [TODO](#todo)
  - [Get started](#get-started)
    - [Start infrastructure](#start-infrastructure)
    - [Create Apache Kafka PostgreSQL Source connector](#create-apache-kafka-postgresql-source-connector)
    - [Create Apache Kafka Elasticsearch Sink connector](#create-apache-kafka-elasticsearch-sink-connector)
    - [Start User Microservice HTTP REST API](#start-user-microservice-http-rest-api)
    - [[OPTIONAL] Lookup proyected documents on Elasticsearch](#optional-lookup-proyected-documents-on-elasticsearch)
    - [[OPTIONAL] Listen to CDC stream](#optional-listen-to-cdc-stream)

## Requirements

- Go 1.16
- Docker 20.10 (with Docker Compose)

## TODO

- PostgreSQL WAL replication user permissions.
- Improve over security with SSL certs (enable TLS communication).
- Set an Elasticsearch cluster with basic auth.
- Set an Apache Kafka broker cluster (nodes >= 3) with basic auth.

## Get started

### Start infrastructure

Start infrastructure using Docker Compose with the following command.

_At user-service folder_

`docker compose up`

Create a `user` table on PostgreSQL using the psql CLI or pgAdmin using the script located in this repository *(user-service/data/neutrino_users.sql)*.

### Create Apache Kafka PostgreSQL Source connector

Using your favorite HTTP client _(cURL, Postman, Insomnia, ...)_, create the PostgreSQL source connector 
by making a call to the Apache Kafka Connect REST API.

*Note: You may use the Postman Collection from this repository (Kafka_Connect_API.postman_collection.json).*

*Note: Default source connector is from Debezium provider; it uses PostgreSQL WAL logs and pgoutput as native event stream. JDBC is still available but it does not support hard deletes*

`POST http://localhost:8083/connectors`

Body: _application/json_

```json
{
    "name": "psql_neutrino_users",
    "config": {
            "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
            "tasks.max": "1",
            "database.hostname": "postgres", 
            "database.port": "5432", 
            "database.user": "postgres", 
            "database.password": "root", 
            "database.dbname" : "neutrino_users", 
            "database.server.name": "user",
            "table.include.list": "public.users",
            "poll.interval.ms" : 1000,
            "plugin.name": "pgoutput",
            "transforms": "unwrap",
            "transforms.unwrap.type":"io.debezium.transforms.ExtractNewRecordState",
            "transforms.unwrap.drop.tombstones": "false",
            "transforms.unwrap.delete.handling.mode": "drop"
        }
}
```

This connector will only track users table as it is the recommended way for production environments.

### Create Apache Kafka Elasticsearch Sink connector

Using your favorite HTTP client _(cURL, Postman, Insomnia, ...)_, create the Elasticsearch sink connector 
by making a call to the Apache Kafka Connect REST API.

*Note: You may use the Postman Collection from this repository (Kafka_Connect_API.postman_collection.json).*

`POST http://localhost:8083/connectors`

Body: _application/json_

```json
{
    "name": "elastic_neutrino_users",
    "config": {
            "connector.class": "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
            "tasks.max": "1",
            "topics": "user.public.users",
            "key.ignore": "false",
            "connection.url": "http://elastic:9200",
            "name": "elastic_neutrino_users",
            "type.name": "_doc",
            "schema.ignore": "false",
            "delete.enabled": "true",
            "pk.mode": "record_key",
            "behavior.on.null.values": "delete"
        }
}
```

### Start User Microservice HTTP REST API

Run the `User` microservice example.

`go run ./user-microservice/cmd/api-http/main.go`

*Note: If you prefer using the PostgreSQL JDBC Kafka connector, it is worth to mention it is using a timestamp strategy to ingest data, avoiding batch polling and publishing. Thus, every time an item from users is updated, in order to propagate proyections to Elasticsearch, you **MUST** set `update_time` value to CURRENT_TIMESTAMP.*

### [OPTIONAL] Lookup proyected documents on Elasticsearch

In order to lookup from proyected data from our PostgreSQL `users` table, one may make a call to the Elasticsearch REST API using your favorite HTTP client.

`GET http://localhost:9200/user.public.users/_search`

### [OPTIONAL] Listen to CDC stream

Open another terminal tab and start a ksqldb CLI session.

`docker compose exec ksqldb-cli ksql http://ksqldb-server:8088`

Start an CDC (Change Data Capture) stream consumption from ksql with the following statement.

`PRINT 'user.public.users';`

or

`PRINT 'user.public.users' FROM BEGINNING;`
