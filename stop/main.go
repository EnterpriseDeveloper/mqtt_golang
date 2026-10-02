package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Start отправляет команду START в CODESYS
func Stop() {

	const broker = "tcp://localhost:1883"
	const topic = "PLC/command"

	// Настройки MQTT клиента
	opts := mqtt.NewClientOptions()

	opts.AddBroker(broker)
	opts.SetClientID("GO_STOP_COMMAND")
	opts.SetCleanSession(true)

	// Создаём MQTT client
	client := mqtt.NewClient(opts)

	// Подключаемся к broker
	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connection error: %v", token.Error())
	}

	fmt.Println("Connected to MQTT broker")

	// Отправляем STOP
	token = client.Publish(
		topic,
		1,      // QoS
		false,  // Retain
		"STOP", // Payload
	)

	if token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT publish error: %v", token.Error())
	}

	fmt.Println("=================================")
	fmt.Println("STOP command sent")
	fmt.Println("Topic   :", topic)
	fmt.Println("Payload : STOP")
	fmt.Println("=================================")

	// Даём MQTT время завершить отправку
	time.Sleep(100 * time.Millisecond)

	client.Disconnect(250)
}

func main() {
	Stop()
}
