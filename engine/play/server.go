package play

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

type Server struct {
	app *fiber.App
	cfg *Config
}

func NewServer(cfg *Config) *Server {
	app := fiber.New()
	server := &Server{
		app: app,
		cfg: cfg,
	}

	app.Use(cors.New())

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", server.HandleWebsocket())

	return server
}

func (s *Server) Run() {
	log.Fatal(s.app.Listen(":8080"))
}

func (s *Server) HandleWebsocket() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		log.Println("Websocket connection established")
		defer log.Println("Websocket connection closed")

		conn := NewConnection(c, s.cfg)

		conn.playUntilPlayerMove()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			message := Message{}
			err = json.Unmarshal(msg, &message)
			if err != nil {
				log.Println("unmarshal:", err)
				break
			}

			conn.ProcessMessage(&message)
		}
	})
}
