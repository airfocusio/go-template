FROM alpine:3.22.0
RUN apk upgrade --update --no-cache
ENTRYPOINT ["/bin/go-template"]
COPY go-template /bin/go-template
WORKDIR /workdir
