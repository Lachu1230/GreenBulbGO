package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Product struct {
	Name    string `json:"name,omitempty"`
	Image   string `json:"image_src,omitempty"`
	Type    string `json:"type,omitempty"`
	Price   string `json:"price,omitempty"`
	Wattage string `json:"wattage,omitempty"`
	Usage   string `json:"usage,omitempty"`
	Desc    string `json:"desc,omitempty"`
}

type Order struct {
	ProductID string `json:"productid,omitempty"`
	UserID    string `json:"userid,omitempty"`
	Status    string `json:"status,omitempty"`
}

type ProductArray struct {
	Products []Product `json:"products,omitempty"`
}

type Energy struct {
	EnergyVal int64 `json:"energy_usage"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/register", RegisterEndPoint).Methods("POST")
	router.HandleFunc("/login", LoginEndPoint).Methods("POST")
	router.HandleFunc("/getAllProducts", GetAllProducts).Methods("GET")
	router.HandleFunc("/createOrder", CreateOrder).Methods("POST")
	router.HandleFunc("/getAllOrders", GetAllOrders).Methods("GET")
	router.HandleFunc("/getEnergyUsage", GetEnergyUsage).Methods("POST")
	fmt.Println("Server started")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func RegisterEndPoint(w http.ResponseWriter, r *http.Request) {

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://root:root@ecobulb-paysk.gcp.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	UserCollection := client.Database("GreenBulb").Collection("User")
	res, err := UserCollection.InsertOne(ctx, user)
	fmt.Println(res)
	fmt.Println(bson.M{"user": "Pavan"})
}
func LoginEndPoint(w http.ResponseWriter, r *http.Request) {
	var user User
	var result User
	_ = json.NewDecoder(r.Body).Decode(&user)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://root:root@ecobulb-paysk.gcp.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	UserCollection := client.Database("GreenBulb").Collection("User")
	ctx = context.TODO()
	filter := bson.M{"email": user.Email}
	errfind := UserCollection.FindOne(ctx, filter).Decode(&result)
	if errfind != nil {
		log.Fatal(errfind)
	}

	if result.Password == user.Password {
		json.NewEncoder(w).Encode(&result)
	} else {
		json.NewEncoder(w).Encode("Not valid")
	}

}
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://root:root@ecobulb-paysk.gcp.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	UserCollection := client.Database("GreenBulb").Collection("Products")
	cur, err := UserCollection.Find(ctx, bson.D{{}})

	var results []Product

	for cur.Next(ctx) {

		var product Product
		errProduct := cur.Decode(&product)
		if errProduct != nil {
			log.Fatal(errProduct)
		}

		results = append(results, product)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	products := ProductArray{
		Products: results,
	}

	json.NewEncoder(w).Encode(products)

}
func CreateOrder(w http.ResponseWriter, r *http.Request)  {}
func GetAllOrders(w http.ResponseWriter, r *http.Request) {}
func GetEnergyUsage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	IncandescentVal, err := strconv.ParseInt(keyVal["Incandescent"], 10, 64)
	HalogenVal, err := strconv.ParseInt(keyVal["Halogen"], 10, 64)
	CFLVal, err := strconv.ParseInt(keyVal["CFL"], 10, 64)
	LEDVal, err := strconv.ParseInt(keyVal["LED"], 10, 64)

	IncandescentWattVal, err := strconv.ParseInt(keyVal["IncandescentWatts"], 10, 64)
	HalogenWattVal, err := strconv.ParseInt(keyVal["HalogenWatts"], 10, 64)
	CFLWattVal, err := strconv.ParseInt(keyVal["CFLWatts"], 10, 64)
	LEDWattVal, err := strconv.ParseInt(keyVal["LEDWatts"], 10, 64)

	AverageVal, err := strconv.ParseInt(keyVal["Average"], 10, 64)

	wattHour := IncandescentVal*IncandescentWattVal + HalogenVal*HalogenWattVal + CFLVal*CFLWattVal + LEDVal*LEDWattVal

	whperYear := wattHour * 365 * AverageVal

	energyData := Energy{
		EnergyVal: whperYear,
	}

	json.NewEncoder(w).Encode(energyData)
}
