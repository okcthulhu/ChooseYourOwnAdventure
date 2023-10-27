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
