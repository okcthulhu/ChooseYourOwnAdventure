build:
	go build -gcflags="all=-N -l" .

run: build
	./ChooseYourOwnAdventure

test:
	go test ./... -v

docker-build:
	docker build -t my_app .

docker-push:
	docker tag my_app gcr.io/$(GCP_PROJECT_ID)/my_app
	docker push gcr.io/$(GCP_PROJECT_ID)/my_app

models:
	oapi-codegen --config=api/models.cfg.yaml cyoa.yaml

client: models
	oapi-codegen --generate client --package main cyoa.yaml > cyoa.gen.go

server:
	oapi-codegen --config=api/server.cfg.yaml cyoa.yaml

generate-all: client server


