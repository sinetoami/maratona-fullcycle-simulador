package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sinetoami/maratona-fullcycle-simulator/entity"
	"github.com/sinetoami/maratona-fullcycle-simulator/queue"
	"github.com/streadway/amqp"
)

// toda ordem que for processada vai ficar dentro desse slice
var active []string

func init() {
	// carrega as variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	in := make(chan []byte)

	// cria uma canal de conexão e comunicação do RabbitMQ
	// toda vez que estiver consumindo uma mensagem, essa mensagem
	// vai ser enviada para o canal 'in'
	ch := queue.Connect()
	queue.StartConsuming(in, ch)

	// tudo que chegar de ordem de serviços:
	for msg := range in {
		// transformando o json recebido como mensagem e transformando em objeto
		var order entity.Order
		err := json.Unmarshal(msg, &order)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("New order Received: ", order.UUID)
		start(order, ch) //
	}
}

//
func start(order entity.Order, ch *amqp.Channel) {
	// se a ordem de serviço ainda não estiver ativa
	if !strInSlice(order.UUID, active) {
		active = append(active, order.UUID)

		// executa assincronamente em background sem precisar ficar esperando
		// uma ordem terminar para iniciar uma nova
		go SimulatorWorker(order, ch)
	} else {
		fmt.Println("Order", order.UUID, "was already complete or is on going...")
	}
}

// SimulatorWorker generate a driver trail
func SimulatorWorker(order entity.Order, ch *amqp.Channel) {
	orderfile, err := os.Open("destinations/" + order.Destination + ".txt")
	if err != nil {
		panic(err.Error())
	}

	defer orderfile.Close()

	scanner := bufio.NewScanner(orderfile)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")
		json := destinationToJSON(order, data[0], data[1])

		// fica enviando infomações de lat e lng a cada 1 segundo
		time.Sleep(1 * time.Second)
		queue.Notify(string(json), ch)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// quando a trajetoria foi concluída
	json := destinationToJSON(order, "0", "0")
	queue.Notify(string(json), ch)
}

func strInSlice(str string, list []string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}

	return false
}

// transforma um objeto 'order' em um json
func destinationToJSON(order entity.Order, lat, lng string) []byte {
	dest := entity.Destination{
		Order: order.UUID,
		Lat:   lat,
		Lng:   lng,
	}

	json, _ := json.Marshal(dest)
	// if err != nil {
	// 	fmt.Errorf("%s: Json destination marshal failed", err)
	// }

	return json
}
