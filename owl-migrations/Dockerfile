FROM alpine:3.17.0

USER root
RUN apk update && \
  apk upgrade && \
  apk add bash && \
  rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.11.2/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose
WORKDIR /
COPY . .

CMD sh entrypoint.sh