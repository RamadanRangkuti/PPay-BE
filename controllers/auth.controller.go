package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
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
}


func Login(c *gin.Context) {
	var formUser  models.User

	if err := c.ShouldBind(&formUser) ; err != nil{
		fmt.Println(err)
	}
	
	response := lib.NewResponse(c)
	godotenv.Load()

	var user models.User
	fmt.Println(formUser.Email)
	if err := initializers.DB.Where("email = ? AND is_deleted = ?", formUser.Email, false).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	var SECRET_KEY = os.Getenv("SECRET_KEY")
	// fmt.Println(user.Password)
	match, err := argon2.ComparePasswordAndHash(formUser.Password, SECRET_KEY, user.Password)
	// fmt.Println(match)
	if err != nil || !match {
		response.BadRequest("Invalid email or password", nil)
		return
	}

	godotenv.Load()

	var JWT_SECRET []byte = []byte(GetMd5Hash(os.Getenv("JWT_SECRET")))

	signer, _ := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: JWT_SECRET}, (nil))
	baseInfo := jwt.Claims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	payload := struct {
		UserId int `json:"userId"`
	}{
		UserId: int(user.ID),
	}

	token, _ := jwt.Signed(signer).Claims(baseInfo).Claims(payload).Serialize()

	tok := models.Token{
		Token: token,
	}

	response.Success("login success", tok)
}