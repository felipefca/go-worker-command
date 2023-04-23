package main

import (
	"fmt"
	"worker/src/config"
	"worker/src/consumer"
)

func main() {
	config.Load()

	fmt.Println("Rodando Worker")

	consumer.ConsumerConnect()
}
