FROM golang

WORKDIR $GOPATH/src/drawlocator

COPY . .

RUN go install

CMD ["drawlocator"]

