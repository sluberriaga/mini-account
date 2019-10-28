package account

import (
	"encoding/json"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type Transaction struct {
	ID            uuid.UUID `json:"id"`
	EffectiveDate time.Time `json:"date"`
	Type          string    `json:"type"`
	Amount        int64     `json:"amount"`
}

func byEffectiveDateComparator(a, b interface{}) int {
	trxa := a.(*Transaction)
	trxb := b.(*Transaction)

	switch {
	case trxa.EffectiveDate.Before(trxb.EffectiveDate):
		return 1
	case trxb.EffectiveDate.Before(trxa.EffectiveDate):
		return -1
	default:
		return 0
	}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID            string    `json:"id"`
		EffectiveDate time.Time `json:"date"`
		Type          string    `json:"type"`
		Amount        int64     `json:"amount"`
	}{
		ID:            t.ID.String(),
		EffectiveDate: t.EffectiveDate,
		Type:          t.Type,
		Amount:        t.Amount,
	})
}

type TransactionBody struct {
	Type   string `json:"type" binding:"required"`
	Amount int64  `json:"amount" binding:"required"`
}

func (t TransactionBody) ToTransaction() Transaction {
	return Transaction{
		ID:     uuid.New(),
		Type:   t.Type,
		Amount: t.Amount,
	}
}

func TransactionBodyValidation(structLevel validator.StructLevel) {
	transactionBody := structLevel.Current().Interface().(TransactionBody)

	if transactionBody.Type != "credit" && transactionBody.Type != "debit" {
		structLevel.ReportError(transactionBody.Type, "Type", "type", "should_be_credit_or_debit", "")
	}

	if transactionBody.Amount == 0 {
		structLevel.ReportError(transactionBody.Type, "Amount", "amount", "should_not_be_0", "")
	}

	if transactionBody.Type == "credit" && transactionBody.Amount < 0 {
		structLevel.ReportError(transactionBody.Type, "Amount", "amount", "credit_should_be_positive", "")
	}

	if transactionBody.Type == "debit" && transactionBody.Amount > 0 {
		structLevel.ReportError(transactionBody.Type, "Amount", "amount", "debit_should_be_negative", "")
	}
}
