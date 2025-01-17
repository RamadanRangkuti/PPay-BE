package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/models"
)

func Topup(c *gin.Context) {
	var input struct {
		Amount          float64 `json:"amount" binding:"required"`
		PaymentMethodID uint    `json:"payment_method_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}

	transaction := models.Transaction{
		Amount:          input.Amount,
		TransactionType: "top_up",
		TopupTransactions: []models.TopupTransaction{
			{
				PaymentMethodID: input.PaymentMethodID,
			},
		},
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Top-up transaction created successfully",
		"transaction": transaction,
	})
}
