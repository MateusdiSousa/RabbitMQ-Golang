package main

import (
	"io"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

func main() {
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

	_, err = channelRabbitMQ.QueueDeclare(
		"service1",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("erro ao declarar fila no RabbitMQ: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			rbody, err := ReadRequestBody(r)
			if err != nil {
				io.WriteString(w, "Request body inválido")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err = channelRabbitMQ.Publish(
				"",
				"service1",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        rbody,
				},
			)
			if err != nil {
				io.WriteString(w, "Service indisponível")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			io.WriteString(w, "mensagem recebida")
			w.WriteHeader(http.StatusAccepted)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":3333", mux))
}

func ReadRequestBody(request *http.Request) ([]byte, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
