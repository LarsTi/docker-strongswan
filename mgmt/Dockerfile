FROM golang:latest as builder
WORKDIR /app
RUN go get github.com/tatsushid/go-fastping && \
    go get github.com/prometheus/client_golang/prometheus && \
    go get github.com/prometheus/client_golang/prometheus/promauto && \
    go get github.com/prometheus/client_golang/prometheus/promhttp && \
    go get github.com/strongswan/govici/vici
COPY app/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM alpine:latest
WORKDIR /app
Copy healthcheck.sh /app/
COPY --from=builder /app .
CMD ["/app/main"]
HEALTHCHECK --start-period=60s CMD /bin/sh healthcheck.sh
