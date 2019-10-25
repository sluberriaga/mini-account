package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Account struct {
	mu               sync.RWMutex
	Total            uint64
	Transactions     []Transaction
	TransactionsByID map[uuid.UUID]*Transaction
}

func New() Account {
	return Account{
		Transactions:     []Transaction{},
		TransactionsByID: map[uuid.UUID]*Transaction{},
	}
}

func (a *Account) ProcessTransaction(transaction *Transaction) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	if transaction.Type == "debit" {
		amount := uint64(transaction.Amount * -1)
		if a.Total < amount {
			return false
		}
		a.Total = a.Total - amount
	} else {
		a.Total = a.Total + uint64(transaction.Amount)
	}

	transaction.EffectiveDate = time.Now()
	a.Transactions = append(a.Transactions, *transaction)
	a.TransactionsByID[transaction.ID] = transaction

	return true
}

func (a Account) GetBalance() uint64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.Total
}

func (a Account) GetTransactions(offset, limit int) []Transaction {
	a.mu.RLock()
	defer a.mu.RUnlock()

	boundedLimit := limit
	boundedOffset := offset
	if offset >= len(a.Transactions) {
		boundedOffset = len(a.Transactions)
		boundedLimit = 0
	} else if limit+offset >= len(a.Transactions) {
		boundedLimit = len(a.Transactions) - boundedOffset
	}

	transactions := make([]Transaction, boundedLimit)
	copy(transactions, a.Transactions[len(a.Transactions)-boundedOffset-boundedLimit:len(a.Transactions)-boundedOffset])

	// Reverse transactions
	for i, j := 0, len(transactions)-1; i < j; i, j = i+1, j-1 {
		transactions[i], transactions[j] = transactions[j], transactions[i]
	}
	return transactions
}

func (a Account) GetTransactionByID(uuid uuid.UUID) *Transaction {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if transaction, ok := a.TransactionsByID[uuid]; ok {
		return transaction
	}

	return nil
}

type Service struct {
	account *Account
}

func NewService(account Account) Service {
	return Service{&account}
}

func (s Service) PostTransactionHandler(c *gin.Context) {
	var transactionBody TransactionBody
	if err := c.ShouldBindJSON(&transactionBody); err != nil {

		var errors []validationError
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, NewValidationError(err.Tag(), err.Field()))
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Input validation failed!",
			"code": "input_validation",
			"errors":  errors,
		})
		return
	}

	transaction := transactionBody.ToTransaction()
	if s.account.ProcessTransaction(&transaction) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": "Insufficient balance!",
		"code":  "insufficient_balance",
	})
	return
}

func (s Service) GetBalance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_balance": s.account.GetBalance(),
	})
	return
}

func (s Service) GetTransactions(c *gin.Context) {
	var err error

	var offset int
	if c.Query("offset") == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(c.Query("offset"))
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Offset malformed, please input number",
			"code":  "invalid_offset",
		})
		return
	}

	var limit int
	if c.Query("limit") == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(c.Query("limit"))
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Limit malformed, please input number",
			"code":  "invalid_limit",
		})
		return
	}

	c.JSON(http.StatusOK, s.account.GetTransactions(offset, limit))
	return
}

func (s Service) GetTransactionByID(c *gin.Context) {
	idString := c.Param("id")
	idString = strings.Replace(idString, "/", "", -1)
	id, err := uuid.Parse(idString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Invalid uuid: %s", idString),
			"code":  err.Error(),
		})
		return
	}

	transaction := s.account.GetTransactionByID(id)
	if transaction == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Not found!",
			"code":  "not_found",
		})
		return
	}

	c.JSON(http.StatusOK, *transaction)
	return
}
