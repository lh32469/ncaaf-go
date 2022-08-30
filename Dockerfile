FROM ubuntu:latest

COPY ncaaf /usr/local/bin

WORKDIR /app
COPY AP-Season.tmpl /app
COPY images         /app/images

CMD [ "ncaaf" ]

