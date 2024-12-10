
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User model
type User struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Username      string `gorm:"unique" json:"username"`
	Email         string `gorm:"unique" json:"email"`
	Password      string `json:"password"`
	EmailVerified bool   `json:"email_verified"` // Field to track if the email is verified
}

// Database connection
var db *gorm.DB

func initDB() {
	var err error
	dsn := "host=localhost user=postgres password=Pass dbname=postgres port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database!")
	}

	// Migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("Failed to migrate database schema!")
	}
}

func main() {
	// Initialize the database
	initDB()

	// Set up the Gin router
	router := gin.Default()

	// Routes
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.PUT("/update-password", updatePassword)
	router.POST("/reset-password", resetPassword)
	router.POST("/verify-email", verifyEmail)
	router.POST("/verify-password", verifyPassword)

	// Start the server
	router.Run(":8080")
}

// Signup handler
func signup(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if user already exists
	var existing User
	if err := db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Create a new user
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Login handler
func login(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Authenticate user
	var user User
	if err := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Update Password handler
func updatePassword(c *gin.Context) {
	var input struct {
		Email           string `json:"email"`
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.NewPassword != input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password and confirm password do not match"})
		return
	}

	// Find the user by email
	var user User
	if err := db.Where("email = ? AND password = ?", input.Email, input.OldPassword).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Update the password
	user.Password = input.NewPassword
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// Reset Password handler
func resetPassword(c *gin.Context) {
	var input struct {
		Email           string `json:"email"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.NewPassword != input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password and confirm password do not match"})
		return
	}

	// Find the user by email
	var user User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update the password
	user.Password = input.NewPassword
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// Verify Email handler
func verifyEmail(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
		// Typically, you'd also want to send a verification token to the user's email
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find the user by email
	var user User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Mark the user's email as verified
	user.EmailVerified = true
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// Verify Password handler
func verifyPassword(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find the user by email and check password
	var user User
	if err := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password verified successfully"})
}
