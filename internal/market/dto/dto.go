package dto

// preciso definir contratos de dados que vou receber da minha ordem e retornar das transações
// para isso vou receber os dados em uma camada "burra" (sem comportamento, apenas valores) que consegue andar pelas camadas da minha aplicação, esse são os dtos
// similares a interfaces
type TradeInput struct {
	OrderID			string	`json:"order_id"`
	AssetID			string	`json:"asset_id"`
	InvestorID		string	`json:"investor_id"`
	CurrentShares	int		`json:"current_shares"`
	Price			float64	`json:"price"`
	Shares			int		`json:"shares"`
	OrderType		string	`json:"order_type"`
}

type OrderOutput struct {
	OrderID				string					`json:"order_id"`
	AssetID				string					`json:"asset_id"`
	InvestorID			string					`json:"investor_id"`
	OrderType			string					`json:"order_type"`
	Status				string					`json:"status"`
	Partial 			int 					`json:"partial"`
	Shares 				int						`json:"shares"`
	TransactionOutput	[]*TransactionOutput	`json:"transactions"`
}

type TransactionOutput struct {
	TransactionID	string	`json:"transaction_id"`
	BuyerID 		string	`json:"buyer_id"`
	SellerID 		string	`json:"seller_id"`
	AssetID 		string	`json:"asset_id"`
	Price 			float64	`json:"price"`
	Shares 			int		`json:"shares"`
}