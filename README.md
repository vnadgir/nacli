# nacli
nacli is a term from colloquial Hindi meaning "not the real thing". It also stands for NATS Command Line Interface

nacli aspires to
- provide easy way of exposing publish and subscribe bindings from the command Line
- provide commands to understand the cluster topology
- provide commands to monitor subscriber group latencies and lags

## Usage
Install nacli
  go install github.com/utilitywarehouse/nacli

Run a nats server
  docker run -it -p 4222:4222 nats

Run the subscriber
  nacli sub -b [dockerhost]:4222 -s test

Run the publisher
  nacli pub -b [dockerhost]:4222 -s test  
