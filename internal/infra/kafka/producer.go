package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

// Producer representa um produtor de mensagens do kafka
type Producer struct {
	ConfigMap *ckafka.ConfigMap
}


// NewKafkaProducer cria um novo produtor de mensagens do kafka
func NewKafkaProducer(configMap *ckafka.ConfigMap) *Producer {
	return &Producer{
		ConfigMap: configMap,
	}
}

// Publish é um método que publica mensagens no kafka ele recebe uma mensagem, uma chave e um tópico a ser publicado a mensagem
func (p *Producer) Publish(msg interface{}, key []byte, topic string) error {
	producer, err := ckafka.NewProducer(p.ConfigMap) // cria um novo produtor de mensagens do kafka
	if err != nil { // se der erro ao criar o produtor
		return err  // retorna o erro
	}

	message := &ckafka.Message{ // cria uma mensagem do kafka
		TopicPartition: ckafka.TopicPartition{ // informa o tópico e a partição que a mensagem deve ser publicada
			Topic: &topic,
			Partition: ckafka.PartitionAny,
		},
		Key: key, // informa a chave da mensagem (particularidade do kafka)
		Value: msg.([]byte), // informa o valor da mensagem em bytes
	}

	err = producer.Produce(message, nil) // publica a mensagem no kafka
	if err != nil { // se der erro ao publicar a mensagem
		return err // retorna o erro
	}

	return nil // se não der erro retorna nil pq sou obrigado a retornar um erro, mesmo que vazio
}