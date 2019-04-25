
FROM alpine:latest
RUN apk --update --no-cache add tzdata  ca-certificates
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /
ADD app /
CMD ["./app","--in=true"]