start:
	minikube start

stop:
	minikube stop

deploy:
	kubectl apply -f deployment/publisher-deployment.yml
	kubectl apply -f deployment/consumer-deployment.yml

setup-keda:
	kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.14.0/keda-2.14.0.yaml

setup-rabbitmq:
	brew install helm
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm install rabbitmq-local bitnami/rabbitmq --version 14.4.4

clear-image:
	minikube image rm minikube-publisher:1.0.2
	minikube image rm minikube-consumer:1.0.2

app-url:
	minikube service publisher-svc --url

rabbitmq-dashboard:
	kubectl port-forward --namespace default service/rabbitmq-local 15672:15672

deploy-keda:
	kubectl apply -f deployment/consumer-keda.yml

deploy-app:
	eval $(minikube docker-env)
	docker build -t minikube-publisher:1.0.2 -f application/publisher/Dockerfile .
	docker build -t minikube-consumer:1.0.2 -f application/consumer/Dockerfile .
	minikube image rm minikube-publisher:1.0.2
	minikube image rm minikube-consumer:1.0.2
	minikube image load minikube-publisher:1.0.2
	minikube image load minikube-consumer:1.0.2

app-url:
	minikube service publisher-svc --url
