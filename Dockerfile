FROM alpine:3.18.0
ENTRYPOINT ["/bin/go-template"]
COPY go-template /bin/go-template
WORKDIR /workdir
