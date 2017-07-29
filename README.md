
# Elasticsearch setup with Nomad and Consul

Default credentials for elastic search

+ username: elastic
+ password: changeme

```bash
$ nomad agent -dev
$ consul agent -dev
```


View the Consul UI at http://127.0.0.1:8500

```bash
$ docker run -e NOMAD_ENABLE=1 -e CONSUL_ENABLE=1 -e CONSUL_BIND_INTERFACE=eth0 -e NOMAD_ADDR=http://127.0.0.1:4646 -e NOMAD_PORT_HTTP="192.168.8.100:3000" -e LISTEN_ADDRESS=192.168.8.100:3000 -p 8000:3000 jippi/hashi-ui:v0.13.6


# Does not work
$ docker run -e NOMAD_ENABLE=1 -e NOMAD_ADDR=http://nomad.service.consul:4646 -p 8000:3000 jippi/hashi-ui:v0.13.6

# Does not work
$ docker run --rm -e NOMAD_PROXY_ADDRESS="localhost:3000/nomad" -e NOMAD_ENABLE=true -e NOMAD_ADDR=http://nomad.service.consul:4646 -p 3000:3000 jippi/hashi-ui:v0.13.6


# Does not work
$ docker run -e NOMAD_ENABLE=1 -p 8000:3000 jippi/hashi-ui

# Does not work
$ docker run --rm -e NOMAD_ENABLE=1 -e NOMAD_ADDR=http://192.168.8.100:4646 -p 3000:3000 jippi/hashi-ui:v0.13.6

# Works
$ docker run --net=host --rm -e NOMAD_ENABLE=1 -e NOMAD_ADDR=http://docker.for.mac.localhost:4646 -p 3000:3000 jippi/hashi-ui:v0.13.6
```



# Hashi UI

If you are planning to use hashi-ui as artifacts:

https://github.com/jippi/hashi-ui/releases/download/v0.13.6/hashi-ui-darwin-amd64