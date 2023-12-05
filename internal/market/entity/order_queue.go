package entity

type OrderQueue struct {
	Orders []*Order
}

//Less - Compara 2 valores
func (oq *OrderQueue) Less(i, j int) bool {
	return oq.Orders[i].Price < oq.Orders[j].Price
}
//Swap - Troca 2 valores
func (oq *OrderQueue) Swap(i, j int) {
	oq.Orders[i], oq.Orders[j] = oq.Orders[j], oq.Orders[i]
}

//Len - Retorna o tamanho da lista
func (oq *OrderQueue) Len() int {
	return len(oq.Orders)
}

//Push - Adiciona um valor na lista
func (oq *OrderQueue) Push(x interface{}) {
	oq.Orders = append(oq.Orders, x.(*Order))
}

//Pop - Remove um valor da lista
func (oq *OrderQueue) Pop() interface{} {
	old := oq.Orders // pega a lista de ordens
	n := len(old) // pega o tamanho da lista
	item := old[n-1] // pega o ultimo item da lista
	oq.Orders = old[0 : n-1] // remove o ultimo item da lista
	return item // retorna o item removido
}

//NewOrderQueue - Cria uma nova fila de ordens
func NewOrderQueue() *OrderQueue {
	return &OrderQueue{} // retorna uma nova fila de ordens vazia - não preciso que tenha nada dentro
}