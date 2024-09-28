package events

import "time"

type PedidoRetiradaEvent struct {
	IdPedido string `json:"idPedido"`
	DataPedido time.Time `json:"dataPedido"`
	Responsavel Responsavel
	Observacoes string `json:"observacoes"`
	Itens []Item `json:"itens"`
}

type Responsavel struct {
	Nome string `json:"nome"`
	Departamento string `json:"departamento"`
}

type Item struct {
     Idproduto string `json:"idProduto"`
	 Nomeproduto string `json:"nomeProduto"`
	 Quantidade int64 `json:"quantidade"`
}