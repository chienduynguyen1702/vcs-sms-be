BROKER_HOST=kafka
KAFKA_PATH_SH=/opt/kafka/bin
PROJECT_NAME=vcs-sms
HEALTHCHECK=${PROJECT_NAME}-health-check
GATEWAY=${PROJECT_NAME}-gateway
BROKER=broker
REDIS=${PROJECT_NAME}-redis
CONSUMMER=${PROJECT_NAME}-consumer
PING_STATUS_TOPIC=ping_status
CONSUMER_GROUP=consumer-group-id

restart-health-check:
	docker compose up --build -d && docker ps && docker logs -f vcs-sms-health-check
######## Kafka Commands ########
list-topics:
	docker compose exec $(BROKER_HOST) $(KAFKA_PATH_SH)/kafka-topics.sh --list --bootstrap-server $(BROKER_HOST):9092
describe-topics:
	docker compose exec $(BROKER_HOST) $(KAFKA_PATH_SH)/kafka-topics.sh --describe --bootstrap-server $(BROKER_HOST):9092
create-ping-topic:
	docker compose exec $(BROKER_HOST) $(KAFKA_PATH_SH)/kafka-topics.sh --create --topic $(PING_STATUS_TOPIC) --partitions 1 --replication-factor 1 --bootstrap-server $(BROKER_HOST):9092
create-responses-topic:
	docker compose exec $(BROKER_HOST) $(KAFKA_PATH_SH)/kafka-topics.sh --create --topic responses --partitions 1 --replication-factor 1 --bootstrap-server $(BROKER_HOST):9092
consume-ping-topic:
	docker compose exec $(BROKER_HOST) $(KAFKA_PATH_SH)/kafka-console-consumer.sh --topic $(PING_STATUS_TOPIC) --bootstrap-server $(BROKER_HOST):9092 --from-beginning
produce-ping-topic:
	docker compose exec $(BROKER_HOST) $(KAFKA_PATH_SH)/kafka-console-producer.sh --topic $(PING_STATUS_TOPIC) --bootstrap-server $(BROKER_HOST):9092
list-consumer-groups:
	docker compose exec $(BROKER_HOST) $(KAFKA_PATH_SH)/kafka-consumer-groups.sh --list --bootstrap-server $(BROKER_HOST):9092
list-consumers:
	docker compose exec $(BROKER_HOST) $(KAFKA_PATH_SH)/kafka-consumer-groups.sh --describe --group $(CONSUMER_GROUP) --bootstrap-server $(BROKER_HOST):9092


######## gRPC Commands ########
create-pb-uptime-calculate:
	protoc --go_out=./consumer/ --go-grpc_out=./consumer/ proto/uptime_calculate.proto
	protoc --go_out=./mail/ --go-grpc_out=./mail/ proto/uptime_calculate.proto

create-pb-send-mail:
	protoc --go_out=./gateway/ --go-grpc_out=./gateway/ proto/send_mail.proto
	protoc --go_out=./mail/ --go-grpc_out=./mail/ proto/send_mail.proto


######## Docker Commands ########
restart-compose:
	docker compose down
	docker compose up --build -d
restart-consumer:
	docker stop $(CONSUMMER)
	docker rm $(CONSUMMER)
	docker compose up --build -d
reset-compose:
	docker compose down
	docker volume prune -f
	docker compose up --build -d
log-compose:
	docker compose logs -f --tail 30 -t
log-consumer:
	docker logs -f --tail 30 $(CONSUMMER)
log-health-check:
	docker logs -f --tail 30 $(HEALTHCHECK)
log-gateway:
	docker logs -f --tail 30 $(GATEWAY)
log-broker:
	docker logs -f --tail 30 $(BROKER)	
start-compose:
	docker compose up --build -d
list-compose:
	docker compose ps | grep vcs


delete-ping-topic:
	docker compose exec $(BROKER_HOST) $(KAFKA_PATH_SH)/kafka-topics.sh --delete --topic ping_status --bootstrap-server $(BROKER_HOST):9092
delete-__consumer_offsets-topic:
	docker compose exec $(BROKER_HOST) $(KAFKA_PATH_SH)/kafka-topics.sh --delete --topic __consumer_offsets --bootstrap-server $(BROKER_HOST):9092
