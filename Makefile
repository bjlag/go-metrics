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
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out