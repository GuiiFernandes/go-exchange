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