FROM golang AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o app 

EXPOSE 4000

CMD [ "./app" ]