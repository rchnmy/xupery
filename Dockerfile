FROM golang:1.22.5 AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV GOMEMLIMIT=800MiB
ENV CGO_ENABLED=0

WORKDIR /app/xupery

COPY cmd ./cmd
COPY log ./log
COPY pkg ./pkg
COPY ./go.mod ./go.sum ./

RUN go build -trimpath -ldflags="-s -w" \
    -o /app/xupery/bin ./cmd

FROM gcr.io/distroless/static-debian12:nonroot-amd64

ENV TZ=UTC

COPY --from=build /app/xupery/bin /usr/local/bin/xupery

EXPOSE 9900
ENTRYPOINT ["/usr/local/bin/xupery"]
