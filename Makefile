APP := metrics

image: clean build-linux
	docker build -t "$(APP)" .

build:
	go build -o $(APP) main.go

build-linux:
	@echo "Building..."
	env GOOS=linux GOARCH=amd64  go build -o $(APP) main.go

clean:
	rm -f $(APP)

undeploy:
	kubectl delete -f deployment.yaml
	kubectl delete -f service.yaml
	kubectl delete -f grafana-dashboard-config.yaml
	kubectl delete -f prometheus.yml

deploy:
	kubectl create -f prometheus.yml
	kubectl create -f grafana-dashboard-config.yaml
	kubectl create -f service.yaml
	kubectl create -f deployment.yaml
