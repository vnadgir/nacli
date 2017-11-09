# nacli
`nacli` is a term from colloquial Hindi meaning "not the real thing". It also stands for NATS Command Line Interface. Suffice to say this is not the

`nacli` aspires to
- provide easy way of exposing publish and subscribe bindings from the command Line
- provide commands to understand the cluster topology
- provide commands to monitor subscriber group latencies and lags

## Usage
Install nacli

    go install github.com/vnadgir/nacli

Run a nats server/cluster

- Single Server

      docker run -it -p 4222:4222 nats

- Cluster Mode

      docker run -d -p 14222:4222 -p 15222:15222 -p 18222:8222 -p 16222:6222 nats --cluster nats://0.0.0.0:15222

      docker run -d -p 24222:4222 -p 25222:25222 -p 28222:8222 -p 26222:6222 nats --cluster nats://0.0.0.0:25222 --routes nats://0.0.0.0:15222

Run the subscriber

    nacli sub -b [dockerhost]:4222 -s test

Run the publisher

    nacli pub -b [dockerhost]:4222 -s test  
