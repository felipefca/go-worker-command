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
		false,           // auto-ack
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
		false,          // auto-ack
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
			fmt.Printf("Mensagem recebida")

			message, err := ConsumeStockMessage(d)
			if err != nil {
				RejectMessage(d)
				continue
			}

			service.ProcessStockService(message)

			AckMessage(d)
		}
	}()

	go func() {
		// loop para consumir mensagens da fila 2
		for d := range fiatMsgs {
			fmt.Printf("Mensagem recebida")

			message, err := ConsumeFiatMessage(d)
			if err != nil {
				RejectMessage(d)
				continue
			}

			service.ProcessFiatService(message)

			AckMessage(d)
		}
	}()

	fmt.Printf("Aguardando mensagens...\n")
	<-forever
}

func ConsumeFiatMessage(d amqp.Delivery) (models.Fiat, error) {
	var fiat models.Fiat

	err := json.Unmarshal(d.Body, &fiat)
	if err != nil {
		fmt.Printf("Erro ao deserializar JSON: %v", err)
		return models.Fiat{}, err
	}

	return fiat, nil
}

func ConsumeStockMessage(d amqp.Delivery) (models.Stock, error) {
	var stock models.Stock

	err := json.Unmarshal(d.Body, &stock)
	if err != nil {
		fmt.Printf("Erro ao deserializar JSON: %v", err)
		return models.Stock{}, err
	}

	return stock, nil
}

func RejectMessage(msg amqp.Delivery) {
	if err := msg.Nack(false, false); err != nil {
		log.Printf("Erro ao enviar mensagem para DLQ: %v", err)
	} else {
		log.Printf("Mensagem enviada para DLQ")
	}
}

func AckMessage(msg amqp.Delivery) {
	if err := msg.Ack(false); err != nil {
		log.Printf("Erro ao confirmar o recebimento da mensagem: %v", err)
	} else {
		log.Printf("Mensagem processada com sucesso")
	}
}
