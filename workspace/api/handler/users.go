package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/T1-SD-Grupo7/database"
	"github.com/T1-SD-Grupo7/model"
)

func GetUsers(c *gin.Context) {
	cursor, err := database.Users.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(500, gin.H{"unable to fetch users": err.Error()})
		return
	}

	var users []model.User
	if err = cursor.All(c, &users); err != nil {
		c.JSON(500, gin.H{"unable to fetch users": err.Error()})
		return
	}

	c.JSON(200, users)
}

func CreateUser(c *gin.Context) {
	var body model.RegisterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"invalid request body": err.Error()})
		return
	}

	// verifying mail
    existingUserByEmail := database.Users.FindOne(c, bson.D{{Key: "email", Value: body.Email}})
    if existingUserByEmail.Err() == nil {
        c.JSON(400, gin.H{"error": "El correo electrónico ya está en uso"})
        return
    }

    // verifying rut
    existingUserByRut := database.Users.FindOne(c, bson.D{{Key: "rut", Value: body.Rut}})
    if existingUserByRut.Err() == nil {
        c.JSON(400, gin.H{"error": "El RUT ya está en uso"})
        return
    }


	res, err := database.Users.InsertOne(c, body)
	if err != nil {
		c.JSON(500, gin.H{"unable to insert user": err.Error()})
		return
	}

	user := model.User{
		ID:        res.InsertedID.(primitive.ObjectID),
		Name:      body.Name,
		Last_name: body.Last_name,
		Rut:       body.Rut,
		Email:     body.Email,
		Password:  body.Password,
	}

	c.JSON(201, user)
}

func Login(c *gin.Context) {
	var loginReq model.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, gin.H{"invalid request body": err.Error()})
		return
	}

	var user model.User
	err := database.Users.FindOne(c, bson.D{{Key: "email", Value: loginReq.Email}}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"invalid email or password": err.Error()})
		return
	}

	if user.Password != loginReq.Password {
		c.JSON(404, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(200, gin.H{"data": gin.H{
		"id":        user.ID.Hex(),
		"name":      user.Name,
		"last_name": user.Last_name,
		"rut":       user.Rut,
		"email":     user.Email,
	}})
}

func GetUserById(c *gin.Context) {
	_id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"invalid id": err.Error()})
		return
	}

	result := database.Users.FindOne(c, bson.D{{Key: "_id", Value: _id}})
	user := model.User{}
	err = result.Decode(&user)
	if err != nil {
		c.JSON(500, gin.H{"unable to find user": err.Error()})
		return
	}

	c.JSON(200, user)
}

func UpdateUserById(c *gin.Context) {
	_id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"invalid id": err.Error()})
		return
	}

	var updateData model.RegisterRequest
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

	collection := database.Users
	result, err := collection.UpdateByID(c, _id, update)
	if err != nil {
		c.JSON(500, gin.H{"failed to update user": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(404, gin.H{"error": "no user found with given ID"})
		return
	}

	c.JSON(200, gin.H{"message": "user updated successfully"})
}

func DeleteUserById(c *gin.Context) {
	_id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"invalid id": err.Error()})
		return
	}

	res, err := database.Users.DeleteOne(c, bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		c.JSON(500, gin.H{"failed to delete user": err.Error()})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(400, gin.H{"error": "no user found with given ID"})
		return
	}

	c.JSON(200, gin.H{"message": "user deleted successfully"})
}
