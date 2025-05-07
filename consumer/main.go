package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main(){
	amqpServerURL := "amqp://guest:guest@localhost:5672/"

	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		log.Fatalf("erro ao conectar com RabbitMQ: %v", err)
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		log.Fatalf("erro ao abrir channel com RabbitMQ: %v", err)
	}
	defer channelRabbitMQ.Close()

	mensagens, err := channelRabbitMQ.Consume(
		"service1",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("erro ao ler mensagens da fila no RabbitMQ: %v", err)
	}

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Wainting for messages")

	forever := make(chan bool)

	go func (){
		for mensagem := range mensagens {
			log.Printf("Mensagem processada: %s\n", string(mensagem.Body))
		}
	}()
	<- forever
}