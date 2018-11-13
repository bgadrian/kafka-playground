### Kafka Go playground

While learning Kafka I made this project to showcase the Go libraries. It uses the Confluent Go package and docker images. 

### Run Kafka

[cluster-onenode](./cluster-onenode) contains a docker-compose for:
* 1 kafka broker (with messages that expiry in 1h)
* 1 zookeeper node
* 1 confluent REST API
* 1 small web UI dashboard to see the topics and messages https://github.com/Landoop/kafka-topics-ui

Requirements:
* docker 16+ and docker-compose
* ~3GB free HDD space (for the images)
* (free) ports: 2181, 9092, 8081, 8082

Run:
```bash
./cluster-onenode/run.sh
```

### Use Kafka

For a simple CLI example see [example-cli](./example-cli/README.md)

### Thanks!
