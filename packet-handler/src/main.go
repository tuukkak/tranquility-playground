package main

import (
    "fmt"
    "net"
    "net/http"
    "log"
    "time"
    "bufio"
    "github.com/streadway/amqp"
)

func main() {
    fmt.Println("Server started")

    pc, err := net.ListenPacket("udp", ":3000")
	if err != nil {
        log.Fatal(err)
    }
    defer pc.Close()
    
    ln, err := net.Dial("tcp", "spell-handler:3002")
    if err != nil {
        log.Fatal(err)
    }

    conn, err := amqp.Dial("amqp://rabbitmq-server")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "movements", // name
        true,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    if err != nil {
        log.Fatal(err)
    }

    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    if err != nil {
        log.Fatal(err)
    }

    timer := time.Now()

    go listenPackets(msgs, &timer)

	for {
		buf := make([]byte, 1024)
        n, addr, err := pc.ReadFrom(buf)
        fmt.Println("Packet received!")
		if err != nil {
			continue
        }
        timer = time.Now()
        go pingSpellHandler(ln, &timer)
        go pingGameData(pc, addr, buf[:n], &timer)
        go pingMovementHandler(q, ch, &timer)
	}
}

func listenPackets(msgs <-chan amqp.Delivery, timer *time.Time) {
    for d := range msgs {
        log.Printf("Received a message: %s", d.Body)
        t := time.Now()
        log.Printf("RabbitMQ took:", t.Sub(*timer).String())
    }
}

func pingSpellHandler(conn net.Conn, timer *time.Time) {
    fmt.Fprintf(conn, "Hello\n")
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Println(message)
    t := time.Now()
    log.Printf("TCP listener took:", t.Sub(*timer).String())
}

func pingGameData(pc net.PacketConn, addr net.Addr, buf []byte, timer *time.Time) {
    http.Get("http://game-data:3001/")
    t := time.Now()
    log.Printf("HTTP REST API took:", t.Sub(*timer).String())
}

func pingMovementHandler(q amqp.Queue, ch *amqp.Channel, timer *time.Time) {
    ch.Publish(
        "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing {
            ContentType: "text/plain",
            Body: []byte("Hi"),
        })
}