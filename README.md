
# Elasticsearch setup with Nomad and Consul

Default credentials for elastic search

+ username: elastic
+ password: changeme

```bash
$ nomad agent -dev
$ consul agent -dev
```


View the Consul UI at http://127.0.0.1:8500

docker run -e NOMAD_ENABLE=1 -e CONSUL_ENABLE=1 -e CONSUL_BIND_INTERFACE=eth0 -e NOMAD_ADDR=http://127.0.0.1:4646 -e NOMAD_PORT_HTTP="127.0.0.1:3000" -e LISTEN_ADDRESS=127.0.0.1:3000 -p 8000:3000 jippi/hashi-ui
