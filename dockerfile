FROM golang:latest
WORKDIR /usr/share/jwt_app/
COPY . .
RUN go mod download && go mod verify
RUN go build main.go
CMD ["./main"]