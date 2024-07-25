install-minikube:
	curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-arm64
	sudo install minikube-darwin-arm64 /usr/local/bin/minikube

start:
	minikube start

stop:
	minikube stop

setup-rabbitmq:
	brew install helm
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm install rabbitmq-local bitnami/rabbitmq --version 14.4.4 --set auth.username=rabbitmq,auth.password=rabbitmqpw

setup-keda:
	kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.14.0/keda-2.14.0.yaml

build-image:
	eval $(minikube docker-env)
	docker build -t minikube-publisher:1.0.2 -f application/publisher/Dockerfile .
	docker build -t minikube-consumer:1.0.2 -f application/consumer/Dockerfile .
	minikube image rm minikube-publisher:1.0.2
	minikube image rm minikube-consumer:1.0.2
	minikube image load minikube-publisher:1.0.2
	minikube image load minikube-consumer:1.0.2

deploy-app:
	kubectl apply -f deployment/publisher-deployment.yml
	kubectl apply -f deployment/consumer-deployment.yml

deploy-keda:
	kubectl apply -f deployment/consumer-keda.yml

get-keda:
	kubectl get scaledobjects.keda.sh

delete-keda:
	kubectl delete scaledobjects.keda.sh consumer-keda

publisher-url:
	minikube service publisher-svc --url

rabbitmq-forward-port:
	kubectl port-forward --namespace default service/rabbitmq-local 15672:15672 5672:5672

rabbitmq-auth-info:
	kubectl exec -it pod/rabbitmq-local-0 -- env | grep 'RABBITMQ_USERNAME\|RABBITMQ_PASSWORD'
