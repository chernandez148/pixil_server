package controllers

import (
	"pixi/api/config"
	"pixi/api/models" // Import the models package
	"pixi/api/utils"

	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go" // Import JWT package
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Login handles user login and generates a JWT token
func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Bind the incoming JSON request data
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate that both username and password are provided
	if loginData.Username == "" || loginData.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Get the database connection
	db := config.GetDB()

	// Find the user by username and preload their posts without user data in the posts
	var user models.User
	if err := db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
		// This ensures the 'User' inside 'Posts' is not loaded
		return db.Omit("User")
	}).Preload("Likes").Preload("Saves").Preload("Posts.Comments").Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the provided password with the stored hashed password
	if err := user.CheckPassword(loginData.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Manually clean up the User field from Posts if it's still present
	for i := range user.Posts {
		user.Posts[i].User = nil // Set User to nil to avoid any nested user data
	}

	// Generate JWT token if credentials are correct
	token, exp, err := generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Return the generated token and user details
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
		"token":   token,
		"exp":     exp, // Returning the expiration time
	})
}

// generateJWT creates a JWT token for the authenticated user
func generateJWT(user models.User) (string, int64, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"sub":      user.ID,                               // Subject: User ID
		"username": user.Username,                         // Username
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Expiration time (3 days)
	}

	// Define the JWT secret key (using an environment variable for security)
	secretKey := []byte(utils.GetEnv("JWT_SECRET_KEY", "ksdjfgvaseyr7465794623pvbdlugsfvbf8o7igw4y65tv9p8q264vbe7rw"))

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", 0, err
	}

	// Return the token and expiration time
	expiration := claims["exp"].(int64)
	return tokenString, expiration, nil
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var newUser models.User
	// Bind the request JSON to the newUser object
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Save the new user to the database
	db := config.GetDB() // Get the database connection
	if err := newUser.Save(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
		return
	}

	// Respond with the created user
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": newUser})
}

// UpdateUser updates an existing user's details
func UpdateUser(c *gin.Context) {
	db := config.GetDB() // Get the database connection
	var user models.User
	userID := c.Param("id") // Get the user ID from the URL params

	// Bind the input JSON data
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find the user by ID
	var existingUser models.User
	if err := db.First(&existingUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the password was updated and hash it
	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		existingUser.Password = user.Password // Update the password field with the hashed password
	}

	// Update the other fields (username, email, etc.) if they are provided
	if user.FullName != "" {
		existingUser.FullName = user.FullName
	}
	if user.Username != "" {
		existingUser.Username = user.Username
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}

	// Save the updated user data
	if err := db.Save(&existingUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, existingUser)
}

// GetUsers retrieves all users
func GetUsers(c *gin.Context) {
	// Initialize the database connection
	db := config.GetDB()

	// Declare a variable to hold the list of users
	var users []models.User

	// Fetch all users from the database
	if err := db.Find(&users).Error; err != nil {
		// If an error occurs, return an error response
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Return the list of users in the response
	c.JSON(http.StatusOK, users)
}

// GetUser retrieves a user by ID
func GetUser(c *gin.Context) {
	db := config.GetDB() // Get the database connection
	var user models.User
	userID := c.Param("id") // Get the user ID from the URL params

	// Find the user by ID
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user by ID
func DeleteUser(c *gin.Context) {
	db := config.GetDB() // Get the database connection
	var user models.User
	userID := c.Param("id") // Get the user ID from the URL params

	// Find the user by ID
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete the user
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
