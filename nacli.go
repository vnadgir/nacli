package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jawher/mow.cli"
	"github.com/nats-io/go-nats"
)

func main() {
	app := cli.App("nacli", "command line interface to work with nats")
	setupSubscriber(app)
	setupPublisher(app)
	app.Run(os.Args)
}

func setupSubscriber(app *cli.Cli) {
	app.Command("sub", "subscribe to a topic", func(cmd *cli.Cmd) {
		app.Spec = "--subject, --brokers"
		subject := cmd.String(cli.StringOpt{
			Name: "subject, s",
			Desc: "Subject to subscribe to",
		})

		brokers := cmd.Strings(cli.StringsOpt{
			Name: "brokers, b",
			Desc: "Brokers to connect to",
		})
		cmd.Action = func() {
			subscribe(*brokers, *subject)
		}
	})
}

func setupPublisher(app *cli.Cli) {
	app.Command("pub", "publish to a topic", func(cmd *cli.Cmd) {
		app.Spec = "--subject, --brokers"
		subject := cmd.String(cli.StringOpt{
			Name: "subject, s",
			Desc: "Subject to subscribe to",
		})

		brokers := cmd.Strings(cli.StringsOpt{
			Name: "brokers, b",
			Desc: "Brokers to connect to",
		})
		cmd.Action = func() {
			publish(*brokers, *subject)
		}
	})
}

func publish(brokers []string, subject string) {
	brokerList := strings.Join(brokers, ",")
	log.Printf("Connecting to %v\n", brokerList)
	natsConnection, err := nats.Connect("nats://" + brokerList)
	if err != nil {
		log.Fatalf("Unable to connect. %v", err)
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		natsConnection.Publish(subject, []byte(s.Text()))
	}
}

func subscribe(brokers []string, subject string) {
	brokerList := strings.Join(brokers, ",")
	log.Printf("Connecting to %v\n", brokerList)
	natsConnection, err := nats.Connect("nats://" + brokerList)
	if err != nil {
		log.Fatalf("Unable to connect. %v", err)
	}

	log.Printf("Subscribing to subject '%v'\n", subject)
	natsConnection.QueueSubscribe(subject, uuid.New().String(), func(msg *nats.Msg) {
		log.Printf("%s\n", string(msg.Data))
	})
}
