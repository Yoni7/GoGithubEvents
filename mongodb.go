package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Events struct {
	ID          	primitive.ObjectID 	`bson:"_id"`
	Name  			string             	`bson:"name"`
	Counter 		uint64           	`bson:"counter"`
}

type Emails struct {
	ID				primitive.ObjectID 	`bson:"_id"`
	Email  			string            	`bson:"email"`
}


type General struct {
	ID           	primitive.ObjectID 	`bson:"_id"`
	General 		string 				`bson:"general"`
	Actors []struct{					
		ID			primitive.ObjectID 	`bson:"_id"`
		Name 		string 				`bson:"name"`
	}									`bson:"actors"`
	Repos []struct{
		ID			primitive.ObjectID 	`bson:"_id"`
		Name 		string 				`bson:"name"`
		Url 		string 				`bson:"url"`
	}									`bson:"repos"`
}

const ACTORS_LIMIT = -50
const REPOS_LIMIT = -20

var MONGO_CLIENT mongo.Client

func ConnectToMongoDB() {
	mongoUser := GetEnv("MONGO_USER", "yoni")
	mongoPassword := GetEnv("MONGO_PASSWORD", "DwXzxWKzu1Dc72EX")

	var connectionStr = "mongodb+srv://" + mongoUser + ":" + mongoPassword + "@cluster0.3ymvqrs.mongodb.net/?retryWrites=true&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionStr))
	if err != nil {
		panic(err)
	}
	MONGO_CLIENT = *client
	fmt.Println("Connected to MongoDB")
}


/* Events Types */
func UpdateEventType(eventName string) {
	fmt.Printf("updateEventType: %v\n", eventName)

	if eventName == "" {
		fmt.Println("Error: missing event name")
		return
	}
	eventsCollection := MONGO_CLIENT.Database("github").Collection("events")

	filter := bson.M{"name": eventName}
	update := bson.M{
		"$set": bson.M{
			"name": eventName,
		},
		"$inc": bson.M{
			"counter": 1,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := eventsCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		fmt.Printf("Error: failed to update event type %v. error %v", eventName, err)
	}
}


func GetEventsDocs() []Events {
	fmt.Println("GetEventsDocs")
	eventsCollection := MONGO_CLIENT.Database("github").Collection("events")

	var events []Events
	eventsCursor, err := eventsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Printf("Error: failed to get event type. error %v", err)
	}

	eventsCursor.All(context.TODO(), &events)
	return events
}


/* Events Emails */
func UpdateEventEmail(eventEmail string) {
	fmt.Printf("UpdateEventEmail: %v\n", eventEmail)
	emailsCollection := MONGO_CLIENT.Database("github").Collection("emails")

	filter := bson.M{"email": eventEmail}
	update := bson.M{
		"$set": bson.M{
			"email": eventEmail,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := emailsCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		fmt.Printf("Error: failed to update event email %v. error %v", eventEmail, err)
	}
}

func GetEmailsDocs() []Emails {
	fmt.Println("GetEmailsDocs")
	emailsCollection := MONGO_CLIENT.Database("github").Collection("emails")

	var emails []Emails
	emailsCursor, err := emailsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Printf("Error: failed to get event email. error %v", err)
	}

	emailsCursor.All(context.TODO(), &emails)
	return emails
}

/* Actors */
func UpdateEventActor(actorName string) {
	fmt.Printf("UpdateEventActor: %v\n", actorName)
	generalCollection := MONGO_CLIENT.Database("github").Collection("general")

	filter := bson.M{"general": "general"}
	update := bson.D{{"$push", bson.D{{"actors", bson.D{{"$sort", bson.D{{"_id", 1}}}, {"$each", bson.A{bson.D{{"name", actorName}, {"_id", primitive.NewObjectID()}}}}, {"$slice", ACTORS_LIMIT}}}}}}

	_, err := generalCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error: failed to update event actor %v. error %v", actorName, err)
	}
}

func GetActorsDocs() General {
	fmt.Println("GetActorsDocs")
	generalCollection := MONGO_CLIENT.Database("github").Collection("general")

	var general General
	filter := bson.M{"general": "general"}
	project := bson.D{{ "actors", 1 }}
	opts := options.FindOne().SetProjection(project)
	err := generalCollection.FindOne(context.TODO(), filter, opts).Decode(&general)
	if err != nil {
		fmt.Printf("Error: failed to get actors. error %v", err)
	}

	return general
}

/* Repos */
func UpdateEventRepo(repoName string, repoUrl string) {
	fmt.Printf("UpdateEventRepo: %v | %v\n", repoName, repoUrl)

	if repoName == "" {
		fmt.Println("Error: missing repo name")
		return
	}

	if repoUrl == "" {
		fmt.Println("Error: missing repo url")
		return
	}

	generalCollection := MONGO_CLIENT.Database("github").Collection("general")

	filter := bson.M{"general": "general"}
	update := bson.D{{"$push", bson.D{{"repos", bson.D{{"$sort", bson.D{{"_id", 1}}}, {"$each", bson.A{bson.D{{"name", repoName}, {"url", repoUrl}, {"_id", primitive.NewObjectID()}}}}, {"$slice", REPOS_LIMIT}}}}}}

	_, err := generalCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error: failed to update event repo %v | %v. error %v", repoName, repoUrl, err)
	}
}

func GetRepoDocs() General {
	fmt.Println("GetRepoDocs")
	generalCollection := MONGO_CLIENT.Database("github").Collection("general")

	var general General
	filter := bson.M{"general": "general"}
	project := bson.D{{ "repos", 1 }}
	opts := options.FindOne().SetProjection(project)
	err := generalCollection.FindOne(context.TODO(), filter, opts).Decode(&general)
	if err != nil {
		fmt.Printf("Error: failed to get repos. error %v", err)
	}

	return general
}