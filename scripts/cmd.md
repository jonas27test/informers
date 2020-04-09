#### get kube config with token
kubectl config view --raw

#### forward dashboard
microk8s kubectl port-forward -n kube-system service/kubernetes-dashboard 10443:443 &

#### Get all Resource APIs
kubectl api-resources -o wide