FROM alpine:3.21.1
RUN apk upgrade --update --no-cache
ENTRYPOINT ["/bin/go-template"]
COPY go-template /bin/go-template
WORKDIR /workdir
