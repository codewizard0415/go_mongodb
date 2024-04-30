package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"mongo/db"
	"mongo/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection = db.Db().Database("admin").Collection("gofirst") // mongodb client instance

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var person models.User // variable of type User in mongodb
	defer r.Body.Close()   // Close the request body ater reading is done
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Extract form values
	name := r.Form.Get("name")
	gender := r.Form.Get("gender")
	age := r.Form.Get("age")
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		http.Error(w, "Invalid age", http.StatusBadRequest)
		return
	}
	person.ID = primitive.NewObjectID()
	person.Age = ageInt
	person.Gender = gender
	person.Name = name

	insertResult, err := userCollection.InsertOne(context.TODO(), person) // insert data into mongo client
	if err != nil {
		log.Fatal(err)
	}
	person.ID = insertResult.InsertedID.(primitive.ObjectID)
	// return to response object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Inserted user with id: ", insertResult.InsertedID, "\n")
	json.NewEncoder(w).Encode(person)

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["id"] //get Parameter value as string

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(404)
	}

	var result primitive.M // an unordered representation of a BSON document which is a Map
	err = userCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}).Decode(&result)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
	}

	// Convert primitive.M to bson
	var person models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	err = bson.Unmarshal(bsonBytes, &person)
	if err != nil {
		log.Fatal(err)
	}

	// return to response object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person) // returns a json object from bson
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["id"] //get Parameter value as string

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(404)
	}

	var result primitive.M // an unordered representation of a BSON document which is a Map
	// get the document
	err = userCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}).Decode(&result)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
	}

	// Convert primitive.M to bson
	var person models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	err = bson.Unmarshal(bsonBytes, &person)
	if err != nil {
		log.Fatal(err)
	}

	// Delete the specified document
	opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison
	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}, opts)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(404)
	}

	// return to response object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted Users: ", res.DeletedCount, "\n")
	json.NewEncoder(w).Encode(person) // returns deleted object
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	//slice for multiple documents
	var results []models.User
	cur, err := userCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
	}

	//Next() gets the next document for corresponding cursor
	for cur.Next(context.TODO()) {
		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		// convert primitive.M to json
		var person models.User
		bsonBytes, err := bson.Marshal(elem)
		if err != nil {
			log.Fatal(err)
		}
		err = bson.Unmarshal(bsonBytes, &person)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, person) // add document to slice
	}
	cur.Close(context.TODO()) // close the cursor
	// return to response object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	type updateBody struct {
		Name   string `json:"name"`   //value that has to be modified
		Gender string `json:"gender"` // value that has to be modified
		Age    int    `json:"age"`    // value that has to be modified
	}

	params := mux.Vars(r)["id"] //get Parameter value as string

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(404)
	}

	// set updateBody with values from request object
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Print(e)
	}

	filter := bson.D{{Key: "_id", Value: _id}} // converting value to BSON type
	after := options.After                     // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: body.Age}, {Key: "name", Value: body.Name}, {Key: "gender", Value: body.Gender}}}}
	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	// convert primitive.M to bson
	var person models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	err = bson.Unmarshal(bsonBytes, &person)
	if err != nil {
		log.Fatal(err)
	}

	// return to response objects
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}
