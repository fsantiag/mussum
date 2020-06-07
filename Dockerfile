FROM debian:stable-slim

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates

RUN mkdir /bot

COPY mussum /bot

WORKDIR /bot

CMD ["/bot/mussum"]
