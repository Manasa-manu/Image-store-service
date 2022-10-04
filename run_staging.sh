## To run below script,ensure minikube is installed 
eval $(minikube -p minikube docker-env)
docker build --tag localhost:5000/image-store:latest .
cd build/k8s/
minikube kubectl -- apply -f mysql-deployment.yaml
minikube kubectl -- apply -f minio-dev.yaml
minikube kubectl -- apply -f image-store-app-deployment.yaml

