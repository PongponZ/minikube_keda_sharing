start:
	minikube start

stop:
	minikube stop

clear-image:
	minikube image rm minikube-publisher:1.0.1
	minikube image rm minikube-consumer:1.0.1

deploy:
	kubectl apply -f deployment/publisher-deployment.yml
	kubectl apply -f deployment/consumer-deployment.yml

app-to-kube:
	eval $(minikube docker-env)
	docker build -t minikube-publisher:1.0.1 -f application/publisher/Dockerfile .
	docker build -t minikube-consumer:1.0.1 -f application/consumer/Dockerfile .
	minikube image rm minikube-publisher:1.0.1
	minikube image rm minikube-consumer:1.0.1
	minikube image load minikube-publisher:1.0.1
	minikube image load minikube-consumer:1.0.1

app-url:
	minikube service publisher-service --url