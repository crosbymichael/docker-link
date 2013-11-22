## Docker engine

A simple PoC to link into docker itself.

To create a container named `engine` just run:

```bash
docker run -d -v /var/run/:/docker -name engine crosbymichael/engine
```

You need to bind mount the docker socket into the engine container then anything linked into 
the container can access the docker remote api on port 4243.
