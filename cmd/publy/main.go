package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/urfave/cli"
	"google.golang.org/api/iterator"
)

func main() {

	app := cli.NewApp()

	app.Name = "publy"
	app.Usage = "Publish messages to Pub/Sub"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "project",
			Usage: "required - project to use",
		},
		cli.StringFlag{
			Name:  "topic",
			Usage: "required - topic to publish to",
		},
		cli.StringFlag{
			Name:  "message",
			Usage: "required - body of the message to send",
		},
		cli.BoolFlag{
			Name:  "create",
			Usage: "optional - create topic if it doesn't exist",
		},
		cli.BoolFlag{
			Name:  "list",
			Usage: "optional - list topics",
		},
	}

	app.Action = func(c *cli.Context) error {
		topicName := c.String("topic")
		projectID := c.String("project")
		messageBody := c.String("message")
		list := c.Bool("list")

		if projectID == "" {
			cli.ShowAppHelpAndExit(c, 0)
		}

		if !list {
			if topicName == "" || messageBody == "" {
				cli.ShowAppHelpAndExit(c, 0)
			}
		}

		ctx := context.Background()
		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			return err
		}

		if list {
			topics := client.Topics(ctx)
			for {
				topic, err := topics.Next()
				if err == iterator.Done {
					return err
				}
				if err != nil {
					return err
				}
				log.Println(topic)
			}
		}

		if c.Bool("create") {
			_, err := client.CreateTopic(ctx, topicName)
			if err != nil {
				log.Printf("attempt to create topic `%s` failed: %v\n", topicName, err)
			}
		}

		topic := client.Topic(topicName)
		result := topic.Publish(ctx, &pubsub.Message{
			Data: []byte(messageBody),
		})
		id, err := result.Get(ctx)
		if err != nil {
			return err
		}

		log.Printf("%s\n", id)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
