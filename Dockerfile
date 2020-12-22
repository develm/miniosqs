FROM golang:1.15-alpine

WORKDIR /app/miniosqs

ENV AWS_ACCESS_KEY_ID=acess
ENV AWS_SECRET_ACCESS_KEY=access123
ENV AWS_DEFAULT_REGION=local

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/minisqs .

ENTRYPOINT ["./out/minisqs"]
