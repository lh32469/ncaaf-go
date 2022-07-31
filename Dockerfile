FROM ubuntu:latest

COPY ncaaf /usr/local/bin

WORKDIR /app
COPY AP-Season.tmpl /app

CMD [ "ncaaf" ]

