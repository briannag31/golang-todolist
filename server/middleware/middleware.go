package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init(){
	loadTheEnv()
	createDBInstance()
}

func loadTheEnv(){
	err := godotenv.Load(".env")
	if err!=nil{
		log.Fatal("Error loading the .env file")
	}
}

func createDBInstance(){
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("DB_COLLECTION_NAME")

	clientOptions := options.Client().ApplyURL(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err!=nil{
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("connected to mongodb")

	collection = client.Dabase(dbName).Collection(collectionName)

	fmt.Println("collection instance created")
}
func GetAllTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// whatever we get in the payload, we will encode it to JSON for the frontend
	payload := getAllTasks()
	json.NewEncoder(w).Encode(payload)
}

func CreateTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var task models.ToDoList
	// send the task as a post request to the json of whatever is in the body
	json.NewDecoder(r.Body).Decode(&task)
	createTask(task)
	json.NewEncoder(w).Encode(task)



}

func GetSingleTask(){
	
}

func TaskComplete(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	taskComplete(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func UndoTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	undoTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	deleteOneTask(params["id"])
}

func DeleteAllTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	count := deleteAllTasks()
	json.NewEncoder(w).Encode(count)

}

func getAllTasks() []primitive.M {
	// call the function and query the database
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err!=nil{
		log.Fatal(err)
	}
	var results []primitive.M
	for cur.Next(context.Background()){
		var result bson.M
		e := cur.Decode(&result)
		if e!=nil{
			log.Fatal(e)
		}
		results = append(results, result)
	}
	if err :=cur.Err(); err!=nil{
		log.Fatal(err)
	}
	cur.Close(context.Background())
	return results
}

func taskComplete(task string){
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status": true}}
	result, err:=collection.UpdateOne(context.Background(), filter, update)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("modified count: ", result.ModifiedCount)
}

func insertOneTask(task models.ToDoList){
	insertResult, err :=collection.insertOne(context.Background(), task)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Inserted a single record", insertResult.InsertedID)
}

func undoTask(task string){
	id, _:=primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status": false}}
	result, err:=collection.UpdateOne(context.Background(), filter, update)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("modified count: ", result.ModifiedCount)
}

func deleteOneTask(task string){
	id, _:=primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	d, err := collection.DeleteOne(context.Background(), filter)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("deleted item: ", d.ModifiedCount)

}

func deleteAllTasks() int64{
	d, err := collection.DeleteMany(context.Background())
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("deleted", d.DeletedCount)
	return d.DeletedCount
}