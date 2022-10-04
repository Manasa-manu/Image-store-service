FROM golang:1.16-alpine
COPY src /go/src
RUN cd /go/src/ && go build -o /go/bin
EXPOSE 8080
CMD [ "/go/bin/image-store-service" ]

