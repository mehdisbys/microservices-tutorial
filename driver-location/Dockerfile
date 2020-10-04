FROM golang:1.15.2-alpine3.12

ADD ./driver-location/config.yaml config.yaml
ADD ./driver-location/main .

EXPOSE 80

CMD ["./main"]
