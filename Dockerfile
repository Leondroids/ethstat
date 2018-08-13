FROM alpine:3.4
MAINTAINER 	Leonard Schellenberg <leonard.schellenberg@gmail.com>

RUN apk update
RUN apk add ca-certificates
RUN apk add curl
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 8060

CMD ["index"]
HEALTHCHECK --interval=10s CMD wget -qO- localhost:8060/healthcheck


COPY static  /usr/local/bin/static
COPY index /usr/local/bin/index

WORKDIR /usr/local/bin

RUN chmod +x /usr/local/bin/index
