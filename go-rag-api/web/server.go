package web

import (
	"log"

	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3"
)

func (s *Server) Routes() *fiber.App {
	s.app = fiber.New(fiber.Config{
		AppName: "go-rag-app",
	})

	s.app.Use(recover.New())

	s.app.Get("/", s.handleIndex)
	s.app.Get("/chat", s.handleChatPage)

	return s.app
}

// handleIndex renders the index.html template.
func (s *Server) handleIndex(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html;charset=utf-8")
	return s.tpl.ExecuteTemplate(c.Response().BodyWriter(), "index.html", fiber.Map{
		"Title": s.title,
	})
}

func (s *Server) handleChatPage(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html;charset=utf-8")

	err := s.tpl.ExecuteTemplate(c.Response().BodyWriter(), "chat.html", fiber.Map{
		"Title": s.title,
		//"CaptionEnabled": s.client.HasVision(),
	})
	if err != nil {
		log.Printf("web template error error %v", err)
	}
	return err
}


