FROM ubuntu:latest

RUN apt update && apt install tzdata ca-certificates openssl -y
ENV TZ="America/Los_Angeles"

COPY ncaaf /usr/local/bin

WORKDIR /app
COPY *.tmpl /app
COPY images         /app/images

CMD [ "ncaaf" ]

