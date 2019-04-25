kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

sh buildImages 创建镜像

运行watch
kubectl run --rm -i watch --image=watch --image-pull-policy=Never

测试监控情况，监控Pod以dubbo开头
kubectl run dubbo --image=nginx:alpine

kubectl scale deploy/dubbo --replicas=1

kubectl scale deploy/dubbo --replicas=10
