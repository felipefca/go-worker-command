package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"worker/src/config"
	"worker/src/models"
	"worker/src/service"

	"github.com/streadway/amqp"
)

func ConsumerConnect() {
	// Conecta com o RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	// Cria um canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	defer ch.Close()

	// Configura a exchange
	err = ch.ExchangeDeclare(
		config.RabbitExchange, // exchange name
		"direct",              // exchange type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Configura a fila de STOCK com DLQ e TTL
	stockQueue, err := ch.QueueDeclare(
		config.RabbitStockQueue, // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    config.RabbitExchangeDLQ,
			"x-dead-letter-routing-key": config.RabbitStockRoutingKey, //config.RabbitQueueDLQ,
			"x-message-ttl":             config.RabbitTTL,             // TTL de 5 segundos
		}, // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Liga a fila de STOCK à exchange
	err = ch.QueueBind(
		stockQueue.Name,              // queue name
		config.RabbitStockRoutingKey, // routing key
		config.RabbitExchange,        // exchange
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Configura a fila de FIAT com DLQ e TTL
	fiatQueue, err := ch.QueueDeclare(
		config.RabbitFiatQueue, // name
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    config.RabbitExchangeDLQ,
			"x-dead-letter-routing-key": config.RabbitFiatRountingKey, //config.RabbitQueueDLQ,
			"x-message-ttl":             config.RabbitTTL,             // TTL de 5 segundos
		}, // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Liga a fila de FIAT à exchange
	err = ch.QueueBind(
		fiatQueue.Name,               // queue name
		config.RabbitFiatRountingKey, // routing key
		config.RabbitExchange,        // exchange
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Configura a exchange de DLQ
	err = ch.ExchangeDeclare(
		config.RabbitExchangeDLQ, // exchange name
		"fanout",                 // exchange type
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	// Configura a fila de DLQ
	_, err = ch.QueueDeclare(
		config.RabbitQueueDLQ, // Nome da fila DLQ
		true,                  // Durable
		false,                 // Auto-delete
		false,                 // Exclusive
		false,                 // No-wait
		nil,                   // Argumentos extras
	)
	if err != nil {
		log.Fatal(err)
	}

	// Liga a fila de DLQ à exchange DLQ
	err = ch.QueueBind(
		config.RabbitQueueDLQ,    // Nome da fila DLQ
		config.RabbitQueueDLQ,    // Routing key para a fila DLQ
		config.RabbitExchangeDLQ, // Nome da exchange DLQ
		false,                    // No-wait
		nil,                      // Argumentos extras
	)
	if err != nil {
		log.Fatal(err)
	}

	stockMsgs, err := ch.Consume(
		stockQueue.Name, // nome da fila
		"worker",        // consumer tag
		true,            // auto-ack
		false,           // exclusividade
		false,           // sem argumentos extras
		false,           // sem pre-fetch
		nil,             // argumentos extras
	)
	if err != nil {
		log.Fatalf("Erro ao registrar o consumidor: %v", err)
	}

	fiatMsgs, err := ch.Consume(
		fiatQueue.Name, // nome da fila
		"worker2",      // consumer tag
		true,           // auto-ack
		false,          // exclusividade
		false,          // sem argumentos extras
		false,          // sem pre-fetch
		nil,            // argumentos extras
	)
	if err != nil {
		log.Fatalf("Erro ao registrar o consumidor: %v", err)
	}

	// Permanecendo conectado e aguardando mensagens
	forever := make(chan bool)

	go func() {
		// loop para consumir mensagens da fila 1
		for d := range stockMsgs {
			fmt.Printf("Mensagem recebida: %s\n", d.Body)
			service.ProcessStockService(ConsumeStockMessage(d))
			fmt.Printf("Mensagem processada: %s\n", d.Body)
		}
	}()

	go func() {
		// loop para consumir mensagens da fila 2
		for d := range fiatMsgs {
			fmt.Printf("Mensagem da fila 2 recebida: %s\n", d.Body)
			service.ProcessFiatService(ConsumeFiatMessage(d))
			fmt.Printf("Mensagem da fila 2 processada: %s\n", d.Body)
		}
	}()

	fmt.Printf("Aguardando mensagens...\n")
	<-forever
}

func ConsumeFiatMessage(d amqp.Delivery) models.Fiat {
	var fiat models.Fiat

	err := json.Unmarshal(d.Body, &fiat)
	if err != nil {
		log.Fatalf("Erro ao deserializar JSON: %v", err)
	}

	return fiat
}

func ConsumeStockMessage(d amqp.Delivery) models.Stock {
	var stock models.Stock

	err := json.Unmarshal(d.Body, &stock)
	if err != nil {
		log.Fatalf("Erro ao deserializar JSON: %v", err)
	}

	return stock
}
