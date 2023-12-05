package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID				string
	SellingOrder	*Order
	BuyingOrder		*Order
	Shares			int
	Price			float64
	Total			float64
	DateTime		time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	return &Transaction{
		ID:				uuid.New().String(),
		SellingOrder:	sellingOrder,
		BuyingOrder:	buyingOrder,
		Shares:			shares,
		Price:			price,
		Total:			price * float64(shares),
		DateTime:		time.Now(),
	}
}

func (t *Transaction) CalculateTotal(shares int, price float64) {
	t.Total = t.Price * float64(t.Shares)
}

func (t *Transaction) CloseOrders() {
	if t.BuyingOrder.PendingShares == 0 { //se a ordem de compra não tiver mais ações pendentes
		t.BuyingOrder.Status = "CLOSED" //muda o status da ordem de compra para fechada
	}
	if t.SellingOrder.PendingShares == 0 { //se a ordem de venda não tiver mais ações pendentes
		t.SellingOrder.Status = "CLOSED" //muda o status da ordem de venda para fechada
	}
}

func (t *Transaction) ExecuteTrans() {
	sellingShares := t.SellingOrder.PendingShares //quantidade de ações pendentes
	buyingShares := t.BuyingOrder.PendingShares //quantidade de ações pendentes
	minShares := sellingShares // quantidade mínima de ações inicia em sellingShares
	if buyingShares < sellingShares { // se a quantidade de ações pendentes da ordem de compra for menor que a quantidade de ações pendentes da ordem de venda
		minShares = buyingShares // a quantidade mínima de ações é a quantidade de ações pendentes da ordem de compra
	} // minShares é a quantidade de ações que podem ser transacionadas nesta operação pq o minimo entre as duas ordens é a quantidade de ações que podem ser transacionadas
	t.SellingOrder.Investor.UpdateAssets(t.SellingOrder.Asset.ID, -minShares) //atualiza os ativos do investidor da ordem de venda
	t.SellingOrder.PendingShares -= minShares //atualiza as ações pendentes da ordem de venda
	t.BuyingOrder.Investor.UpdateAssets(t.BuyingOrder.Asset.ID, minShares) //atualiza os ativos do investidor da ordem de compra
	t.BuyingOrder.PendingShares -= minShares //atualiza as ações pendentes da ordem de compra
}