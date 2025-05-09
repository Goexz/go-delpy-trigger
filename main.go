package main

import (
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
)

type Message struct {
    ID           string    `json:"id"`
    Content      string    `json:"message"`
    ReceiverUUID string    `json:"receiverUuid"`
    Timestamp    time.Time `json:"timestamp"`
}

var messages = make(map[string][]Message)

func init() {
    log.SetFormatter(&log.JSONFormatter{})
    log.SetLevel(log.InfoLevel)
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.WithField("port", port).Info("Starting server")

    r := gin.Default()
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    r.LoadHTMLGlob("templates/*")
    r.Static("/static", "./static")
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })

    r.GET("/generate-uuid", func(c *gin.Context) {
        newUUID := uuid.New().String()
        c.JSON(http.StatusOK, gin.H{"uuid": newUUID})
    })

    r.POST("/messages", func(c *gin.Context) {
        var newMessage Message
        if err := c.BindJSON(&newMessage); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        newMessage.ID = uuid.New().String()
        newMessage.Timestamp = time.Now()

        messages[newMessage.ReceiverUUID] = append(messages[newMessage.ReceiverUUID], newMessage)
        c.JSON(http.StatusOK, newMessage)
    })

    r.GET("/messages/:uuid", func(c *gin.Context) {
        uuid := c.Param("uuid")
        if userMessages, exists := messages[uuid]; exists {
            c.JSON(http.StatusOK, userMessages)
        } else {
            c.JSON(http.StatusOK, []Message{})
        }
    })

    r.DELETE("/messages/:id", func(c *gin.Context) {
        messageID := c.Param("id")

        for uuid, userMessages := range messages {
            for i, msg := range userMessages {
                if msg.ID == messageID {
                    messages[uuid] = append(userMessages[:i], userMessages[i+1:]...)
                    c.JSON(http.StatusOK, gin.H{"message": "ลบข้อความสำเร็จ"})
                    return
                }
            }
        }

        c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อความที่ต้องการลบ"})
    })

    r.Run(":" + port)
}
