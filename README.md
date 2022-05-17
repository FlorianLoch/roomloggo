# Roomloggo

Small service that read temperature and humidity information from the `dnt RoomLogg Pro/ELV Raumklimastation RS 500`
via USB in order to feed it into an InfluxDB afterwards. 
Additionally, it supports exporting the measurements as Prometheus metrics.

The implementation has been inspired by/is based on [roomlogg-go](https://github.com/jhendess/roomlogg-go) and the research
done for the [Raumklima project](https://github.com/juergen-rocks/raumklima).

## Deploy to a Raspberry

The following setup assumes `roomloggo` will be deployed on the Raspberry as a Docker container.
There are several good reason to use a container image as deployment artifact, but in order for this to work some preparation is necessary:
- Install Docker
- Run the post-installation steps [[1]]
- Create the configuration file expected by `roomloggo`; it will be mounted into the container later

Once this is done build the image with `make image` and run the container as a daemon. This will configure Docker to restart the container after a failure or reboot.

```bash
docker run --restart always --privileged -v /home/floch/roomloggo.config.yaml:/app/roomloggo.config.yaml roomloggo
```


[1]: https://docs.docker.com/engine/install/linux-postinstall/