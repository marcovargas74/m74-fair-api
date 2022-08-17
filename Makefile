all: clean deps test build docker-build docker-deploy-up

stop: docker-deploy-down

restart: stop build docker-build docker-deploy-up


clean:
	@go clean
	@rm -rf build

deps:
	@go get -v ./...

test:
	@go test -v ./...


build:
	@cd src/fair && go vet && go build main.go 

docker-build:
	@docker build -t m74fairapi .

docker-deploy-up:
	@cd docker && docker-compose up -d

docker-deploy-down:
	@cd docker && docker-compose down