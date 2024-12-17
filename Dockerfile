FROM golang:1.23 as builder
WORKDIR /app
COPY go.mod . 
RUN go mod download
COPY . .
RUN go build -o /main ./cmd/main.go

FROM ubuntu:22.04
COPY --from=builder main /bin/main
COPY ./configs ./configs
COPY ./.env .
RUN chmod +x /bin/main

CMD ["/bin/main"]