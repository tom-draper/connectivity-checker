package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Ping struct {
	Loss     float64   `json:"loss"`
	Response int64     `json:"response"`
	Time     time.Time `json:"time"`
}

type Data struct {
	Name  string `json:"name"`
	Pings []Ping `json:"pings"`
}

var Websites = [...]string{"pldashboard.com", "pldashboard-backend.vercel.app", "tomdraper.dev", "my-api-analytics.vercel.app", "api-analytics-server.vercel.app", "notion-courses.netlify.app", "colour-themes.netlify.app", "array-3d-viz.vercel.app", "codedevblog.netlify.app", "persona-api.vercel.app", "connectivity-checker.netlify.app"}

func getEnv(key string, production bool) string {
	if !production {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Some error occured. Err: %s", err)
		}
	}
	val := os.Getenv(key)

	return val
}

func validAddress(address string) bool {
	for _, v := range Websites {
		if v == address {
			return true
		}
	}

	return false
}

func FetchData(id string) Data {
	if !validAddress(id) {
		return Data{}
	}

	var data Data
	collection := ConnectToDatabase(true)

	filter := bson.D{{"name", id}}

	err := collection.FindOne(context.TODO(), filter).Decode(&data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	return data
}

func FetchAllData() []Data {
	var data []Data
	collection := ConnectToDatabase(true)

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem Data
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, elem)
	}

	return data
}

func UpdateDatabase(collection *mongo.Collection, address string, ping Ping) {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"name", address}}
	update := bson.D{{"$push", bson.D{{"pings", ping}}}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}

	// Remove oldest ping
	// filter = bson.D{{"name", address}}
	// update = bson.D{{"$pop", bson.D{{"pings", -1}}}}
	// _, err = collection.UpdateOne(context.TODO(), filter, update, opts)
	// if err != nil {
	// 	panic(err)
	// }
}

func ConnectToDatabase(production bool) *mongo.Collection {
	username := getEnv("USERNAME", production)
	password := getEnv("PASSWORD", production)

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://" + username + ":" + password + "@main.pvnry.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("Connectivity").Collection("Pings")
	return collection
}
