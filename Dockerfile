FROM golang:1.18-alpine as builder
WORKDIR /build
ADD . .
RUN CGO_ENABLED=0 go build -o weatherapi -buildvcs=false

FROM scratch
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=0:0 --from=builder /build/* /
USER 65534
EXPOSE 8080:8080
ENTRYPOINT [ "/weatherapi" ]
