FROM golang:1.15.6-alpine3.12
RUN apk add git

COPY . /home/src
WORKDIR /home/src
RUN go build -o /bin/action ./

ENTRYPOINT [ "/bin/action" ]