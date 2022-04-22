test:
	go test ./... -v

docker-build:
	docker build -t brunodeluk/kube-config:latest .

docker-push:
	docker push brunodeluk/kube-config:latest