package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

type UserAPI struct {
	db *mongo.Database
}

func (u *UserAPI) getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	collection := u.db.Collection("users")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(users)
}

func (u *UserAPI) getUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get("id")
	var user User
	collection := u.db.Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(user)
}

// เปลี่ยนชื่อฟังก์ชัน CreateUser เป็นเป็น public method โดยเปลี่ยนชื่อเป็น CreateUser
func (u *UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	collection := u.db.Collection("users")
	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)
}


func (u *UserAPI) updateUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get("id")
	var updatedUser User
	json.NewDecoder(r.Body).Decode(&updatedUser)
	collection := u.db.Collection("users")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"name": updatedUser.Name, "email": updatedUser.Email, "password": updatedUser.Password}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(updatedUser)
}

func (u *UserAPI) deleteUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get("id")
	collection := u.db.Collection("users")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "User deleted successfully")
}



func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://drink179531:FPUMbLQ7bolMiaau@cluster0.ign7q1r.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	defer client.Disconnect(context.Background())

	db := client.Database("test")
	userAPI := &UserAPI{db: db}

	http.HandleFunc("/users", userAPI.getUsers)
	http.HandleFunc("/users/get", userAPI.getUser)
	http.HandleFunc("/users/create", userAPI.CreateUser) // เปลี่ยนจาก userAPI.createUser เป็น userAPI.createUser
	http.HandleFunc("/users/update", userAPI.updateUser)
	http.HandleFunc("/users/delete", userAPI.deleteUser)

	http.ListenAndServe(":8080", nil)
}