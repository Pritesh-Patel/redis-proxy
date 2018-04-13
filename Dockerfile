FROM alpine

WORKDIR /app
COPY bin/app /app/

CMD ["./app"]