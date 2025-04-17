IMAGE_NAME=multi-nic-bw-checker
TAG=latest

build:
	go build -o bandwidth-checker ./cmd/main.go

docker-build:
	docker build -t $(IMAGE_NAME):$(TAG) .

docker-push:
	docker tag $(IMAGE_NAME):$(TAG) your-repo/$(IMAGE_NAME):$(TAG)
	docker push your-repo/$(IMAGE_NAME):$(TAG)

run:
	IPERF_MODE=client TARGET_IP=192.168.1.10 ./bandwidth-checker

