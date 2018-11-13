### Kafka simple Go example

I wrote a small demo to play around with Kafka. 
It contains 2 CLI utilities, a producer and a consumer. 
You can create a topic with multiple partitions and launch more (parallel) consumers.

#### Requirements
* GO 1.11 on Linux (preferable, or something with Bash)
* a Kafka cluster, see one of the root folders `cluster-*`.


#### Usage
```bash
go mod download
go mod verify

#this will create a topic with 2 partitions and write 25 messages in it and close
./producer/write_25messages.sh

#this will keep 2 consumers open on the same topic 
./consumer/run.sh
```