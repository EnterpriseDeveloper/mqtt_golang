package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {

	const broker = "tcp://127.0.0.1:1883"
	const topic = "PLC/command"

	fmt.Println("========================================")
	fmt.Println("MQTT START COMMAND")
	fmt.Println("========================================")
	fmt.Println("Connecting to:", broker)

	opts := mqtt.NewClientOptions()

	opts.AddBroker(broker)
	opts.SetClientID("GO_START_COMMAND")
	opts.SetCleanSession(true)

	client := mqtt.NewClient(opts)

	// Connect to Mochi broker
	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		log.Fatalf("Connection failed: %v", token.Error())
	}

	fmt.Println("Connected!")

	// Send START
	token = client.Publish(
		topic,
		1,
		false,
		"START",
	)

	if token.Wait() && token.Error() != nil {
		log.Fatalf("Publish failed: %v", token.Error())
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("COMMAND SENT")
	fmt.Println("========================================")
	fmt.Println("Topic   :", topic)
	fmt.Println("Payload :", "START")
	fmt.Println("========================================")

	// Даем MQTT завершить отправку
	time.Sleep(100 * time.Millisecond)

	client.Disconnect(250)

	fmt.Println("Disconnected")
}
