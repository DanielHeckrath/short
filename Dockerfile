FROM busybox:ubuntu-14.04

WORKDIR /usr/bin
ADD ./build/short /usr/bin/short

EXPOSE 8000
EXPOSE 8001
EXPOSE 8002

ENTRYPOINT ["short"]

CMD [""]
