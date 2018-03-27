FROM alpine:3.7

ENV PORT 80
EXPOSE $PORT

ADD release/lumber /

CMD ["/lumber"]
