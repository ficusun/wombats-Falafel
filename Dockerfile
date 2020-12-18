FROM golang:1.15.6-alpine3.12 as builder
RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN go build -o server .

#scratch  alpine:latest
FROM golang:1.15.6-alpine3.12
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
