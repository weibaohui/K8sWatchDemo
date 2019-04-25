kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

kubectl run --rm -i watch --image=watch --image-pull-policy=Never

kubectl run dubbo --image=nginx:alpine