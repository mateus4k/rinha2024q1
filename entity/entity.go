package entity

type TransactionInput struct {
	Value       int32  `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type ExtractOutput struct {
	Balance          BalanceOutput            `json:"saldo"`
	LastTransactions []LastTransactionsOutput `json:"ultimas_transacoes"`
}

type BalanceOutput struct {
	Total int32  `json:"total"`
	Date  string `json:"data_extrato"`
	Limit int32  `json:"limite"`
}

type LastTransactionsOutput struct {
	Value       int32  `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
	Date        string `json:"realizada_em"`
}
