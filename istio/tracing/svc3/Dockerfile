FROM golang:1.17.2 as builder
COPY . /build/
WORKDIR /build
RUN go env -w GOPROXY=https://goproxy.cn,direct && \ 
mkdir -p bin/amd64 && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64/httpserver .

# alpine:3.16.2
FROM ubuntu
COPY --from=builder /build/bin/amd64/httpserver / 
ENTRYPOINT ["/httpserver"]

# docker build -t parvinshi/httpserver:v4-svc3 .
# docker run -d -p 80:80 parvinshi/httpserver:v4 httpserver /bin/bash
# docker logs -f eb7ec007f550