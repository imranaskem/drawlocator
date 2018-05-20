FROM golang AS builder

WORKDIR $GOPATH/src/drawlocator

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine
COPY --from=builder /app ./
COPY --from=builder /go/src/drawlocator/frontend ./frontend
ENTRYPOINT ["./app"]

EXPOSE 80