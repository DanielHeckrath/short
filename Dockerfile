FROM busybox:ubuntu-14.04

WORKDIR /usr/bin

EXPOSE 80
EXPOSE 8000
EXPOSE 8001

ADD ./build/short /usr/bin/short

ENTRYPOINT ["short"]

CMD [""]
