FROM golang:1.16

WORKDIR /go/src

RUN git clone https://github.com/dbogunia/anonymizer

WORKDIR /go/src/anonymizer

CMD sh ./start.sh "username:password@protocol(address)/dbname"
