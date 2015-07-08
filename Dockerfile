FROM busybox:ubuntu-14.04

WORKDIR /usr/bin
ADD ./build/short /usr/local/bin/short
ADD ./main.sh /usr/local/bin/main.sh

EXPOSE 8000
EXPOSE 8001
EXPOSE 8002

ENTRYPOINT ["/usr/local/bin/main.sh"]

CMD [""]
