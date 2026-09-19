package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type RealtimeHandler struct {
	Redis *redis.Client
}

func (h *RealtimeHandler) PollUpdates(c *gin.Context) {
	pollID := c.Param("id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	channel := "poll:" + pollID + ":updates"

	pubsub := h.Redis.Subscribe(
		c.Request.Context(),
		channel,
	)

	defer pubsub.Close()

	// Confirm Redis subscription.
	_, err = pubsub.Receive(c.Request.Context())
	if err != nil {
		return
	}

	ch := pubsub.Channel()

	for {
		msg, ok := <-ch

		if !ok {
			return
		}

		err := conn.WriteMessage(
			websocket.TextMessage,
			[]byte(msg.Payload),
		)

		if err != nil {
			return
		}
	}
}