package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"live-polling-backend/models"
)

type PollHandler struct {
	DB    *mongo.Database
	Redis *redis.Client
}

type CreatePollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

func (h *PollHandler) CreatePoll(c *gin.Context) {
	var req CreatePollRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Backend validation
	if len(req.Question) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Question must contain at least 3 characters",
		})
		return
	}

	if len(req.Options) < 2 || len(req.Options) > 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Poll must have between 2 and 6 options",
		})
		return
	}

	// Get authenticated user ID
	userIDString, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDString.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Create options
	pollOptions := make([]models.PollOption, 0, len(req.Options))

	for i, optionText := range req.Options {
		if len(optionText) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Option cannot be empty",
			})
			return
		}

		pollOptions = append(pollOptions, models.PollOption{
			ID:    string(rune('a' + i)),
			Text:  optionText,
			Count: 0,
		})
	}

	poll := models.Poll{
		ID:        bson.NewObjectID(),
		Question:  req.Question,
		Options:   pollOptions,
		CreatedBy: userID,
		CreatedAt: time.Now().Unix(),
	}

	_, err = h.DB.Collection("polls").InsertOne(
		c.Request.Context(),
		poll,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create poll",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Poll created successfully",
		"poll":   poll,
	})
}
func (h *PollHandler) GetPoll(c *gin.Context) {
	pollID := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(pollID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid poll ID",
		})
		return
	}

	var poll models.Poll

	err = h.DB.Collection("polls").
		FindOne(c.Request.Context(), bson.M{"_id": objectID}).
		Decode(&poll)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Poll not found",
		})
		return
	}

	// Get the latest vote counts from Redis
	redisKey := "poll:" + pollID + ":votes"

	for i := range poll.Options {
		count, err := h.Redis.HGet(
			c.Request.Context(),
			redisKey,
			poll.Options[i].ID,
		).Int64()

		if err == nil {
			poll.Options[i].Count = count
		} else {
			poll.Options[i].Count = 0
		}
	}

	c.JSON(http.StatusOK, poll)
}
func (h *PollHandler) Vote(c *gin.Context) {
	pollID := c.Param("id")

	id, err := bson.ObjectIDFromHex(pollID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid poll ID",
		})
		return
	}

	var req struct {
		OptionID string `json:"optionId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Get poll from MongoDB to validate the option.
	var poll models.Poll

	err = h.DB.Collection("polls").FindOne(
		c.Request.Context(),
		bson.M{"_id": id},
	).Decode(&poll)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Poll not found",
		})
		return
	}

	validOption := false

	for _, option := range poll.Options {
		if option.ID == req.OptionID {
			validOption = true
			break
		}
	}

	if !validOption {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid poll option",
		})
		return
	}

	// Redis stores the live vote count.
	redisKey := "poll:" + pollID + ":votes"

	count, err := h.Redis.HIncrBy(
		c.Request.Context(),
		redisKey,
		req.OptionID,
		1,
	).Result()
	// Publish the new count to all live viewers.
message := `{"optionId":"` + req.OptionID + `","count":` + fmt.Sprint(count) + `}`

err = h.Redis.Publish(
	c.Request.Context(),
	"poll:"+pollID+":updates",
	message,
).Err()

if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Vote recorded but live update failed",
	})
	return
}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not record vote",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Vote recorded",
		"optionId": req.OptionID,
		"count":    count,
	})
}