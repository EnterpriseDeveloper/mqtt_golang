package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

// ----------------------------------------------------
// MQTT Debug Hook
// ----------------------------------------------------

type DebugHook struct {
	mqtt.HookBase
}

// ID нашего hook
func (h *DebugHook) ID() string {
	return "debug-console"
}

// Говорим Mochi, какие события нас интересуют
func (h *DebugHook) Provides(b byte) bool {
	return b == mqtt.OnPacketRead ||
		b == mqtt.OnPacketSent ||
		b == mqtt.OnConnect ||
		b == mqtt.OnDisconnect ||
		b == mqtt.OnPublish ||
		b == mqtt.OnPublished ||
		b == mqtt.OnSubscribed
}

func (h *DebugHook) OnSubscribed(
	cl *mqtt.Client,
	pk packets.Packet,
	reasonCodes []byte,
) {
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("MQTT SUBSCRIPTION")
	fmt.Println("========================================")

	fmt.Printf("Client ID : %s\n", cl.ID)
	fmt.Printf("Filters   : %v\n", pk.Filters)
	fmt.Printf("QoS       : %v\n", reasonCodes)

	fmt.Println("========================================")
	fmt.Println()
}

func (h *DebugHook) OnPublished(
	cl *mqtt.Client,
	pk packets.Packet,
) {
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("MQTT MESSAGE DELIVERED")
	fmt.Println("========================================")

	fmt.Printf("Source Client : %s\n", cl.ID)
	fmt.Printf("Topic         : %s\n", pk.TopicName)
	fmt.Printf("Payload       : %s\n", string(pk.Payload))

	fmt.Println("========================================")
	fmt.Println()
}

// Когда клиент подключился
func (h *DebugHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("MQTT CLIENT CONNECTED")
	fmt.Println("========================================")
	fmt.Printf("Client ID : %s\n", cl.ID)
	fmt.Printf("Remote    : %s\n", cl.Net.Remote)
	fmt.Println("========================================")
	fmt.Println()

	return nil
}

// Когда клиент публикует сообщение
func (h *DebugHook) OnPublish(
	cl *mqtt.Client,
	pk packets.Packet,
) (packets.Packet, error) {

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("MQTT MESSAGE RECEIVED")
	fmt.Println("========================================")

	fmt.Printf("Client ID : %s\n", cl.ID)
	fmt.Printf("Remote    : %s\n", cl.Net.Remote)

	fmt.Printf("Topic     : %s\n", pk.TopicName)

	fmt.Printf("Payload   : %s\n", string(pk.Payload))

	fmt.Printf("Payload bytes: %v\n", pk.Payload)

	fmt.Printf("Retain    : %v\n", pk.FixedHeader.Retain)

	fmt.Println("========================================")
	fmt.Println()

	return pk, nil
}

func (h *DebugHook) OnPacketRead(
	cl *mqtt.Client,
	pk packets.Packet,
) (packets.Packet, error) {

	fmt.Println()
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("MQTT PACKET RECEIVED")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	fmt.Printf("Client : %s\n", cl.ID)
	fmt.Printf("Type   : %s\n", packets.PacketNames[pk.FixedHeader.Type])

	if pk.FixedHeader.Type == packets.Subscribe {
		fmt.Println("SUBSCRIBE FILTERS:")

		for _, filter := range pk.Filters {
			fmt.Printf(
				"  Topic: %s | QoS: %d\n",
				filter.Filter,
				filter.Qos,
			)
		}
	}

	if pk.FixedHeader.Type == packets.Publish {
		fmt.Printf("Topic   : %s\n", pk.TopicName)
		fmt.Printf("Payload : %s\n", string(pk.Payload))
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println()

	return pk, nil
}

func (h *DebugHook) OnPacketSent(
	cl *mqtt.Client,
	pk packets.Packet,
	b []byte,
) {

	fmt.Println()
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("MQTT PACKET SENT")
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

	fmt.Printf("Client : %s\n", cl.ID)
	fmt.Printf("Type   : %s\n", packets.PacketNames[pk.FixedHeader.Type])

	if pk.FixedHeader.Type == packets.Suback {
		fmt.Printf("Reason codes: %v\n", pk.ReasonCodes)
	}

	if pk.FixedHeader.Type == packets.Publish {
		fmt.Printf("Topic   : %s\n", pk.TopicName)
		fmt.Printf("Payload : %s\n", string(pk.Payload))
	}

	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println()
}

// Когда клиент отключился
func (h *DebugHook) OnDisconnect(
	cl *mqtt.Client,
	err error,
	expire bool,
) {

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("MQTT CLIENT DISCONNECTED")
	fmt.Println("========================================")

	fmt.Printf("Client ID : %s\n", cl.ID)

	if err != nil {
		fmt.Printf("Error     : %v\n", err)
	}

	fmt.Println("========================================")
	fmt.Println()
}

// ----------------------------------------------------
// MAIN
// ----------------------------------------------------

func main() {

	// Ctrl+C
	sigs := make(chan os.Signal, 1)

	signal.Notify(
		sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	// MQTT broker
	server := mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	// Разрешаем всем подключаться
	err := server.AddHook(
		new(auth.AllowHook),
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	// Наш Debug Hook
	err = server.AddHook(
		new(DebugHook),
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	// TCP listener
	tcp := listeners.NewTCP(listeners.Config{
		ID:      "tcp1",
		Address: ":1883",
	})

	err = server.AddListener(tcp)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("========================================")
	log.Println("MQTT SERVER")
	log.Println("========================================")
	log.Println("Address: 0.0.0.0:1883")
	log.Println("Status : RUNNING")
	log.Println("========================================")
	log.Println("Waiting for MQTT messages...")
	log.Println()

	// Запускаем сервер
	go func() {

		err := server.Serve()

		if err != nil {
			log.Fatal(err)
		}

	}()

	// Ждём Ctrl+C
	<-sigs

	log.Println()
	log.Println("Stopping MQTT server...")

	server.Close()

	log.Println("MQTT server stopped")
}
