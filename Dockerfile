FROM golang:alpine

COPY . /app

WORKDIR /app

RUN go install github.com/trap-bytes/gourlex@latest

ENTRYPOINT ["gourlex"]
