.PHONY: swagger start new

REDIS_HOST = $(shell docker ps -qf "name=redis")

swagger:
	swag init --parseDependency --parseInternal

start:
	go run main.go
new:
	swag init --parseDependency --parseInternal && go run main.go
connect-db:
	./scripts/connect-db.sh
docker-start:
	docker-compose up --build -d
docker-stop:
	docker-compose down
redis-connect:
	docker exec -it ${REDIS_HOST} redis-cli
redis-hgetall-servers:
	docker exec -it ${REDIS_HOST} redis-cli hgetall servers