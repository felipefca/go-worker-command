package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	RabbitUserName        = ""
	RabbitPassword        = ""
	RabbitHost            = ""
	RabbitExchange        = ""
	RabbitStockQueue      = ""
	RabbitFiatQueue       = ""
	RabbitExchangeDLQ     = ""
	RabbitQueueDLQ        = ""
	RabbitTTL             = 0
	RabbitFiatRountingKey = ""
	RabbitStockRoutingKey = ""
	UrlLastContation      = ""
	UrlLastDaysContation  = ""
	LastDaysContation     = 0
	UrlNasdaq             = ""
	PatchNasdaqQuote      = ""
	NasdaqToken           = ""
)

func Load() {
	var erro error

	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	LastDaysContation, erro = strconv.Atoi(os.Getenv("LAST_DAYS_CONTATION"))
	if erro != nil {
		LastDaysContation = 2
	}

	RabbitTTL, erro = strconv.Atoi(os.Getenv("RABBITMQ_TTL"))
	if erro != nil {
		RabbitTTL = 5000
	}

	RabbitUserName = os.Getenv("RABBITMQ_USER")
	RabbitPassword = os.Getenv("RABBITMQ_PASS")
	RabbitHost = os.Getenv("RABBITMQ_HOST")
	RabbitExchange = os.Getenv("RABBITMQ_EXCHANGE")
	RabbitFiatQueue = os.Getenv("RABBITMQ_FIAT_QUEUE")
	RabbitStockQueue = os.Getenv("RABBITMQ_STOCK_QUEUE")
	RabbitExchangeDLQ = os.Getenv("RABBITMQ_DLQ_EXCHANGE")
	RabbitQueueDLQ = os.Getenv("RABBITMQ_DLQ_QUEUE")
	RabbitFiatRountingKey = os.Getenv("RABBITMQ_FIAT_ROUTING_KEY")
	RabbitStockRoutingKey = os.Getenv("RABBITMQ_STOCK_ROUTING_KEY")
	UrlLastContation = os.Getenv("URL_LAST_COTATION")
	UrlLastDaysContation = os.Getenv("URL_LAST_DAY_COTATION")
	UrlNasdaq = os.Getenv("URL_NASDAQ")
	PatchNasdaqQuote = os.Getenv("PATCH_QUOTE")
	NasdaqToken = os.Getenv("NASDAQ_TOKEN")
}
