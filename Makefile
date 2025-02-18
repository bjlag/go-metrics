BUILD_VERSION = "v1.0.0"
BUILD_DATE = $(shell date +'%Y/%m/%d %H:%M:%S')
BUILD_COMMIT = $(shell git rev-parse --short HEAD)

up:
	docker run -d --rm \
		--name postgres \
		-p 5432:5432 \
		-v ${PWD}/data/psql:/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD=local \
		-e POSTGRES_DB=metrics \
		postgres:16.4-alpine3.20

down:
	docker stop postgres

exec:
	docker exec -it postgres psql -U postgres

run:
	make -j 2 run-server run-agent

run-agent:
	go run -ldflags "-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit=$(BUILD_COMMIT)'" ./cmd/agent/.

run-server:
	go run -ldflags "-X main.buildVersion=$(BUILD_VERSION) -X 'main.buildDate=$(BUILD_DATE)' -X 'main.buildCommit=$(BUILD_COMMIT)'" ./cmd/server/.

build:
	 go build -o ./bin/multichecker ./cmd/staticlint/multichecker.go

fmt:
	goimports -local "github.com/bjlag/go-metrics" -d -w $$(find . -type f -name '*.go' -not -path "*_mock.go")
	swag fmt --dir ./cmd/server,./internal/handler

lint:
	$(if $(wildcard ./bin/multichecker),,$(error "Binary './bin/multichecker' not found. Please run 'make build'"))
	./bin/multichecker -c 2 ./...

doc:
	godoc -http=:8888 -play

swagger:
	swag init --parseDependency --parseDepth 1 --dir ./cmd/server,./internal/handler

test:
	go test ./...

cover:
	go test -coverpkg='./internal/...','./cmd' -coverprofile coverage.out.tmp ./... \
    	&& cat coverage.out.tmp | egrep -v "_mock.go" > coverage.out \
    	&& rm coverage.out.tmp \
    	&& go tool cover -func coverage.out

gen-proto:
	protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import --go-grpc_opt=require_unimplemented_servers=false proto/metric.proto
