package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"

	"github.com/mateus4k/rinha2024q1/db"
	"github.com/mateus4k/rinha2024q1/entity"
)

type Controller struct {
	queries *db.Queries
	conn    *pgxpool.Pool
}

func NewController(queries *db.Queries, conn *pgxpool.Pool) *Controller {
	return &Controller{
		queries: queries,
		conn:    conn,
	}
}

func (ctrl *Controller) GetTransactions(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	if idInt > 5 {
		return c.SendStatus(http.StatusNotFound)
	}

	accountId := int32(idInt)
	account, err := ctrl.queries.GetAccount(c.Context(), accountId)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	transactions, err := ctrl.queries.GetTransactions(c.Context(), accountId)
	if err != nil {
		transactions = []db.GetTransactionsRow{}
	}

	return c.Status(200).JSON(entity.ExtractOutput{
		Balance: entity.BalanceOutput{
			Total: account.Balance,
			Date:  time.Now().UTC().Format(time.RFC3339Nano),
			Limit: account.Lim,
		},
		LastTransactions: lo.Map(transactions, func(t db.GetTransactionsRow, index int) entity.LastTransactionsOutput {
			return entity.LastTransactionsOutput{
				Value:       t.Amount,
				Type:        t.Type,
				Description: t.Description,
				Date:        t.Date.Time.Format(time.RFC3339Nano),
			}
		}),
	})
}

func (ctrl *Controller) CreateTransaction(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	if idInt > 5 {
		return c.SendStatus(http.StatusNotFound)
	}

	transaction := new(entity.TransactionInput)
	if err := c.BodyParser(transaction); err != nil {
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	if transaction.Value <= 0 || transaction.Description == "" || len(transaction.Description) > 10 {
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	accountId := int32(idInt)
	account, err := ctrl.queries.GetAccount(c.Context(), accountId)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	switch transaction.Type {
	case "c":
		account.Balance += transaction.Value
	case "d":
		if account.Balance-transaction.Value < -account.Lim {
			return c.SendStatus(http.StatusUnprocessableEntity)
		}
		account.Balance -= transaction.Value
	default:
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	if err = ctrl.queries.InsertTransaction(c.Context(), db.InsertTransactionParams{
		PAccountID:   accountId,
		PAmount:      transaction.Value,
		PType:        transaction.Type,
		PDescription: transaction.Description,
	}); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"limite": account.Lim,
		"saldo":  account.Balance,
	})
}
