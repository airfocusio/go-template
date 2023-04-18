FROM alpine:3.17.3
ENTRYPOINT ["/bin/go-template"]
COPY go-template /bin/go-template
WORKDIR /workdir
