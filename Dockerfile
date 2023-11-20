FROM alpine:3.18.4
RUN apk update --upgrade --no-cache
ENTRYPOINT ["/bin/go-template"]
COPY go-template /bin/go-template
WORKDIR /workdir
