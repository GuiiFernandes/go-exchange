package transformer

import (
	"github.com/guiifernandes/go-exchange/internal/market/dto"
	"github.com/guiifernandes/go-exchange/internal/market/entity"
)

// transforma o input do kafka em um objeto do domínio (dados crus em dados de negócio)
func TransformInput(input dto.TradeInput) *entity.Order { // recebe um input do kafka e retorna um objeto do domínio
	asset := entity.NewAsset(input.AssetID, input.AssetID, 1000) // cria um novo ativo a partir do input do kafka
	investor := entity.NewInvestor(input.InvestorID) // cria um novo investidor a partir do input do kafka
	order := entity.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType) // cria uma nova ordem a partir do input do kafka
	if input.CurrentShares > 0 { // se o investidor já tiver ações
		assetPosition := entity.NewAssetPosition(input.AssetID, input.CurrentShares) // cria uma nova posição de ativo a partir do input do kafka
		investor.AddAsset(assetPosition) // adiciona a posição de ativo ao investidor
	}
	return order // retorna a ordem
}

// transforma o output do domínio em um objeto de saída (dados de negócio em dados crus)
func TransformOutput(order *entity.Order) *dto.OrderOutput { // recebe uma ordem do domínio e retorna um objeto de saída
	output := &dto.OrderOutput{ // cria um novo objeto de saída
		OrderID: order.ID,
		InvestorID: order.Investor.ID,
		AssetID: order.Asset.ID,
		OrderType: order.OrderType,
		Status: order.Status,
		Partial: order.PendingShares,
		Shares: order.Shares,
	} // preenche os campos do objeto de saída com os dados da ordem

	// como uma ordem pode ter várias transações, preciso criar um array de transações no formato de saída
	var transactionsOutput []*dto.TransactionOutput // cria um array de transações de saída
	for _, t := range order.Transactions { // para cada transação da ordem
		transactionOutput := &dto.TransactionOutput{ // cria uma nova transação de saída
			TransactionID: t.ID,
			BuyerID: t.BuyingOrder.Investor.ID,
			SellerID: t.SellingOrder.Investor.ID,
			AssetID: t.SellingOrder.Asset.ID,
			Price: t.Price,
			Shares: t.SellingOrder.Shares - t.SellingOrder.PendingShares,
		}
		transactionsOutput = append(transactionsOutput, transactionOutput) // adiciona a transação de saída ao array de transações de saída
	}
	output.TransactionOutput = transactionsOutput // adiciona o array de transações de saída ao objeto de saída
	return output // retorna o objeto de saída
}