FROM golang:1.17.2 as builder
COPY . /build/
WORKDIR /build
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN mkdir -p bin/amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64/httpserver .
 
FROM scratch
COPY --from=builder /build/bin/amd64/httpserver / 
ENTRYPOINT ["/httpserver"]

# docker build -t="parvinshi/httpserver:v1" .
# docker run -d -p 80:80 parvinshi/httpserver:v1 httpserver /bin/bash
# docker logs parvinshi/httpserver