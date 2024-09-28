package main

import (
    "dojo/events"
    "encoding/json"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "net"
    "net/http"
    "dojo/entities"
    "github.com/gofiber/fiber/v2/log"
	"github.com/lauro-ss/goe"
	"github.com/lauro-ss/postgres"
    "time"
    "github.com/go-stomp/stomp/v3"
)

var db *Database


type Database struct {
	PedidoRetirada *entities.PedidoRetirada
	*goe.DB
}

func connectDatabase() {
    db = &Database{DB: &goe.DB{}}
    dsn := "host=localhost user=go_go_power_rangers password=powerrangers dbname=alameda port=5432 "
    err := goe.Open(db, postgres.Open(dsn))
    if err != nil{
        log.Error("Erro ao conectar ao banco")
    }
    err = db.Migrate(goe.MigrateFrom(db))
    if err != nil {
        log.Error("Erro realizar migrations")
        return
    }
}

func main() {
    connectDatabase()

    app := fiber.New()

    // Definir a rota para criar um pedido
    app.Post("/pedidos", createPedido)

    // Iniciar o servidor
    app.Listen(":3000")
}

func createPedido(c *fiber.Ctx) error {
    var pedido entities.PedidoRetirada

    // Parseia o corpo da requisição
    if err := c.BodyParser(&pedido); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid input",
        })
    }

    pedido.CreateAt = time.Now()
    pedido.UpdatedAt = time.Now()

    if _, err := db.Insert(db.PedidoRetirada).Value(&pedido); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
           "error": err.Error()       })
    }

    // publish event!
    publish(&pedido)

    return c.Status(http.StatusOK).JSON(pedido)

}

func publish(pedido *entities.PedidoRetirada) {

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
    destination := "/queue/pedidoretirada-queue" // Change this to your queue name

    event :=  events.PedidoRetiradaEvent{
        IdPedido: uuid.New().String(),
        DataPedido: pedido.CreateAt,
        Responsavel: events.Responsavel{
            Nome: "Ranger vermelho",
            Departamento: "Heroi",
        },
        Itens: make([]events.Item, 0),
        Observacoes: "pow power rangeeessssss",
    }

    pedidoJSON, err := json.Marshal(event)
    if err != nil {
        log.Fatalf("Failed to marshal pedido: %v", err)
    }

    err = conn.Send(destination, "application/json", pedidoJSON, stomp.SendOpt.Header("persistent", "true"))
    if err != nil {
        log.Fatalf("Failed to send message: %v", err)
    }
}