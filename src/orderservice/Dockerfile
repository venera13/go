FROM debian:9-slim

ADD ./cmd/orderservice/orderservice /app/bin/
WORKDIR /app

EXPOSE 8000

CMD [ "/app/bin/orderservice" ]