FROM golang:alpine as builder
WORKDIR /go/src/K8sWatchDemo/
COPY . .
RUN ls
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-d -w -s ' -a -installsuffix cgo -o app .
RUN ls

FROM busybox
WORKDIR /app/
COPY --from=builder /go/src/K8sWatchDemo/app .

CMD ["./app","--in=true"]