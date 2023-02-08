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

# docker build -t parvinshi/httpserver:v2 .
# docker run -d -p 8081:8081 parvinshi/httpserver:v2 httpserver /bin/bash
# docker logs -f eb7ec007f550