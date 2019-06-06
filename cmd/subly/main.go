package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/urfave/cli"
	"google.golang.org/api/iterator"
)

func main() {

	app := cli.NewApp()

	app.Name = "subly"
	app.Usage = "Receive messages from Pub/Sub"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "project",
			Usage: "required - project to use",
		},
		cli.StringFlag{
			Name:  "subscription",
			Usage: "required - subscription to use",
		},
		cli.IntFlag{
			Name:  "timeout",
			Usage: "optional - timeout after seconds waiting for messages",
			Value: 10,
		},
		cli.BoolFlag{
			Name:  "create",
			Usage: "optional - create topic if it doesn't exist",
		},
		cli.BoolFlag{
			Name:  "list",
			Usage: "optional - list subscriptions for topic",
		},
		cli.StringFlag{
			Name:  "topic",
			Usage: "required if create or list are true - topic to publish to",
		},
	}

	app.Action = func(c *cli.Context) error {
		subName := c.String("subscription")
		projectID := c.String("project")
		timeout := time.Duration(c.Int("timeout")) * time.Second
		create := c.Bool("create")
		list := c.Bool("list")
		topicName := c.String("topic")

		if projectID == "" {
			cli.ShowAppHelpAndExit(c, 0)
		}

		if !list {
			if subName == "" {
				cli.ShowAppHelpAndExit(c, 0)
			}
		}

		ctx := context.Background()
		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			log.Fatal(err)
		}

		if list {
			subscriptions := client.Subscriptions(ctx)
			for {
				subscription, err := subscriptions.Next()
				if err == iterator.Done {
					return err
				}
				if err != nil {
					return err
				}
				subscriptionConfig, err := subscription.Config(ctx)
				if err != nil {
					return err
				}
				log.Printf("%v\n\t%v\n", subscription, subscriptionConfig.Topic)
			}
		}

		if create {
			if topicName == "" {
				log.Println("topic is required when creating subscription")
				cli.ShowAppHelpAndExit(c, 0)
			}

			_, err := client.CreateSubscription(ctx, subName,
				pubsub.SubscriptionConfig{
					Topic:       client.Topic(topicName),
					AckDeadline: 20 * time.Second,
				})
			if err != nil {
				log.Printf("Attempt to create subscription `%s` failed: %v", subName, err)
			}
		}

		seen := make(chan int, 1)
		cctx, cancel := context.WithCancel(ctx)
		var mu sync.Mutex
		received := 0
		sub := client.Subscription(subName)

		// Schedule cancel after n seconds
		go func() {
			for {
				select {
				case <-seen:
				case <-time.After(timeout):
					cancel()
					log.Fatalf("Timeout after %v", timeout)
				}
			}
		}()

		err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
			select {
			case seen <- 1:
			default:
			}
			msg.Ack()
			fmt.Printf("%v\n", string(msg.Data))
			mu.Lock()
			defer mu.Unlock()
			received++
			if received == 1 {
				cancel()
			}
		})

		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
