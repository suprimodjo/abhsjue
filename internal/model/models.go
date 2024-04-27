package model

type Transaction struct {
	ID                      int64
	TransactionDate         string
	TransactionType         int
	ReferenceNo             int
	SACode                  string
	InvestorFundUnitACNo    string
	FundCode                string
	AmountNominal           float64
	AmountUnit              float64
	AmountAllUnits          float64
	FeeNominal              float64
	FeeUnit                 float64
	FeePercent              int
	RedmPaymentACSequential int
	RedmPaymentBankBICCode  string
	RedmPaymentBankBIMember string
	RedmPaymentACNo         string
	PaymentDate             string
	TransferType            int
	SAReferenceNo           string
	CreatedAt               *string
	UpdatedAt               *string
}

// Define CustomerTransaction struct to store customer transaction data
type CustomerTransaction struct {
	Customer           string
	TotalAmountNominal float64
	TransactionCount   int
	Transactions       []TransactionDetail // Menyimpan detail transaksi
}

// TransactionDetail menunjukkan detail transaksi, termasuk tanggal dan nominal
type TransactionDetail struct {
	TransactionDate string
	AmountNominal   float64
}
