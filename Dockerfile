FROM busybox

ADD engine /engine

EXPOSE 4243
ENTRYPOINT ["/engine", "-s", "/docker/docker.sock", "-p", "4243"]
