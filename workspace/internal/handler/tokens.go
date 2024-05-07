package handler

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/T1-SD/internal/database"
	"github.com/T1-SD/internal/model"
)

func CreateToken(c *gin.Context) {
	var body model.CreateTokenRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"invalid request body": err.Error()})
		return
	}

	res, err := database.Token.InsertOne(c, body)
	if err != nil {
		c.JSON(500, gin.H{"unable to insert token": err.Error()})
		return
	}

	token := model.Token{
		ID:    res.InsertedID.(primitive.ObjectID),
		Token: body.Token,
	}

	c.JSON(201, token)
}

func GetToken(c *gin.Context) {
	var token model.Token
	err := database.Token.FindOne(c.Request.Context(), bson.D{}).Decode(&token)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, gin.H{"error": "no tokens found"})
		} else {
			c.JSON(500, gin.H{"error": "failed to retrieve token", "details": err.Error()}) // Otro error
		}
		return
	}

	c.JSON(200, gin.H{"token": token.Token})
}

func DeleteAllTokens(c *gin.Context) {
	res, err := database.Token.DeleteMany(c, bson.D{})
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to delete token", "details": err.Error()})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(404, gin.H{"message": "no tokens found to delete"})
		return
	}

	c.JSON(200, gin.H{"message": "token deleted successfully"})
}
