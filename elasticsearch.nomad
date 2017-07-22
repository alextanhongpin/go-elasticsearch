
job "search" {
  datacenters = ["dc1"]

  type = "service"

  update {
    stagger = "10s"
    max_parallel = 1
  }

  group "cache" {
    count = 1
    restart {
      attempts = 10
      interval = "5m"
      delay = "25s"
      mode = "delay"
    }
    task "elasticsearch" {
      driver = "docker"
      config {
        image = "docker.elastic.co/elasticsearch/elasticsearch:5.4.3"
        port_map {
          http = 9200
          tcp = 9300
        }
      }

      env {
        bootstrap.memory_lock = "true"
        ES_JAVA_OPTS = "-Xms512m -Xmx512m"
        cluster.name = "docker-cluster"
      }

      resources {
        cpu    = 1000 # 1000 MHz
        memory = 1024 # 1024MB
        disk = 512
        network {
          mbits = 10
          port "http" {}
          port "tcp" {}
        }
      }

      service {
        name = "global-elastichsearch-check"
        tags = ["global", "search"]
        port = "http"
        check {
          name     = "alive"
          type     = "tcp"
          interval = "10s"
          timeout  = "2s"
        }
      }
// template {
//         data = <<EOH
// cluster.name: "escluster"
// network.host: {{ env "attr.unique.network.ip-address" }}
// discovery.zen.minimum_master_nodes: 2
// network.publish_host: {{ env "attr.unique.network.ip-address" }}
// {{ if service "escluster-transport"}}discovery.zen.ping.unicast.hosts:{{ range service "escluster-transport" }}
//   - {{ .Address }}:{{ .Port }}{{ end }}{{ end }}
// http.port: {{ env "NOMAD_HOST_PORT_http" }}
// http.publish_port: {{ env "NOMAD_HOST_PORT_http" }}
// transport.tcp.port: {{ env "NOMAD_HOST_PORT_transport" }}
// transport.publish_port: {{ env "NOMAD_HOST_PORT_transport" }}

// action.auto_create_index: filebeat*

// readonlyrest:
//   enable: false
// EOH
//         destination = "local/elasticsearch.yml"
//         change_mode = "noop"
//       }

    }
  }
}