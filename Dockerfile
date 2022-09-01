FROM docker.io/library/golang:1.18.5 AS builder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/cache/apt/* && \
    useradd --home /app/ -M appuser && \
    update-ca-certificates

WORKDIR /src/
COPY . .
RUN go build -ldflags "-extldflags \"-static\"" -a -installsuffix cgo -o /app/portainer-git-redeploy .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

WORKDIR /app/
COPY --from=builder /app/* ./
USER appuser
ENTRYPOINT ["/app/portainer-git-redeploy"]
