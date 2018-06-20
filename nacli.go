package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/jawher/mow.cli"
	"github.com/nats-io/go-nats"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/go-nats-streaming/pb"
)

func main() {
	app := cli.App("nacli", "command line interface to work with nats")

	isStan := app.Bool(cli.BoolOpt{
		Name:  "streaming s",
		Desc:  "Connects to streaing server if true",
		Value: true,
	})

	setupSubscriber(app, isStan)
	setupPublisher(app, isStan)
	setupctl(app)
	app.Run(os.Args)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, closing connection...\n\n")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func setupctl(app *cli.Cli) {
	app.Command("ctl", "control plane for nats", func(cmd *cli.Cmd) {
		cmd.Command("routes", "routes on nats", func(cmd *cli.Cmd) {
			cmd.Action = func() {

				log.Printf("Listing route info\n")
			}
		})

		cmd.Command("subs", "subscriptions on nats", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				log.Printf("Listing subscriptions info\n")
			}
		})

	})
}

func setupSubscriber(app *cli.Cli, isStan *bool) {
	app.Command("sub", "subscribe to a topic", func(cmd *cli.Cmd) {
		cmd.Spec = "--topic --brokerURL --clusterID --from-beginning --persistent"
		subject := cmd.String(cli.StringOpt{
			Name: "topic t",
			Desc: "Subject to subscribe to",
		})

		brokers := cmd.String(cli.StringOpt{
			Name: "brokerURL b",
			Desc: "Brokers to connect to",
		})

		clusterID := cmd.String(cli.StringOpt{
			Name: "clusterID c",
			Desc: "clusterID to connect to",
		})

		fromBeginning := cmd.Bool(cli.BoolOpt{
			Name:  "from-beginning f",
			Desc:  "subscribe from the beginning of time",
			Value: true,
		})

		persistent := cmd.String(cli.StringOpt{
			Name: "persistent p",
			Desc: "Keep the subscription persistent and dont kill the subscription",
		})
		cmd.Action = func() {
			subscribe(*brokers, *subject, *clusterID, *isStan, *fromBeginning, persistent)
		}
	})
}

func setupPublisher(app *cli.Cli, isStan *bool) {
	app.Command("pub", "publish to a topic", func(cmd *cli.Cmd) {
		cmd.Spec = "--topic --brokerURL --clusterID"
		subject := cmd.String(cli.StringOpt{
			Name: "topic t",
			Desc: "Subject to subscribe to",
		})

		brokerURL := cmd.String(cli.StringOpt{
			Name: "brokerURL b",
			Desc: "Brokers to connect to",
		})

		clusterID := cmd.String(cli.StringOpt{
			Name: "clusterID c",
			Desc: "clusterID to connect to",
		})
		cmd.Action = func() {
			publish(*brokerURL, *subject, *clusterID, *isStan)
		}
	})
}

func publish(brokerURL string, subject string, clusterID string, isStan bool) {
	log.Printf("Connecting to %v\n", brokerURL)
	if !isStan {
		natsConnection, err := nats.Connect(brokerURL)
		if err != nil {
			log.Fatalf("Unable to connect. %v", err)
		}
		defer natsConnection.Close()

		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			natsConnection.Publish(subject, []byte(s.Text()))
		}

	} else {
		natsConnection, err := stan.Connect(clusterID, uuid.New().String(), stan.NatsURL(brokerURL))
		if err != nil {
			log.Fatalf("Unable to connect. %v", err)
		}
		defer natsConnection.Close()

		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			natsConnection.Publish(subject, []byte(s.Text()))
		}
	}
}

func subscribe(brokerURL string, subject string, clusterID string, isStan bool, fromBeginning bool, persistentName *string) {
	log.Printf("Connecting to %v\n", brokerURL)

	if !isStan {
		natsConnection, err := nats.Connect(brokerURL)
		if err != nil {
			log.Fatalf("Unable to connect. %v", err)
		}

		//defer natsConnection.Close()

		log.Printf("Subscribing to subject '%v'\n", subject)

		_, err = natsConnection.QueueSubscribe(subject, uuid.New().String(), func(msg *nats.Msg) {
			log.Printf("%s\n", string(msg.Data))
		})
		if err != nil {
			panic(err)
		}
		//defer sub.Unsubscribe()

	} else {
		natsConnection, err := stan.Connect(clusterID, uuid.New().String(), stan.NatsURL(brokerURL))
		if err != nil {
			log.Fatalf("Unable to connect. %v", err)
		}

		//defer natsConnection.Close()
		var startPos pb.StartPosition
		if fromBeginning {
			startPos = pb.StartPosition_First
		} else {
			startPos = pb.StartPosition_LastReceived
		}

		_, err = natsConnection.QueueSubscribe(subject, uuid.New().String(), func(msg *stan.Msg) {
			log.Printf("%s\n", msg.String())
		}, stan.StartAt(startPos))
		if err != nil {
			panic(err)
		}
		//defer sub.Unsubscribe()
	}

}
