APP?=wbl0

GOOS?=linux
GOARCH?=amd64

clean:
	rm -f ${APP}

build:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o $(APP) cmd/main.go

test:
	go test -v ./...

run: build
	./${APP}

stop_compose:
	docker compose down

run_compose: stop_compose
	docker compose up -d --force-recreate --build