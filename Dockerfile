FROM alpine:3.14

COPY ./bin/randOME_linux /app/randOME
COPY ./scripts/run.sh /app/run.sh
COPY ./sample.yaml /app/sample.yaml

RUN chmod +x /app/run.sh

CMD ["/app/run.sh"]
