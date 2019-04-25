CGO_ENABLED=0 GOOS=linux go build -ldflags '-d -w -s ' -a -installsuffix cgo -o app .
docker build -t watch .