package info

import (
	"net/http"

	"auth/internal/domain"
	"auth/internal/services/info"

	"github.com/gin-gonic/gin"
)

type InfoResponse struct {
	Coins       int                `json:"coins"`
	Inventory   []domain.Inventory `json:"inventory"`
	CoinHistory struct {
		Sent     []domain.Transaction `json:"sent"`
		Received []domain.Transaction `json:"received"`
	} `json:"coin_history"`
}

type InfoHandler interface {
	GetUserInfo(c *gin.Context)
}

type handler struct {
	infoSrv info.Service
}

func New(
	infoSrv info.Service,
) InfoHandler {
	return &handler{
		infoSrv: infoSrv,
	}
}

func (h *handler) GetUserInfo(c *gin.Context) {
	infoData, err := h.infoSrv.GetInfo(c) //.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user info"})
		c.Error(err)
		return
	}

	response := InfoResponse{
		Coins:     infoData.Coins,
		Inventory: infoData.Inventory,
		CoinHistory: struct {
			Sent     []domain.Transaction `json:"sent"`
			Received []domain.Transaction `json:"received"`
		}{
			Sent:     infoData.CoinHistory.Sent,
			Received: infoData.CoinHistory.Received,
		},
	}

	c.JSON(http.StatusOK, response)
}
