FROM golang:1.22.2 as builder
WORKDIR /application
COPY go.mod .
COPY go.sum .
COPY /application/publisher .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o publisher .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /application/publisher .
CMD ["./publisher"]
EXPOSE 3000