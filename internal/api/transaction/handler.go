package api

import (
	"net/http"

	"auth/internal/services/transaction"

	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	ToUser string `json:"to_user" binding:"required"`
	Amount int    `json:"amount" binding:"required,min=1"`
}

type BuyRequest struct {
	ItemType string `json:"item_type" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type TransactionHandler interface {
	TransferCoins(c *gin.Context)
	BuyItem(c *gin.Context)
}

type transactionHandler struct {
	transactionSrv transaction.Service
}

func New(
	transactionSrv transaction.Service,
) TransactionHandler {
	return &transactionHandler{
		transactionSrv: transactionSrv,
	}
}

func (h *transactionHandler) TransferCoins(c *gin.Context) {
	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		c.Error(err)
		return
	}

	if err := h.transactionSrv.TransferCoins(c, req.ToUser, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coins transferred successfully"})
}

func (h *transactionHandler) BuyItem(c *gin.Context) {
	itemType := c.Param("item")
	if itemType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing item type"})
		return
	}

	if err := h.transactionSrv.BuyItem(c, itemType, 1); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item purchased successfully"})
}
