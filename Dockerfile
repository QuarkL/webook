# basic image
FROM ubuntu:24.04
# copy project to /app
COPY webook /app/webook
WORKDIR /app
CMD ["/app/webook"]
