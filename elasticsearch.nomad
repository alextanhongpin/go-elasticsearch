
job "search" {
  datacenters = ["dc1"]

  # region = "us"

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
    task "hashi-ui" {
      driver = "docker"
      config {
        image = "jippi/hashi-ui:latest"
        port_map {
          http = 3000
        }
      }
      env {
        NOMAD_ENABLE = 1
        NOMAD_ADDR = "http://docker.for.mac.localhost:4646"
        // NOMAD_LOG_LEVEL = "debug"
        // NOMAD_READ_ONLY = "false"
        // PROXY_ADDRESS   = "/hashi-ui"

        // CONSUL_ENABLE    = "true"
        // CONSUL_ADDR      = "169.254.1.1:8500"
        // CONSUL_READ_ONLY = "true"
      }
      resources {
        # 1000 MHz
        cpu = 1000
        # 512 MB
        disk = 512
        # 1024 MB
        memory = 1024
        network {
          mbits = 10
          port "http" {}
        }
      }
    }
    // task "hashi-ui" {
    //   driver = "raw_exec"
    //   # TODO: Find out where nomad download the artifact
    //   // artifact {
    //   //   source      = "https://github.com/jippi/hashi-ui/releases/download/v0.13.6/hashi-ui-darwin-amd64"
    //   //   destination = "./local/dir"
    //   // }
    //   config {
    //     # Path must be absolute
    //     command = "/Users/alextanhongpin/Documents/golang/src/github.com/alextanhongpin/go-elasticsearch/hashi-ui-darwin-amd64"
    //     args = ["--nomad-enable", "--consul-enable"]
    //   }
    //   resources {
    //     # 1000 MHz
    //     cpu = 1000
    //     # 512 MB
    //     disk = 512
    //     # 1024 MB
    //     memory = 1024
    //     network {
    //       mbits = 10
    //       port "http" {}
    //     }
    //   }
    // }
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
        # 1000 MHz
        cpu = 1000
        # 512 MB
        disk = 512
        # 1024 MB
        memory = 1024
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