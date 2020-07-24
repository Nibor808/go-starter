FROM golang:latest
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go-starter

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o app .

FROM alpine:latest
WORKDIR /root/

COPY --from=0 /go-starter/app .
COPY --from=0 /go-starter/.env .

RUN mkdir -p ./view
COPY --from=0 /go-starter/view/index.html ./view

CMD ["./app"]
