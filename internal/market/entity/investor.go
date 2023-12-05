package entity

// struct (classe) do investidor, com id, nome e lista de ativos
type Investor struct {
	ID 		string
	Name	string
	Assets	[]*AssetPosition
}
// construtor do investidor com id e lista de ativos, cria um novo investidor.
func NewInvestor(id string) *Investor {
	return &Investor{
		ID: id,
		Assets: []*AssetPosition{},
	}
}
// metodo para adicionar ativos ao investidor,
func (i *Investor) AddAsset(assetPosition *AssetPosition) {
	i.Assets = append(i.Assets, assetPosition) // adiciona ativos ao investidor pegando a lista de ativos e adicionando o ativo no final da lista
}
// metodo para atualizar ativos do investidor, recebe o id do ativo e a quantidade de ações
func (i *Investor) UpdateAssets(assetID string, qtdShares int) {
	assetPosition := i.GetAsset(assetID) // pega o ativo pelo id
	if assetPosition == nil { // se o ativo for nulo, cria um novo ativo
		i.Assets = append(i.Assets, NewAssetPosition(assetID, qtdShares)) // adiciona o ativo na lista de ativos
	} else { // se o ativo não for nulo, atualiza o ativo
		assetPosition.Shares += qtdShares // atualiza a quantidade de ações do ativo
	}
}
// metodo para pegar ativos do investidor, recebe o id do ativo e retorna o ativo
func (i *Investor) GetAsset(assetID string) *AssetPosition {
	for _, asset := range i.Assets { // percorre a lista de ativos do investidor
		if asset.AssetID == assetID { // se o id do ativo for igual ao id do ativo passado como parametro, retorna o ativo
			return asset
		}
	}
	return nil // se não encontrar o ativo, retorna nulo
}
// struct (classe) do ativo, com id e quantidade de ações
type AssetPosition struct {
	AssetID	string
	Shares	int
}
// construtor do ativo com id e quantidade de ações, cria um novo ativo para o Assets do investidor.
func NewAssetPosition(assetID string, shares int) *AssetPosition {
	return &AssetPosition{
		AssetID: assetID,
		Shares: shares,
	}
}
//ponteiro indica o endereço de memoria do ativo para que se alterado reflita no investidor em qualquer lugar do codigo