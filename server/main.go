package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"server/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Product struct to represent a product in MongoDB
// Product struct to represent a product in MongoDB
type Product struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"` // Include MongoDB's ObjectID
	Name  string             `json:"name" bson:"name"`
	Price float64            `json:"price" bson:"price"`
}

func main() {
	// MongoDB URI (change if needed for remote MongoDB or authentication)

	clientOptions := options.Client().ApplyURI("mongodb://mongo-db:27017")

	// Create MongoDB client
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping the MongoDB server to ensure the connection is successful
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	} else {
		fmt.Println("Successfully connected to MongoDB")
	}

	// Choose a database and collection (change to your database and collection names)
	db := client.Database("shopweb")
	collection := db.Collection("products")

	// Set up Gin router
	r := gin.Default()

	// Simple Ping route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// MongoDB connection test route

	// GET route to fetch data from MongoDB
	r.GET("/api/v1/products", func(c *gin.Context) {
		var results []Product

		cursor, err := collection.Find(context.TODO(), bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error fetching data from MongoDB",
				"error":   err.Error(),
			})
			return
		}
		defer cursor.Close(context.TODO())

		for cursor.Next(context.TODO()) {
			var product Product
			if err := cursor.Decode(&product); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error decoding document",
					"error":   err.Error(),
				})
				return
			}
			results = append(results, product)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Cursor iteration error",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": results,
		})
	})

	r.POST("/api/v1/products", func(c *gin.Context) {
		var newProduct Product
		if err := c.ShouldBindJSON(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request data",
				"error":   err.Error(),
			})
			return
		}

		// Create a new ObjectID for the product
		newProduct.ID = primitive.NewObjectID()

		// Insert the product into the MongoDB collection
		_, err := collection.InsertOne(context.TODO(), newProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error inserting product into MongoDB",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Product inserted successfully",
			"data":    newProduct, // Return the full product with ID
		})
	})

	r.GET("/api/v1/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid ID format",
			})
			return
		}

		var product Product
		err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Product not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": product,
		})
	})

	r.DELETE("/api/v1/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid ID format",
			})
			return
		}

		_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error deleting product",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Product deleted successfully",
		})
	})

	r.PUT("/api/v1/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid ID format",
			})
			return
		}

		var updatedProduct Product
		if err := c.ShouldBindJSON(&updatedProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request data",
				"error":   err.Error(),
			})
			return
		}

		update := bson.M{"$set": bson.M{
			"name":  updatedProduct.Name,
			"price": updatedProduct.Price,
		}}

		_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error updating product",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Product updated successfully",
		})
	})

	// Add CORS middleware
	r.Use(middleware.CORS())
	// Start the server
	r.Run(":8080") // The server will run on localhost:8080
}
