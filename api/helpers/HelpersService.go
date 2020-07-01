package helpers

import (
	"math/rand"
	"strings"
	"time"
)

type HelperServicer interface {
	GenerateID() string
	GenerateRandom(num int) string
	CreateAutoCode(name string) string
}

type HelperService struct {
}

type HelperServiceConfig struct {
}

func NewHelperService(config HelperServiceConfig) HelperService {
	return HelperService{}
}

func (s HelperService) GenerateID() string {
	return s.GenerateRandom(8)
}

func (s HelperService) GenerateRandom(num int) string {
	characters := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, num)
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = characters[seed.Intn(len(characters))]
	}

	return string(b)
}

func (s HelperService) CreateAutoCode(name string) string {
	result := ""

	if len(name) > 3 {
		result = strings.ToLower(name[0:4])
	} else {
		result = s.GenerateRandom(4)
	}

	return result
}
