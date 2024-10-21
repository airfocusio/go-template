FROM alpine:3.20.3
RUN apk upgrade --update --no-cache
ENTRYPOINT ["/bin/go-template"]
COPY go-template /bin/go-template
WORKDIR /workdir
