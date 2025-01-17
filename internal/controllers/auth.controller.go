package controllers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/dto"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/models"
	"github.com/ppay/lib"
)

func Register(c *gin.Context) {
	response := lib.NewResponse(c)
	var input dto.RegisterDTO

	if err := c.ShouldBind(&input); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}
	fmt.Println(input)
	if !IsValidEmail(input.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}

	if !IsValidPassword(input.Password) {
		response.BadRequest("Password must be at least 8 characters long", nil)
		return
	}

	hasher := lib.GenerateHash(input.Password)
	input.Password = hasher

	// Mulai transaksi
	tx := initializers.DB.Begin()

	// Buat User
	var user = models.User{
		Email:    input.Email,
		Password: input.Password,
	}

	// Periksa apakah email sudah digunakan
	if _, err := GetUserByEmail(input.Email); err == nil {
		response.BadRequest("Email is already registered", nil)
		return
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to create user", nil)
		return
	}

	// Buat Wallet untuk User
	var wallet = models.Wallet{
		UserID:  user.ID,
		Balance: 0.00, // Saldo awal
	}

	if err := tx.Create(&wallet).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to create user", nil)
		return
	}

	// Commit transaksi
	tx.Commit()

	response.Created("Success register user", nil)
}

func Login(c *gin.Context) {
	response := lib.NewResponse(c)
	var input dto.LoginDTO

	if err := c.ShouldBind(&input); err != nil {
		fmt.Println(err)
	}
	user, err := GetUserByEmail(input.Email)
	if err != nil {
		response.BadRequest("Invalid email or password", nil)
		return
	}

	//compare password
	if user == nil || !lib.VerifyHash(input.Password, user.Password) {
		response.BadRequest("Invalid email or password", nil)
		return
	}

	// Generate token JWT
	token, err := lib.GenerateToken(int(user.ID))
	if err != nil {
		response.InternalServerError("Failed to generate token", nil)
		return
	}

	response.Success("Login successful", gin.H{"token": token})
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// Query ke database untuk mencari user berdasarkan email
	if err := initializers.DB.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error; err != nil {
		return nil, err // Kembalikan error jika user tidak ditemukan
	}
	return &user, nil
}

func IsValidEmail(email string) bool {
	if len(email) < 5 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

func IsValidPassword(password string) bool {
	return len(password) >= 8
}
