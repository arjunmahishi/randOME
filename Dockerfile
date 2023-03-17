FROM alpine:3.14

COPY ./bin/randOME_linux /app/randOME
COPY ./sample.yaml /app/sample.yaml

CMD ["/app/randOME", "emit", "-c", "/app/sample.yaml"]
