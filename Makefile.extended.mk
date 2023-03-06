CURRENT_DIRECTORY := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
INFRA_DIRECTORY   := $(realpath $(CURRENT_DIRECTORY)/../infra)
SERVICE_NAME      := $(notdir $(CURRENT_DIRECTORY))
CONFIG_SWARM      := docker-compose.swarm.yml
CLUSTER           := swarm
DOTENV            := .env
SERVICES          := server
REGISTRY_HOST     := $(shell cat $(INFRA_DIRECTORY)/$(DOTENV) | grep "REGISTRY_HOST" | cut -d "=" -f 2)
STORAGE_PATH      := /storage
STAGING_SSH       := ubuntu@cargo.b-resh.ru

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
	@docker run --rm -v $$(pwd):/app \
		-e GOCACHE=/cache/go \
		-e GOLANGCI_LINT_CACHE=/cache/go \
		-v $$(go env GOCACHE):/cache/go \
		-v $$(go env GOPATH)/pkg:/go/pkg \
		-w /app golangci/golangci-lint:latest-alpine \
		golangci-lint run --verbose --timeout 5m

.PHONY: lintfix
# Run linter fo Golang files with --fix flag
lintfix:
	@docker run --rm -v $$(pwd):/app \
		-e GOCACHE=/cache/go \
		-e GOLANGCI_LINT_CACHE=/cache/go \
		-v $$(go env GOCACHE):/cache/go \
		-v $$(go env GOPATH)/pkg:/go/pkg \
		-w /app golangci/golangci-lint:latest-alpine \
		golangci-lint run --verbose --timeout 5m --fix

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

.PHONY: jwt
# Generate JWT token for access service HTTP methods
jwt:
	@set -e; if [ ! -f ./bin/jwt ] ; \
		then go build -o ./bin/jwt cmd/jwt/main.go ; \
	fi
	@./bin/jwt

.PHONY: schema
# Generates Golang types from OpenAPI from api/storage/schema.yaml
schema: #
	@oapi-codegen -config ./schema/common/generate.yaml ./schema/common/schema.yaml > ./schema/common/schema.gen.go
	@oapi-codegen -config ./schema/storage/generate.yaml ./schema/storage/schema.yaml > ./schema/storage/schema.gen.go
	@oapi-codegen -config ./schema/generate.yaml ./schema/schema.yaml > ./schema/schema.gen.go

.PHONY: check
# Make all checks (recommended before commit and push)
check: vendor all ent schema combine lint test

.PHONY: combine
# Generates a combined swagger.yaml TODO make it from internal /swagger handler
combine:
	@swagger-combine ./schema/schema.yaml -f yaml -o static/swagger/swagger.yaml

.PHONY: storage
# Make a symbolic link from /storage to current ./storage directory
storage:
	@set -e; if [ ! -d ${STORAGE_PATH} ] ; \
		then mkdir -m 0775 ${STORAGE_PATH} && chgrp docker ${STORAGE_PATH} ; \
	fi
	@set -e; if [ ! -d $$(echo ${STORAGE_PATH} | cut -d '/' -f 2) ] ; \
		then ln -s ${STORAGE_PATH} $$(pwd) ; \
	fi

.PHONY: deploy-staging
# Build current state of service and deploy to staging server
deploy-staging:
	@make push && ssh ${STAGING_SSH} 'cd /var/www/${SERVICE_NAME} && make update'
