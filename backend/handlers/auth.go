package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"live-polling-backend/models"
)

type AuthHandler struct {
	DB *mongo.Database
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Name == "" || req.Email == "" || len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name, email and password of at least 6 characters are required",
		})
		return
	}

	ctx := c.Request.Context()

	collection := h.DB.Collection("users")

	var existing models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existing)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already registered",
		})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not secure password",
		})
		return
	}

	user := models.User{
		ID:       bson.NewObjectID(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	_, err = collection.InsertOne(ctx, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := c.Request.Context()

	collection := h.DB.Collection("users")

	var user models.User

	err := collection.FindOne(
		ctx,
		bson.M{"email": req.Email},
	).Decode(&user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": user.ID.Hex(),
			"exp":    time.Now().Add(24 * time.Hour).Unix(),
		},
	)

	signedToken, err := token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create authentication token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   signedToken,
		"user": gin.H{
			"id":   user.ID.Hex(),
			"name": user.Name,
			"email": user.Email,
		},
	})
}