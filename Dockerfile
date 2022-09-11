FROM ubuntu:latest

RUN apt update && apt install tzdata -y
ENV TZ="America/Los_Angeles"

COPY ncaaf /usr/local/bin

WORKDIR /app
COPY AP-Season.tmpl /app
COPY images         /app/images

CMD [ "ncaaf" ]

