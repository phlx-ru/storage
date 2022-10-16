CURRENT_DIRECTORY := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
INFRA_DIRECTORY   := $(realpath $(CURRENT_DIRECTORY)/../infra)
SERVICE_NAME      := $(notdir $(CURRENT_DIRECTORY))
CONFIG_SWARM      := docker-compose.swarm.yml
CLUSTER           := swarm
DOTENV            := .env
SERVICES          := server
REGISTRY_HOST     := $(shell cat $(INFRA_DIRECTORY)/$(DOTENV) | grep "REGISTRY_HOST" | cut -d "=" -f 2)

.PHONY: path
# Show command how add Go binaries to PATH making it accessible after `go install ...`
path:
	@echo 'export PATH="$$PATH:$$(go env GOPATH)/bin"'

.PHONY: run-server
# Run cmd/server
run-server:
	@go run cmd/server/wire_gen.go cmd/server/main.go -conf=./configs -dotenv=.env.local

.PHONY: vendor
# Make ./vendor folder with dependencies
vendor:
	@go mod tidy && go mod vendor && go mod verify

.PHONY: test
# Makes go test ./...
test:
	@go test -race -parallel 10 ./...

.PHNOY: wire
# Wire dependencies with google/wire
wire:
	@go run -mod=mod github.com/google/wire/cmd/wire ./...

.PHONY: ent
# Run ent for generate schema
ent:
	@go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema

.PHONY: lint
# Run linter fo Golang files
lint:
	@docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run

.PHONY: update
# Update service in Docker Swarm without downtime
update:
	@set -e; for service in ${SERVICES}; \
		do docker pull ${REGISTRY_HOST}/${SERVICE_NAME}-$${service}:latest \
			&& docker service update \
			--with-registry-auth \
			--image ${REGISTRY_HOST}/${SERVICE_NAME}-$${service}:latest \
			${CLUSTER}_${SERVICE_NAME}-$${service} ; \
	done

.PHONY: deploy
# Deploy to Docker Swarm
deploy:
	@env \
		$$(cat ${INFRA_DIRECTORY}/${DOTENV} | sed '/^[[:blank:]]*#/d;s/#.*//' | xargs) \
		docker stack deploy \
		--orchestrator swarm \
		--with-registry-auth \
		-c "${CURRENT_DIRECTORY}"/${CONFIG_SWARM} \
		${CLUSTER}

.PHONY: undeploy
# Remove service from Docker Swarm
undeploy:
	@set -e; for service in ${SERVICES}; \
		do if docker service ls | grep -q "${CLUSTER}_${SERVICE_NAME}-$${service}" ; \
			then docker service rm ${CLUSTER}_${SERVICE_NAME}-$${service} ; \
			else echo "${CLUSTER}_${SERVICE_NAME}-$${service} is already undeployed" ; \
		fi ; \
	done

.PHONY: push
# Build and push image to registry
push:
	@set -e; for service in ${SERVICES}; \
		do docker build -t ${REGISTRY_HOST}/${SERVICE_NAME}-$${service}:latest \
			-f ${CURRENT_DIRECTORY}/Dockerfile-$${service} ${CURRENT_DIRECTORY}/. \
			&& docker push ${REGISTRY_HOST}/${SERVICE_NAME}-$${service}:latest ; \
	done

.PHONY: push
# Build and push image to registry
pull:
	@set -e; for service in ${SERVICES}; \
		do docker pull ${REGISTRY_HOST}/${SERVICE_NAME}-$${service}:latest ; \
	done

.PHONY: env
# Display environment variables from infra .env
env:
	@echo $$(cat ${INFRA_DIRECTORY}/${DOTENV} | sed '/^[[:blank:]]*#/d;s/#.*//' | xargs)

.PHONY: grafana
# Copy Grafana Dashboard file ./grafana_dashboard.json to infra Grafana dashboards directory
grafana:
	@cp -u ./grafana_dashboard.json ${INFRA_DIRECTORY}/grafana/dashboards/service-${SERVICE_NAME}.json

.PHONY: schema
# Generates Golang types from OpenAPI from api/storage/schema.yaml
schema: #
	@oapi-codegen -config ./api/storage/common/generate.yaml ./api/storage/common/schema.yaml > ./api/storage/common/schema.gen.go
	@oapi-codegen -config ./api/storage/storage/generate.yaml ./api/storage/storage/schema.yaml > ./api/storage/storage/schema.gen.go
	@oapi-codegen -config ./api/storage/generate.yaml ./api/storage/schema.yaml > ./api/storage/schema.gen.go

.PHONY: check
# Make all checks (recommended before commit and push)
check: vendor all lint test

.PHONY: combine
# Generates a combined swagger.yaml TODO make it from internal /swagger handler
combine:
	@swagger-combine ./api/storage/schema.yaml -f yaml -o static/swagger/swagger.yaml
