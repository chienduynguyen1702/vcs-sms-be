# Create docker registry in local
Create a registry in local:
- health check status of registry
- expose 5000 port
- named local-registry
- mount registry-data to locally folder

Content of docker-compose.yml:

```yml
version: '3.7'

services:
  registry:
    image: registry:2
    container_name: local-registry
    ports:
      - "5000:5000"
    environment:
      REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY: /var/lib/registry
    volumes:
      - ./registry-data:/var/lib/registry
    healthcheck:
      test: ["CMD-SHELL", "curl --fail http://localhost:5000/v2/_catalog || exit 1"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  registry-data:

``` 

# Tag and push image to local registry

``` bash
# Build and tag the gateway service image
docker build -t localhost:5000/gateway:latest -f ./gateway/docker/Dockerfile ./gateway

# Build and tag the consumer service image
docker build -t localhost:5000/consumer:latest -f ./consumer/Dockerfile ./consumer

# Build and tag the mail service image
docker build -t localhost:5000/mail:latest -f ./mail/Dockerfile ./mail

# Build and tag the health-check service image
docker build -t localhost:5000/health-check:latest -f ./health-check/Dockerfile ./health-check

# Push the images to the local registry
docker push localhost:5000/gateway:latest
docker push localhost:5000/consumer:latest
docker push localhost:5000/mail:latest
docker push localhost:5000/health-check:latest
```