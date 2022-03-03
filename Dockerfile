FROM alpine:3.15
ENTRYPOINT ["/bin/go-template"]
COPY go-template /bin/go-template
WORKDIR /workdir
