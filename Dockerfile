FROM alpine:3.7

ADD release/lumber /

CMD ["/lumber"]
