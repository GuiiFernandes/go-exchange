package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order []*Order
	Transactions []*Transaction
	OrdersChan chan *Order //canal para receber todas as ordens para adicionar no book
	OrdersChanOut chan *Order //canal para enviar todas as ordens que foram adicionadas no book para enviar para o kafka
	Wg *sync.WaitGroup //Ajuda a sincronizar nossas threads, quando duas operações acontecerem (terminarem) falo para o waitgroup que terminou e ele libera a thread
}

// NewBook cria um novo book
func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order: []*Order{},
		Transactions: []*Transaction{},
		OrdersChan: orderChan,
		OrdersChanOut: orderChanOut,
		Wg: wg,
	}
}

// Trade - Faz a troca de ordens
// Vai rodar em uma thread separada - então toda vez que uma ordem chegar no canal, ele vai executar a função
func (b *Book) Trade() {
	buyOrders := make(map[string]*OrderQueue) //Cria um mapa de ordens de compra para que cada ativo tenha sua fila de ordens de compra
	sellOrders := make(map[string]*OrderQueue) //Cria um mapa de ordens de venda para que cada ativo tenha sua fila de ordens de venda
	// buyOrders := NewOrderQueue() //Cria uma fila de oderns de compra
	// sellOrders := NewOrderQueue() //Cria uma fila de ordens de venda

	// heap nos permirte usar toda a estrutura de dados do go de forma mais performática
	// heap.Init(buyOrders) //inicializa a fila de ordens de compra
	// heap.Init(sellOrders) //inicializa a fila de ordens de venda

	//loop infinito para ficar escutando o canal de ordens
	for order := range b.OrdersChan{ //para cada ordem que chegar no canal
		asset := order.Asset.ID
		if buyOrders[asset] == nil { //se não tiver ordens de compra para o ativo
			buyOrders[asset] = NewOrderQueue() //cria uma nova fila de ordens de compra para o ativo
			heap.Init(buyOrders[asset]) //inicializa a fila de ordens de compra para o ativo
		}
		if sellOrders[asset] == nil { //se não tiver ordens de compra para o ativo
			sellOrders[asset] = NewOrderQueue() //cria uma nova fila de ordens de compra para o ativo
			heap.Init(sellOrders[asset]) //inicializa a fila de ordens de compra para o ativo
		}
		if order.OrderType == "BUY" { //se for uma ordem de compra
			buyOrders[asset].Push(order) //adiciona na fila de ordens de compra
			if sellOrders[asset].Len() > 0 && sellOrders[asset].Orders[0].Price <= order.Price { //se tiver ordens de venda e o preço da ordem de venda for menor ou igual ao preço da ordem de compra
				sellOrder := sellOrders[asset].Pop().(*Order) // remove a ordem de venda da fila e salva na variavel sellOrder
				if sellOrder.PendingShares > 0 { //se a ordem de venda tiver ações pendentes para vender
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price) //cria uma nova transação com a ordem de venda, a ordem de compra, a quantidade de ações da ordem de compr e o preço da ordem de venda
					b.AddTransaction(transaction, b.Wg) // adiciona a transação no book com o waitgroup
					sellOrder.Transactions = append(sellOrder.Transactions, transaction) //adiciona a transação na ordem de venda
					order.Transactions = append(order.Transactions, transaction) //adiciona a transação na ordem de compra
					// preciso adicionar a transação nas duas ordens para mostrar que ouve uma transação entre elas
					b.OrdersChanOut <- sellOrder //retorna a ordem de venda para o canal de saída
					b.OrdersChanOut <- order //retorna a ordem de compra para o canal de saída
					// preciso retornar as duas ordens para o canal de saída para mostrar que ouve uma transação entre elas
					if sellOrder.PendingShares > 0 { // se a ordem de venda ainda tiver ações pendentes, adiciona ela na fila de ordens de venda
						sellOrders[asset].Push(sellOrder) //adiciona a ordem de venda na fila de ordens de venda
						// para quando aparecer outra ordem de compra com o preço maior ou igual ao preço da ordem de venda ela seja liquidada.
					}
				}
			}
		} else { //se for uma ordem de venda
			sellOrders[asset].Push(order) //adiciona na fila de ordens de venda
			if buyOrders[asset].Len() > 0 && buyOrders[asset].Orders[0].Price >= order.Price { //se tiver ordens de compra e o preço da ordem de compra for maior ou igual ao preço da ordem de venda
				buyOrder := buyOrders[asset].Pop().(*Order) // remove a ordem de compra da fila e salva na variavel buyOrder
				if buyOrder.PendingShares > 0 { //se a ordem de compra tiver ações pendentes para comprar
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price) //cria uma nova transação com a ordem de venda, a ordem de compra, a quantidade de ações da ordem de venda e o preço da ordem de compra
					b.AddTransaction(transaction, b.Wg) // adiciona a transação no book com o waitgroup
					buyOrder.Transactions = append(buyOrder.Transactions, transaction) //adiciona a transação na ordem de compra
					order.Transactions = append(order.Transactions, transaction) //adiciona a transação na ordem de venda
					b.OrdersChanOut <- buyOrder //retorna a ordem de compra para o canal de saída
					b.OrdersChanOut <- order //retorna a ordem de venda para o canal de saída
					if buyOrder.PendingShares > 0 { // se a ordem de compra ainda tiver ações pendentes, adiciona ela na fila de ordens de compra
						buyOrders[asset].Push(buyOrder) //adiciona a ordem de compra na fila de ordens de compra
					}
				}
			}
		}
	}
}

// AddTransaction - adiciona uma transação no book, nas ordens e altera os ativos dos investidores
func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) { //recebe uma transação e o waitgroup para gerenciar as threads quando terminar a operação.
	defer wg.Done() //quando terminar a operação, falo para o waitgroup que terminou (defer é executado quando a função termina) (Done termina o waitgroup)
	transaction.ExecuteTrans() //executa a transação alterando as ações pendentes das ordens e os ativos dos investidores
	// o ponteiro ajuda nesse caso para a partir de transaction eu conseguir acessar o modificar os ativos do investidor (outro objeto) a partir do local da memória onde ele está armazenado
	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price) //calcula o total da transação que é a quantidade de ações transacionadas vezes o preço da ordem de compra
	transaction.CloseOrders() //fecha as ordens da transação
	b.Transactions = append(b.Transactions, transaction) //adiciona a transação no book
}