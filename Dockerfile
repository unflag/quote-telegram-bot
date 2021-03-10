FROM alpine:3.13

MAINTAINER Vyacheslav Mitrofanov <unflag@ymail.com>

COPY quote-telegraf-bot /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/quote-telegraf-bot"]
