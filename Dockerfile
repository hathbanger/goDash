# command to generate binary
# docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.8.3 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o goDash'

FROM alpine:latest 
WORKDIR /app
COPY goDash /app/
EXPOSE 1323
ENTRYPOINT ["./goDash"]

