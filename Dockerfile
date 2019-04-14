FROM golang AS builder
WORKDIR $HOME/etc/drawlocator
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine
COPY --from=builder /app ./
ENTRYPOINT ["./app"]