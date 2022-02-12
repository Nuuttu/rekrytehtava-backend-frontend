package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/darahayes/go-boom"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Coffee struct {
	Id         int      `json:"Id"`
	Name       *string  `json:"Name"`
	Weight     *int     `json:"Weight"`
	Price      *float64 `json:"Price"` // mikä valuutta :> ?
	RoastLevel *int     `json:"RoastLevel"`
}

var CoffeeList []Coffee

// ENDPOINT lista
func apiIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Api endpoints available: \n/coffeelist \n/coffee/{id} \n/coffee/add POST \n/coffee/delete/{id} DELETE \n/becoffeelist")
	fmt.Println("Endpoint Hit: apiIndex")
}

// ENDPOINT palauttaa koko listan
func backendCoffeeList(w http.ResponseWriter, r *http.Request) {
	coffees := "Kahvit:\n"
	for _, c := range CoffeeList {
		coffees = coffees + fmt.Sprintf("Name: %s Hinta: %.2f  Paino: %d Paahtoaste: %d Id: %d\n", *c.Name, *c.Price, *c.Weight, *c.RoastLevel, *&c.Id)
	}
	fmt.Fprintf(w, "%s", coffees)
}

// ENDPOINT palauttaa koko listan
func returnCoffees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Endpoint Hit: All")
	json.NewEncoder(w).Encode(CoffeeList)
}

// ENDPOINT palauttaa yhden kahvin
func returnSingleCoffee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	key := stringToint(vars["id"])

	for _, Coffee := range CoffeeList {
		if Coffee.Id == key {
			json.NewEncoder(w).Encode(Coffee)
		}
	}
}

// ENDPOINT Lisää listan
func createNewCoffee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newCoffee Coffee
	json.Unmarshal(reqBody, &newCoffee)

	// Tsekataan, että kaikki on kunnossa pyynnössä
	err := verifyCoffee(newCoffee)
	if err != nil {
		boom.BadData(w, err)
	} else {
		err = writeNewCoffee(w, newCoffee)
		if err != nil {
			boom.Internal(w, "Error while trying to create new Coffee")
		}
	}
}

// funktio createNewCoffee ENDPOINTille
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
	fmt.Printf("Coffee added to the Coffee List: %s\n", *newCoffee.Name)
	updateCoffeeList()
	return nil
}

// ENDPOINT - poista kahvi
func deleteCoffee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	id := stringToint(vars["id"])
	if !containsCoffeeById(CoffeeList, id) {
		boom.BadRequest(w, "No coffee by that ID")
	} else {
		var deletedCoffee Coffee
		for index, coffee := range CoffeeList {
			if coffee.Id == id {
				deletedCoffee = coffee
				CoffeeList = append(CoffeeList[:index], CoffeeList[index+1:]...)
			}
		}
		err := setCoffeeList(CoffeeList)
		if err != nil {
			fmt.Println("Error happened", err)
		}
		json.NewEncoder(w).Encode(deletedCoffee)
		getCoffeeList()
		fmt.Printf("Coffee deleted: %s\n", *deletedCoffee.Name)
	}
}

// ROUTER - mux - Tämä jäi vähän heikkoon malliin. Melkein toimii, mutta yliarvioin edistymiseni.
// Ensimmäinen golang backendi, niin jäi hieman hakemiseksi :) Ei ihan täytä kaikkia turvallisuusvaatimuksia.
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", apiIndex)
	myRouter.HandleFunc("/becoffeelist", backendCoffeeList)
	myRouter.HandleFunc("/coffeelist", returnCoffees).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/coffee/{id}", returnSingleCoffee).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/coffee/add", createNewCoffee).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/coffee/delete/{id}", deleteCoffee).Methods("DELETE", "OPTIONS")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin: *"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Fatal(http.ListenAndServe(":10000", toLowerCase(handlers.CORS(originsOk, headersOk, methodsOk)(myRouter))))
}

// Tiedostofunktio - päivittää ohjelman sisäisen muuttujan
func setCoffeeList(coffees []Coffee) error {
	file, err := json.MarshalIndent(coffees, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("CoffeeList.json", file, 0777)
	if err != nil {
		return err
	}
	return nil
}

// Tiedostofunktio - päivittää ohjelman sisäisen muuttujan
func getCoffeeList() []Coffee {
	coffeeJson, _ := ioutil.ReadFile("CoffeeList.json")
	var cList []Coffee
	json.Unmarshal(coffeeJson, &cList)
	return cList
}

// MAIN func
func main() {
	fmt.Println("Setting up a server on http://localhost:10000")
	updateCoffeeList()
	handleRequests()
}

/*









	HELPER FUNCTIONS
*/

func updateCoffeeList() {
	CoffeeList = getCoffeeList()
}

// Middleware, jotta saataisiin poistettua case sensitiviteetti
func toLowerCase(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
} // https://groups.google.com/g/gorilla-web/c/rpwSDVOxkyc

func containsCoffeeById(cl []Coffee, u int) bool {
	for _, c := range cl {
		if c.Id == u {
			return true
		}
	}
	return false
}

func setId() int {
	var idList []int
	for _, c := range CoffeeList {
		idList = append(idList, c.Id)
	}
	newId := rand.Intn(10000)
	for containsId(idList, newId) {
		newId = rand.Intn(10000)
	}
	return newId
}

func containsId(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func stringToint(str string) int {
	int, err := strconv.Atoi(str)
	if err != nil {
		log.Println("error happened", err)
	}
	return int
}

func verifyCoffee(coffee Coffee) error {
	if coffee.Name == nil {
		return errors.New("Bad Name")
	}
	if *coffee.Weight < 0 || coffee.Weight == nil {
		return errors.New("Bad Weight")
	}
	if *coffee.Price < 0 || coffee.Price == nil {
		return errors.New("Bad Price")
	}
	if *coffee.RoastLevel < 1 || *coffee.RoastLevel > 5 || coffee.RoastLevel == nil {
		return errors.New("Bad Roastlevel")
	}
	return nil
}
