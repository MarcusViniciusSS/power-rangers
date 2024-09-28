package main

import (
    "net"
    "github.com/gofiber/fiber/v2/log"
    "github.com/go-stomp/stomp/v3"
)


func main() {
    consumer()
}

func consumer() {

    tcpConn, err := net.Dial("tcp", "localhost:61613")
	if err != nil {
		log.Fatalf("Erro ao conectar ao ActiveMQ: %v", err)
	}
	defer tcpConn.Close()

    conn, err := stomp.Connect(tcpConn, stomp.ConnOpt.Login("admin", "admin"))
    if err != nil {
        log.Fatalf("Failed to connect to ActiveMQ: %v", err)
    }
    defer conn.Disconnect()

    // Define the destination queue
    //destination := "/queue/pedidoretirada-queue" // Change this to your queue name

    //pedidoJSON, err := json.Marshal(pedido)
    if err != nil {
        log.Fatalf("Failed to marshal pedido: %v", err)
    }

    //err = conn.Send(destination, "application/json", pedidoJSON, stomp.SendOpt.Header("persistent", "true"))
    if err != nil {
        log.Fatalf("Failed to send message: %v", err)
    }
}