all: clean deps test build docker-build docker-deploy-up

stop: docker-deploy-down

restart: stop build docker-build docker-deploy-up

dev: clean deps build docker-dev-up run

stop-dev: docker-dev-down

tests: test cover


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


docker-dev-up:
	@cd docker-dev && docker-compose up -d

docker-dev-down:
	@cd docker-dev && docker-compose down	

run:
	@cd src/fair && TYPE_APP=DEV go run main.go 

cover:
	@go test ./... -coverprofile="go-cover-fair.tmp"

