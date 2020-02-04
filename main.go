package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

var (
	subscriptionNameInput *string
	topicNameInput        *string
)

func init() {
	topicNameInput = flag.String("topicName", "", "The topic name")
	subscriptionNameInput = flag.String("subscriptionName", "", "The subscription name")
}

func main() {
	flag.Parse()

	if *topicNameInput == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *subscriptionNameInput == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}

	// Create a client to communicate with a Service Bus Namespace.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	// build the topic for sending priority messages
	tm := ns.NewTopicManager()

	te, err := tm.Get(ctx, *topicNameInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	sm, err := ns.NewSubscriptionManager(te.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	se, err := sm.Get(ctx, *subscriptionNameInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *se.SubscriptionDescription.CountDetails.DeadLetterMessageCount == 0 {
		fmt.Println("No dead letters found")
		return
	}

	fmt.Printf("Found %d deadletters\n", *se.SubscriptionDescription.CountDetails.DeadLetterMessageCount)

	// build the topic for sending priority messages
	topic, err := ns.NewTopic(*topicNameInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	subscription, err := topic.NewSubscription(*subscriptionNameInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	qdl := subscription.NewDeadLetter()
	fmt.Println("Looking for deadletters")
	if err := qdl.ReceiveOne(ctx, servicebus.HandlerFunc(func(ctx context.Context, message *servicebus.Message) error {
		fmt.Println("Found dead letter:")
		fmt.Println(string(message.Data))
		return message.Complete(ctx)
	})); err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = qdl.Close(ctx)
	}()
}
