package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/initializers"
	"github.com/ppay/lib"
	"github.com/ppay/models"
)

func Register(c *gin.Context) {
	var input struct {
		Fullname string  `json:"fullname"`
		Email    string  `json:"email" binding:"required,email"`
		Password string  `json:"password" binding:"required,min=6"`
		Pin      *string `json:"pin"`
		Phone    *string `json:"phone"`
		Image    *string `json:"image"`
	}

	response := lib.NewResponse(c)

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest("Invalid input", nil)
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasher := lib.CreateHash(input.Password)

	input.Password = hasher

	// Mulai transaksi
	tx := initializers.DB.Begin()

	// Buat User
	user := models.User{
		Email:    input.Email,
		Password: input.Password,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to create user", nil)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Buat Wallet untuk User
	wallet := models.Wallet{
		UserID:  user.ID,
		Balance: 0.00, // Saldo awal
	}

	if err := tx.Create(&wallet).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to create wallet", nil)
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	// Commit transaksi
	tx.Commit()

	response.Created("Register success", nil)
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "User and wallet created successfully",
	// 	"user":    user,
	// 	"wallet":  wallet,
	// })
}
