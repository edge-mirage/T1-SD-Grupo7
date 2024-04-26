package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/T1-SD-Grupo7/database"
	"github.com/T1-SD-Grupo7/model"
)

func GetClients(c *gin.Context) {
	cursor, err := database.Clients.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(500, gin.H{"unable to fetch clients": err.Error()})
		return
	}

	var clients []model.Client
	if err = cursor.All(c, &clients); err != nil {
		c.JSON(500, gin.H{"unable to fetch clients": err.Error()})
		return
	}

	c.JSON(200, clients)
}

func CreateClient(c *gin.Context) {
	var body model.CreateClientRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"invalid request body": err.Error()})
		return
	}

	res, err := database.Clients.InsertOne(c, body)
	if err != nil {
		c.JSON(500, gin.H{"unable to insert client": err.Error()})
		return
	}

	client := model.Client{
		ID:        res.InsertedID.(primitive.ObjectID),
		Name:      body.Name,
		Last_name: body.Last_name,
		Rut:       body.Rut,
		Email:     body.Email,
	}

	c.JSON(201, client)
}

func GetClientById(c *gin.Context) {
	_id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"invalid id": err.Error()})
		return
	}

	result := database.Clients.FindOne(c, bson.D{{Key: "_id", Value: _id}})
	client := model.Client{}
	err = result.Decode(&client)
	if err != nil {
		c.JSON(500, gin.H{"unable to find client": err.Error()})
		return
	}

	c.JSON(200, client)
}

func UpdateClientById(c *gin.Context) {
	_id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"invalid id": err.Error()})
		return
	}

	var updateData model.CreateClientRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(500, gin.H{"invalid request body": err.Error()})
		return
	}

	updateFields := bson.M{}
	if updateData.Name != "" {
		updateFields["name"] = updateData.Name
	}
	if updateData.Last_name != "" {
		updateFields["last_name"] = updateData.Last_name
	}
	if updateData.Rut != "" {
		updateFields["rut"] = updateData.Rut
	}
	if updateData.Email != "" {
		updateFields["email"] = updateData.Email
	}

	update := bson.M{
		"$set": updateFields,
	}

	if len(updateFields) == 0 {
		c.JSON(400, gin.H{"error": "no fields to update"})
		return
	}

	collection := database.Clients
	result, err := collection.UpdateByID(c, _id, update)
	if err != nil {
		c.JSON(500, gin.H{"failed to update client": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(404, gin.H{"error": "no client found with given ID"})
		return
	}

	c.JSON(200, gin.H{"message": "client updated successfully"})
}

func DeleteClientById(c *gin.Context) {
	_id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"invalid id": err.Error()})
		return
	}

	res, err := database.Clients.DeleteOne(c, bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		c.JSON(500, gin.H{"failed to delete client": err.Error()})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(400, gin.H{"error": "no client found with given ID"})
		return
	}

	c.JSON(200, gin.H{"message": "client deleted successfully"})
}
