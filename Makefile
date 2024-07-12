pull-latest-mac:
	docker pull --platform linux/x86_64 ghcr.io/isd-sgcu/rpkm67-gateway:latest
	docker pull --platform linux/x86_64 ghcr.io/isd-sgcu/rpkm67-auth:latest
	docker pull --platform linux/x86_64 ghcr.io/isd-sgcu/rpkm67-backend:latest
	docker pull --platform linux/x86_64 ghcr.io/isd-sgcu/rpkm67-checkin:latest
	docker pull --platform linux/x86_64 ghcr.io/isd-sgcu/rpkm67-store:latest

pull-latest-windows:
	docker pull ghcr.io/isd-sgcu/rpkm67-gateway:latest
	docker pull ghcr.io/isd-sgcu/rpkm67-auth:latest
	docker pull ghcr.io/isd-sgcu/rpkm67-backend:latest
	docker pull ghcr.io/isd-sgcu/rpkm67-checkin:latest
	docker pull ghcr.io/isd-sgcu/rpkm67-store:latest

docker:
	docker rm -v -f $$(docker ps -qa) || echo "No containers found. Skipping removal."
	docker-compose up

docker-qa:
	docker rm -v -f $$(docker ps -qa) || echo "No containers found. Skipping removal."
	docker-compose -f docker-compose.qa.yml up

server:
	go run cmd/main.go

watch: 
	air

mock-gen:
	mockgen -source ./internal/object/object.repository.go -destination ./mocks/object/object.repository.go
	mockgen -source ./internal/object/object.service.go -destination ./mocks/object/object.service.go
	mockgen -source ./internal/client/http/http.client.go -destination ./mocks/client/http/http.client.go
	mockgen -source ./internal/client/store/store.client.go -destination ./mocks/client/store/store.client.go
	mockgen -source ./internal/utils/random.utils.go -destination ./mocks/utils/random/random.utils.go

test:
	go vet ./...
	go test  -v -coverpkg ./internal/... -coverprofile coverage.out -covermode count ./internal/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

proto:
	go get github.com/isd-sgcu/rpkm67-go-proto@latest

swagger:
	swag init -d ./internal/file -g ../../cmd/main.go -o ./docs -md ./docs/markdown --parseDependency --parseInternal