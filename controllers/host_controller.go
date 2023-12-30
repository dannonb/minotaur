package controllers

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/dannonb/go-network-monitor/config"
	"github.com/dannonb/go-network-monitor/helpers"
	"github.com/dannonb/go-network-monitor/models"
	"github.com/dannonb/go-network-monitor/responses"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var hostCollection *mongo.Collection = config.GetCollection(config.DB, "hosts")

func AddHost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var host models.Host
		var existingHost models.Host
		userId := c.Param("userId")
		userObjId, _ := primitive.ObjectIDFromHex(userId)

		defer cancel()
		if err := c.BindJSON(&host); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(host)
		defer cancel()
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := hostCollection.CountDocuments(ctx, bson.M{"hostname": host.Hostname})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for existing subscriber"})
			return
		}

		if count > 0 {
			hostCollection.FindOne(ctx, bson.M{"hostname": host.Hostname}).Decode(&existingHost)
		
			alreadyRegistered := slices.Contains(existingHost.Subscriber_ids, userObjId)
			if alreadyRegistered == true {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user has already registered this host"})
				return
			} else {
				subscriber_ids := append(existingHost.Subscriber_ids, userObjId)
				updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

				update := bson.M{"updated_at": updated_at, "subscriber_ids": subscriber_ids}
				result, err := hostCollection.UpdateOne(ctx, bson.M{"host_id": existingHost.Id.Hex()}, bson.M{"$set": update})
				if err != nil {
					c.JSON(http.StatusInternalServerError, responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
					return
				}

				var updatedHost models.Host
				if result.MatchedCount == 1 {
					err := hostCollection.FindOne(ctx, bson.M{"host_id": existingHost.Id.Hex()}).Decode(&updatedHost)
					if err != nil {
						c.JSON(http.StatusInternalServerError, responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
						return
					}
				}

				c.JSON(http.StatusOK, updatedHost)
			}

		} else {
			host.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			host.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			host.Id = primitive.NewObjectID()
			host.Host_id = host.Id.Hex()
			host.Subscriber_ids = append(host.Subscriber_ids, userObjId)

			resultInsertionNumber, insertErr := hostCollection.InsertOne(ctx, host)
			if insertErr != nil {
				msg := fmt.Sprintf("Host item was not created")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			defer cancel()

			err = helpers.AddSingleHostToRedis(host.Hostname)
			if err != nil {
				msg := fmt.Sprintf("Host was not added to cache, please wait until server reset to view your host's status")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}

			c.JSON(http.StatusOK, resultInsertionNumber)
		}
	}
}

func GetAllHosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var hosts []models.Host
		defer cancel()

		results, err := hostCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.Host
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			hosts = append(hosts, singleUser)
		}

		c.JSON(http.StatusOK,
			responses.APIResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": hosts}},
		)
	}
}
