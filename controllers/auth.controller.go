package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pilinux/argon2"
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

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Buat Wallet untuk User
	wallet := models.Wallet{
		UserID:  user.ID,
		Balance: 0.00, // Saldo awal
	}

	if err := tx.Create(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	// Commit transaksi
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "User and wallet created successfully",
		"user":    user,
		"wallet":  wallet,
	})
}

type Response struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Results  any    `json:"results,omitempty"`
}

func Login(c *gin.Context) {
	var formUser  models.User

	c.ShouldBind(&formUser)

	var user models.User
	if err := initializers.DB.Where("email = ? AND is_deleted = ?", formUser.Email, false).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	godotenv.Load()
	var SECRET_KEY = os.Getenv("SECRET_KEY")
	match, err := argon2.ComparePasswordAndHash(formUser.Password, SECRET_KEY, user.Password)
	if err != nil || !match {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Wrong email or password",
		})
		return
	}
}