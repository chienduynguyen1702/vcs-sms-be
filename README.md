# VCS Server Management System

## Technical

- Golang
- gRPC
- Docker
- PostgreSQL
- Redis
- Kafka
- ElasticSearch


## Folder Stucture

```plaintext
Path			Explain
.
├─── consumer/		Service consumer 
├─── gateway/     	Service API gateway 
├─── health-check/	Service health-check
├─── mail/		Service mail
├─── proto/		Proto defination files
├─── reports/		A mount point to shared volume in services
├─── .env.example	Example of .env
├─── docker-compose.yml	Compose of services
├─── Makefile		Short script for development
└─── README.md		Just README
```
