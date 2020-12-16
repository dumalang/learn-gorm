package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/darahayes/go-boom"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"net/http"
)

var db *gorm.DB
var err error

type Product struct {
	ID     uint64          `json:"id" gorm:"primaryKey;type:bigint(20) AUTO_INCREMENT;"`
	Uuid   uuid.UUID       `json:"uuid" type:uuid;`
	Code   string          `json:"code"`
	Name   string          `json:"name"`
	Price  decimal.Decimal `json:"price" sql:"type:decimal(15,2)"`
	Price2 decimal.Decimal `json:"price2" gorm:"type:numeric" sql:"type:decimal(15,2)"`
	Stock  uint64          `json:"stock"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	fmt.Println("Hello world")
	db, err = gorm.Open("mysql", "root@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True")

	if err != nil {
		fmt.Println(err)
		log.Println("Connection Failed to open")
	} else {
		log.Println("Connection Established")
	}

	db.Debug().DropTable(&Product{})
	db.Debug().AutoMigrate(&Product{})

	handleRequest()

}

func handleRequest() {
	log.Println("Start the deevelopment server at http://127.0.0.1:9999")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/products", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", showProduct).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	myRouter.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9999", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	payload, _ := ioutil.ReadAll(r.Body)
	var product Product

	json.Unmarshal(payload, &product)

	product.Uuid = uuid.Must(uuid.NewV4())
	product.Price2 = product.Price

	db.Debug().Create(&product)

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Success create product",
	}

	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}

	db.Find(&products)

	res := Result{
		Code:    200,
		Data:    products,
		Message: "Success get products",
	}

	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func showProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	product := Product{}

	db.First(&product, productID)

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Success get product",
	}

	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	product := Product{}

	updateResult := db.First(&product, productID)

	log.Println(updateResult.RowsAffected) // returns found records count

	if errors.Is(updateResult.Error, gorm.ErrRecordNotFound) {
		//notFoundResult, _ := json.Marshal(Result{
		//	Code:    http.StatusNotFound,
		//	Message: "Product not found",
		//})
		//http.Error(w, string(notFoundResult), http.StatusNotFound)

		boom.NotFound(w, "Product not found bro.")
	}

	payload, _ := ioutil.ReadAll(r.Body)
	var productData Product

	json.Unmarshal(payload, &productData)

	db.Model(&product).Updates(productData).Debug()

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Success update product",
	}

	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	product := Product{}

	updateResult := db.First(&product, productID)

	log.Println(updateResult.RowsAffected) // returns found records count

	if errors.Is(updateResult.Error, gorm.ErrRecordNotFound) {
		//notFoundResult, _ := json.Marshal(Result{
		//	Code:    http.StatusNotFound,
		//	Message: "Product not found",
		//})
		//http.Error(w, string(notFoundResult), http.StatusNotFound)

		boom.NotFound(w, "Product not found bro.")
	}

	db.Delete(&product).Debug()

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Success update product",
	}

	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}
