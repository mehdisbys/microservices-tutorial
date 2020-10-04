FROM golang:1.15.2-alpine3.12

ADD ./zombie-driver/config.yaml config.yaml
ADD ./zombie-driver/main .

EXPOSE 80

CMD ["./main"]