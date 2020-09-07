FROM golang:alpine AS builder
WORKDIR /app
ADD . .
RUN mkdir dist
RUN go build -o ./dist/gop .

FROM alpine
WORKDIR /app
COPY --from=builder /app/gop /bin/
ENTRYPOINT ["gop"]
