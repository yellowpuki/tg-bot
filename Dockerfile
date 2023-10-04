FROM golang:alpine as builder

ARG TOKEN

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"'

FROM scratch
COPY --from=builder /app /app
EXPOSE 4000

RUN app $TOKEN