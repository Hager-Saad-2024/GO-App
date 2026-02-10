FROM golang:1.21-alpine 
WORKDIR /app
COPY main.go .
COPY go.sum .
COPY go.mod .
RUN go mod download && \
    go mod verify
RUN go build -o survey
CMD [ "./survey" ]
EXPOSE 8080