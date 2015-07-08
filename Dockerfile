FROM busybox:ubuntu-14.04

WORKDIR /usr/bin
ADD ./build/short /usr/bin/short

EXPOSE 80
EXPOSE 8000
EXPOSE 8001

ENTRYPOINT ["short"]

CMD [""]
