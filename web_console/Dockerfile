FROM wzshiming/base:latest

COPY ./ /tmp/

RUN chmod +x /tmp/web_console

WORKDIR /tmp/

VOLUME /var/run/docker.sock

EXPOSE 8888

CMD /tmp/web_console
