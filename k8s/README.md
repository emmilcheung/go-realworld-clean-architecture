### Kubernetes setup


#### Recomendation for local development most comfortable usage:
- [minikube](https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Farm64%2Fstable%2Fbinary+download)

####  Command
```cmd
# apply
minikube kubectl apply -f .

# expose http
minikube addons enable ingress
minikube tunnel 

# kubernetes dashboard
minikube dashboard

# clean
minikube kubectl delete -f .

# more info: 
https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Farm64%2Fstable%2Fbinary+download
```

### Main app (API):

http://localhost/api/

### Jaeger UI:

http://jaeger.localhost/

### Adminer UI:

http://adminer.localhost/