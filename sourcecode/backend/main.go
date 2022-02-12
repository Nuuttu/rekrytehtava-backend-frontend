package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Coffee struct {
	Id         uint64   `json:"Id"` // Hieman ylilyöntiä hinnassa
	Name       *string  `json:"Name"`
	Weight     *int     `json:"Weight"`
	Price      *float64 `json:"Price"` // mikä valuutta :> ?
	RoastLevel *int     `json:"RoastLevel"`
}

var CoffeeList []Coffee

func stringToUint64(str string) uint64 {
	uint, err := strconv.ParseUint(str, 0, 64)
	if err != nil {
		log.Println("error happened", err)
	}
	return uint
}

func apiIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Api endpoints available: \n/coffeelist \n/coffee/{id} \n/coffee POST \n/coffee/{id} DELETE")
	fmt.Println("Endpoint Hit: apiIndex")
}

func backendCoffeeList(w http.ResponseWriter, r *http.Request) {
	coffees := "Kahvit:\n"
	for _, c := range CoffeeList {
		coffees = coffees + fmt.Sprintf("Name: %s Hinta: %.2f  Paino: %d Paahtoaste: %d Id: %d\n", *c.Name, *c.Price, *c.Weight, *c.RoastLevel, *&c.Id)
	}
	fmt.Fprintf(w, "%s", coffees)
}

func returnAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All")
	json.NewEncoder(w).Encode(CoffeeList)
}

func returnSingleCoffee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := stringToUint64(vars["id"])

	for _, Coffee := range CoffeeList {
		if Coffee.Id == key {
			json.NewEncoder(w).Encode(Coffee)
		}
	}
}

func createNewCoffee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Adding new coffee")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var newCoffee Coffee
	json.Unmarshal(reqBody, &newCoffee)

	// Tsekataan, että kaikki on kunnossa pyynnössä

	err := writeNewCoffee(w, newCoffee)
	if err != nil {
		JSONError(w, err, 3)
	}

}

func writeNewCoffee(w http.ResponseWriter, newCoffee Coffee) error {
	updateCoffeeList()
	// Luodaan uusi ID uudelle oliolle
	newCoffee.Id = setId()
	newCoffeeList := append(CoffeeList, newCoffee)

	file, err := json.MarshalIndent(newCoffeeList, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("CoffeeList.json", file, 0777)
	if err != nil {
		return err
	}
	json.NewEncoder(w).Encode(newCoffee)
	fmt.Println("Coffee added to the Coffee List: ", newCoffee)
	updateCoffeeList()
	return nil
}

func deleteCoffee(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := stringToUint64(vars["id"])

	for index, Coffee := range CoffeeList {
		if Coffee.Id == id {
			CoffeeList = append(CoffeeList[:index], CoffeeList[index+1:]...)
		}
	}

}

func setId() uint64 {
	var idList []uint64
	for _, c := range CoffeeList {
		idList = append(idList, c.Id)
	}
	newId := rand.Uint64()
	for containsId(idList, newId) {
		newId = rand.Uint64()
	}
	return newId
}

func containsId(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// https://stackoverflow.com/questions/59763852/can-you-return-json-in-golang-http-error
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", apiIndex)
	myRouter.HandleFunc("/becoffeelist", backendCoffeeList)
	myRouter.HandleFunc("/coffeelist", returnAll)
	myRouter.HandleFunc("/coffee/{id}", returnSingleCoffee)
	myRouter.HandleFunc("/coffee", createNewCoffee).Methods("POST")
	myRouter.HandleFunc("/coffee/{id}", deleteCoffee).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":10000", toLowerCase(myRouter)))
}

// https://groups.google.com/g/gorilla-web/c/rpwSDVOxkyc
func toLowerCase(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

/*
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/route", Handler)

	log.Fatal(http.ListenAndServe("localhost:8000", LowerCaseURI(r)))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello\n"))
}

func LowerCaseURI(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
*/

func updateCoffeeList() {
	CoffeeList = getCoffeeList()
}

func getCoffeeList() []Coffee {
	coffeeJson, _ := ioutil.ReadFile("CoffeeList.json")
	var cList []Coffee
	json.Unmarshal(coffeeJson, &cList)
	return cList
}

func main() {
	updateCoffeeList()
	handleRequests()
}
