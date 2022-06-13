package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection URI
const uri = "mongodb://127.0.0.1:27017/infinity"

var ctx = context.Background()

func main() {
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	col := client.Database("infinity").Collection("bots")

	for x := range time.Tick(10 * time.Second) {
		fmt.Println(x)

		cur, _ := col.Find(ctx, bson.M{})

		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var a struct {
				ID   string `bson:"botID"`
				Date any    `bson:"date"`
			}

			cur.Decode(&a)

			fmt.Println(a.Date, a.ID, reflect.TypeOf(a.Date))

			if d, ok := a.Date.(primitive.DateTime); ok {
				fmt.Println("Found bad bot", a.ID)
				proper := d.Time().UnixMilli()
				fmt.Println("Corrected time:", proper)
				client.Database("infinity").Collection("bots").UpdateOne(ctx, bson.M{"botID": a.ID}, bson.M{"$set": bson.M{"date": proper}})
			}
		}
	}
}
