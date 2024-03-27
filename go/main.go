package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	hosts      = "catalogdb:27017"
	database   = "db"
	username   = ""
	password   = ""
	collection = "catalogs"
)

type Product struct {
	Name    string `json:"name"`
	Company string `json:"company"`
}

type MongoStore struct {
	session *mgo.Session
}

var mongoStore = MongoStore{}

func main() {

	fmt.Println("Server Starting.... ")
	session := initialiseMongo()
	mongoStore.session = session

	fmt.Println("Mongo DB session init .... ")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/catalogs", catalogsGetHandler).Methods("GET")
	router.HandleFunc("/catalogs", catalogsPostHandler).Methods("POST")

	fmt.Println("Handler init done .... ")

	log.Fatal(http.ListenAndServe(":9090", router))

	fmt.Println("Server Started")

}

func initialiseMongo() (session *mgo.Session) {

	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}

	return

}

func catalogsGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET Request Received")

	col := mongoStore.session.DB(database).C(collection)

	results := []Product{}
	col.Find(bson.M{"name": bson.RegEx{"", ""}}).All(&results)
	jsonString, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(jsonString))

}

func catalogsPostHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("POST Request Received")

	col := mongoStore.session.DB(database).C(collection)

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	var _product Product
	err = json.Unmarshal(b, &_product)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = col.Insert(_product)
	if err != nil {
		panic(err)
	}

	jsonString, err := json.Marshal(_product)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//Set content-type http header
	w.Header().Set("content-type", "application/json")

	//Send back data as response
	w.Write(jsonString)

}
