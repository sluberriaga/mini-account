package account

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/emirpasic/gods/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var TransactionParseError = errors.New("error_parsing_transaction")
var TransactionNotFoundError = errors.New("transaction_not_found_error")

type Account struct {
	mu           sync.RWMutex
	Total        uint64
	Transactions *treebidimap.Map
}

func uuidComparator(a, b interface{}) int {
	uuidA := a.(uuid.UUID)
	uuidB := b.(uuid.UUID)

	return utils.StringComparator(uuidA.String(), uuidB.String())
}

func New() Account {
	return Account{
		Transactions: treebidimap.NewWith(uuidComparator, byEffectiveDateComparator),
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
	a.Transactions.Put(transaction.ID, transaction)

	return true
}

func (a Account) GetBalance() uint64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.Total
}

func (a Account) GetTransactions(offset, limit int) ([]Transaction, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	boundedLimit := limit
	boundedOffset := offset
	if offset >= a.Transactions.Size() {
		boundedOffset = a.Transactions.Size()
		boundedLimit = 0
	} else if limit+offset >= a.Transactions.Size() {
		boundedLimit = a.Transactions.Size() - boundedOffset
	}

	rawTransactions := a.Transactions.Values()
	initial := boundedOffset
	final := boundedOffset + boundedLimit

	rawTransactions = rawTransactions[initial:final]
	transactions := make([]Transaction, 0)
	for _, rawTransaction := range rawTransactions {
		trx, ok := rawTransaction.(*Transaction)
		if !ok {
			return nil, TransactionParseError
		}

		transactions = append(transactions, *trx)
	}

	return transactions, nil
}

func (a Account) GetTransactionByID(uuid uuid.UUID) (Transaction, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if rawTransaction, found := a.Transactions.Get(uuid); found {
		if transaction, ok := rawTransaction.(*Transaction); ok {
			return *transaction, nil
		}

		return Transaction{}, TransactionParseError
	}

	return Transaction{}, TransactionNotFoundError
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
			"code":    "input_validation",
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
		"code":    "insufficient_balance",
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
			"code":    "invalid_offset",
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
			"code":    "invalid_limit",
		})
		return
	}

	transactions, err := s.account.GetTransactions(offset, limit)
	if err == TransactionParseError {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Malformed transactions",
			"code":    "invalid_transactions",
		})
		return
	}

	c.JSON(http.StatusOK, transactions)
	return
}

func (s Service) GetTransactionByID(c *gin.Context) {
	idString := c.Param("id")
	idString = strings.Replace(idString, "/", "", -1)
	id, err := uuid.Parse(idString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Invalid uuid: %s", idString),
			"code":    err.Error(),
		})
		return
	}

	transaction, err := s.account.GetTransactionByID(id)
	if err == TransactionNotFoundError {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Not found!",
			"code":    "not_found",
		})
		return
	}

	if err == TransactionParseError {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad request!",
			"code":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transaction)
	return
}
