FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod ./
COPY *.go ./
RUN go build -o /my-image-store
EXPOSE 8080
CMD [ "/my-image-store" ]