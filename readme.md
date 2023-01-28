## PreRequirement

1. install go
2. run docker
3. install kind
    ```cmd
    go install sigs.k8s.io/kind@v0.17.0
    ```
4. create cluster
    ```cmd
   kind create cluster
    ```
5. deploy k8s dashboard
    ```cmd
   kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml
   kubectl apply -f dashboard-admin.yaml
   
    ``` 
6. login to dashboard
    ```cmd
   kubectl proxy
   open http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/.
   kubectl -n kubernetes-dashboard create token dashboard-admin
    ```
   
## build solution

```cmd
docker build -t docker.io/smartpcr/leaderelection:202301281312 . 
docker login
docker push docker.io/smartpcr/leaderelection:202301281312
```

## Run 

```cmd
kubectl config current-context
kubectl apply -f rbac.yaml
kubectl apply -f deploy.yaml
kubectl describe lease my-lease
```