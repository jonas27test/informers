FROM golang:1.13.9 as builder
EXPOSE 8080
COPY . .
RUN cp -r conf /conf
RUN useradd scratchuser && \
    export GOPATH="" && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /informer .

FROM scratch
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
COPY --from=builder /informer /informer
COPY --from=builder /conf /conf
COPY --from=builder /etc/passwd /etc/passwd
USER scratchuser
CMD ["/informer"]
