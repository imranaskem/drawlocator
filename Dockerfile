FROM golang AS builder
WORKDIR $HOME/etc/drawlocator
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine
COPY --from=builder /app ./
ENTRYPOINT ["./app"]