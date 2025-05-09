FROM golang:1.24.3-alpine3.21

WORKDIR /app

COPY . .

RUN go build -o mainapi 

EXPOSE 8080

CMD [ "./mainapi" ]