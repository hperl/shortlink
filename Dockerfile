FROM golang:onbuild
RUN mkdir /data
VOLUME /data
ENV STORE_FILE /data/store.txt
EXPOSE 80
