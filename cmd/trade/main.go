package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"

	"github.com/guiifernandes/go-exchange/internal/infra/kafka"
	"github.com/guiifernandes/go-exchange/internal/market/dto"
	"github.com/guiifernandes/go-exchange/internal/market/entity"
	"github.com/guiifernandes/go-exchange/internal/market/transformer"
)

// função principal da aplicação. Aqui é onde tudo começa
func main() {
	err := godotenv.Load()
    if err != nil {
        panic("Error to load .env file")
    }
	ordersIn := make(chan *entity.Order) // canal de entrada de ordens
	ordersOut := make(chan *entity.Order) // canal de saída de ordens
	wg := &sync.WaitGroup{} // cria um grupo de espera
	defer wg.Wait() // espera todas as goroutines terminarem

	kafkaMsgChan := make(chan *ckafka.Message) // canal de mensagens do kafka
	configMap := &ckafka.ConfigMap{ // cria um mapa de configurações do kafka
		"bootstrap.servers": os.Getenv("KAFKA_CONN_SERVER"), // endereço de conexão do kafka com variável de ambiente
		"group.id":          os.Getenv("KAFKA_GROUP_ID"), // server para juntar diversos consumidores em um grupo, vamos trabalhar apenas com um
		"auto.offset.reset": "latest", // offset da ordem de leitura do kafka, escolhemos ler pelo final dos dados (opcional)
	}
	producer := kafka.NewKafkaProducer(configMap) // cria um novo produtor kafka
	kafka := kafka.NewConsumer(configMap, []string{"input"}) // cria um novo consumidor kafka com o configMap e os tópicos que quero ler em um array, nesse caso é só o "input"

	// Vamos começar a consumir as mensagens do kafka
	go kafka.Consume(kafkaMsgChan) // consome as mensagens que estão no chanel do kafka
	// Como o Consume tem um loop infinito preciso gerar uma nova thread para ele não travar o restante da aplicação

	// Vamos começar a processar as mensagens do kafka
	// recebe no canal do kafka, joga no input, processa, joga no output e depois publica no kafka
	book := entity.NewBook(ordersIn, ordersOut, wg) // cria um novo livro de ordens
	go book.Trade() // inicia o trade do livro de ordens em uma nova thread pois também tem um loop infinito.

	//Vamos criar um loop infinito que vai receber os dados do kafka e publicar na nossa order/input

	go func() {
		for msg := range kafkaMsgChan { // para cada mensagem que chegar no canal do kafka
			wg.Add(1) // adiciona uma goroutine no grupo de espera
			fmt.Println(string(msg.Value)) // imprime a mensagem que chegou no canal do kafka
			tradeInput := dto.TradeInput{} // recebe os dados crus em JSON
			err := json.Unmarshal(msg.Value, &tradeInput) // converte os dados crus em JSON para a struct TradeInput alterando o valor do tradeInput por referência
			if err != nil { // se der erro ao converter os dados crus em JSON para a struct TradeInput
				panic(err) // erro para restartar a aplicação
			}
			order := transformer.TransformInput(tradeInput) // transforma os dados crus em JSON para a struct Order
			ordersIn <- order // envia a ordem para o canal de entrada de ordens
			// Nosso book lê esse canal e processa a ordem. Depois de processar a ordem ele envia para o canal de saída de ordens para ser publicado no kafka e ser visualizado pelos consumers
		}
	}()

	// Vamos criar um loop infinito que vai receber os dados do canal de saída de ordens e publicar no kafka
	for res := range ordersOut {
		output := transformer.TransformOutput(res) // transforma a ordem para o formato de saída
		outputJSON, err := json.MarshalIndent(output, "", "   ") // converte a ordem no formato cru para JSON
		fmt.Println(string(outputJSON)) // imprime a ordem em JSON
		if err != nil { // se der erro ao converter a ordem para JSON
			fmt.Println(err) // erro para restartar a aplicação
		}
		producer.Publish(outputJSON, []byte("orders"), "output") // publica a ordem no kafka
	}
}