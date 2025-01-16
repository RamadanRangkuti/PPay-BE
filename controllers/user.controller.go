package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilinux/argon2"
	"github.com/ppay/initializers"
	"github.com/ppay/models"
)

type UserResponse struct {
	ID        uint    `json:"id"`
	Fullname  string  `json:"fullname"`
	Email     string  `json:"email"`
	Pin       *string `json:"pin,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Image     *string `json:"image,omitempty"`
	IsDeleted bool    `json:"is_deleted"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func GetMd5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Create User and Wallet
func CreateUser(c *gin.Context) {
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

	if input.Password != "" {
		input.Password, _ = argon2.CreateHash(input.Password, "", argon2.DefaultParams)
	}

	// Mulai transaksi
	tx := initializers.DB.Begin()

	// Buat User
	user := models.User{
		Fullname: input.Fullname,
		Email:    input.Email,
		Password: input.Password,
		Pin:      input.Pin,
		Phone:    input.Phone,
		Image:    input.Image,
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

// Get All Users
func GetUsers(c *gin.Context) {
	var users []models.User
	search := c.Query("search")

	query := initializers.DB.Where("is_deleted = ?", false)
	if search != "" {
		query = query.Where("fullname ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := initializers.DB.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Fullname string  `json:"fullname"`
		Email    string  `json:"email" binding:"omitempty,email"`
		Password string  `json:"password" binding:"omitempty,min=6"`
		Pin      *string `json:"pin"`
		Phone    *string `json:"phone"`
		Image    *string `json:"image"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields
	if input.Fullname != "" {
		user.Fullname = input.Fullname
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Password != "" {
		user.Password = input.Password
	}
	if input.Pin != nil {
		user.Pin = input.Pin
	}
	if input.Phone != nil {
		user.Phone = input.Phone
	}
	if input.Image != nil {
		user.Image = input.Image
	}

	// Save the user
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create response with UserResponse struct
	response := UserResponse{
		ID:        user.ID,
		Fullname:  user.Fullname,
		Email:     user.Email,
		Pin:       user.Pin,
		Phone:     user.Phone,
		Image:     user.Image,
		IsDeleted: user.IsDeleted,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"), // Format waktu jika perlu
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"), // Format waktu jika perlu
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": response})
}

// Delete User
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.IsDeleted = true
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
