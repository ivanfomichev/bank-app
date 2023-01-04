FROM golang:1.19 AS builder

RUN apt update && \
  apt install -y git

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  GOPROXY=$GOPROXY

WORKDIR /build
COPY . .

ARG LDFLAGS
RUN go build -ldflags="${LDFLAGS}" -o /go/bin/bank-app


FROM scratch

COPY --from=builder /go/bin/bank-app /go/bin/bank-app
COPY config.* /

EXPOSE 58001

CMD ["/go/bin/bank-app"]
