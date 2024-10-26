package usecase

import "github.com/ei-sugimoto/logresdb/api/internal/ports/in"

type GreetService struct {
}

func NewGreetService() in.GreetService {
	return &GreetService{}
}

func (g *GreetService) Greet(name string) string {
	return "Hello, " + name
}
