package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

// Consumer representa um consumidor de mensagens do kafka
type Consumer struct {
	ConfigMap *ckafka.ConfigMap // ConfigMap é um mapa de configurações do kafka
	Topics []string // Topics é um array de strings que representa os tópicos que o consumer vai consumir
}

// NewConsumer cria um novo consumidor de mensagens do kafka
func NewConsumer(configMap *ckafka.ConfigMap, topics []string) *Consumer {
	return &Consumer{
		ConfigMap: configMap,
		Topics: topics,
	}
}


// Consume é um método que consome mensagens do kafka ela recebe um canal de mensagens do kafka
func (c *Consumer) Consume(msgChan chan *ckafka.Message) error {
	consumer, err := ckafka.NewConsumer(c.ConfigMap) // cria um novo consumidor de mensagens do kafka
	if err != nil { // se der erro ao criar o consumidor
		panic(err) // erro para restartar a aplicação
	}

	err = consumer.SubscribeTopics(c.Topics, nil) // se não der erro, o consumidor se inscreve nos tópicos que ele vai consumir
	if err != nil { // se der erro ao se inscrever nos tópicos
		panic(err) // erro para restartar a aplicação
	}

	for { // loop infinito para ficar sempre consumindo mensagens
		msg, err := consumer.ReadMessage(-1) // lê a mensagem do kafka
		if err == nil { // se não der erro ao ler a mensagem
			msgChan <- msg // envia a mensagem para o canal de mensagens do kafka
		}
	}
}