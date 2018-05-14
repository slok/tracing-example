
FROM golang:1.10-alpine AS build-stage

WORKDIR /go/src/github.com/slok/tracing-example
COPY . .
RUN ./hack/scripts/build.sh

# Final image.
FROM alpine:latest
RUN apk --no-cache add \
  ca-certificates
COPY --from=build-stage /go/src/github.com/slok/tracing-example/bin/tracing-example /usr/local/bin/tracing-example
ENTRYPOINT ["/usr/local/bin/tracing-example"]
CMD ["--help"]