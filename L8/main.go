package main

import (
	"context"
	"errors"
	"fmt"
	"l8es/domain"
	"l8es/espackage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

var (
	es  *espackage.ES
	err error
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	//http
	mux := http.NewServeMux()
	mux.HandleFunc("/add", AddHandler)
	mux.HandleFunc("/search", SearchHandler)
	mux.HandleFunc("/getall", SearchGetAllHandler)
	mux.HandleFunc("/searchallf", SearchAllFields)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	//es stuff
	es, err = espackage.New()
	if err != nil {
		log.Println("error on es new: ", err)
		return
	}

	esInfo, err := es.ES.Info()
	if err != nil {
		log.Println("error on es info: ", err)
		return
	}
	log.Printf("es created; info: %v\n", esInfo)
	esInfo.Body.Close()

	// stop stuff
	wg := sync.WaitGroup{}
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer signal.Stop(sigterm)

		<-sigterm
		log.Println("terminating: via signal")

		err := s.Shutdown(context.Background())
		if err != nil {
			log.Println("error on http shutdown: ", err)
		}
		log.Println("http server stopped")
	}()

	log.Println("http server started on port 8080")
	err = s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println("error on http server: ", err)
	}
	log.Println("http server ListenAndServe stopped")

	wg.Wait()
	log.Println("main stop")
}

// curl -d "name=name1&price=10.12&description=desc1" http://localhost:8080/add
func AddHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	price := r.FormValue("price")
	description := r.FormValue("description")
	// strconv.FormatFloat(price, 'f', 2, 64)
	f, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Println("error on parse float: ", err)
		w.Write([]byte(err.Error()))
		return
	}
	product := domain.Product{
		Name:        name,
		Price:       f,
		Description: description,
	}
	err = es.Add(&product)
	if err != nil {
		log.Println("error on product add: ", err)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("add: name: %s, price: %s", name, price)
	w.Write([]byte("add ok"))
}

// curl -d "fieldName=name&value=na" http://localhost:8080/search
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	fieldName := r.FormValue("fieldName")
	value := r.FormValue("value")
	log.Printf("SearchHandler: fieldName: %s, price: %s", fieldName, value)

	res := es.Search(fieldName, value, "prefix")
	log.Printf("SearchHandler: res: %v\n", len(res))
	if len(res) == 0 {
		w.Write([]byte("SearchHandler: no results"))
		return
	}

	str := ""
	for _, v := range res {
		str += Print(v)
		log.Println(v)
	}

	// w.Write([]byte("search ok"))
	w.Write([]byte(str))
}

func Print(item map[string]any) string {
	name := item["name"]
	price := 0.0
	if item["price"] != nil {
		price = item["price"].(float64)
	}
	description := ""
	if item["description"] != nil {
		description = item["description"].(string)
	}

	return fmt.Sprintf("name: %s, price: %f, description: %s\n", name, price, description)
}

// curl -d "fieldName=name&value=na" http://localhost:8080/getall
func SearchGetAllHandler(w http.ResponseWriter, r *http.Request) {
	res := es.SearchGetAll()
	log.Printf("SearchGetAll: res: %v\n", len(res))

	if len(res) == 0 {
		w.Write([]byte("SearchGetAll: no results"))
		return
	}

	str := ""
	for _, v := range res {
		str += Print(v)
		log.Println(v)
	}
	w.Write([]byte(str))
}

// curl -d "name=name1&price=10.12&desc=desc1" http://localhost:8080/searchallf
func SearchAllFields(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	price := r.FormValue("price")
	desc := r.FormValue("desc")

	res := es.SearchAllFields(name, price, desc)
	log.Printf("SearchAllFields: res: %v\n", len(res))

	if len(res) == 0 {
		w.Write([]byte("SearchAllFields: no results"))
		return
	}

	str := ""
	for _, v := range res {
		str += Print(v)
		log.Println(v)
	}
	w.Write([]byte(str))
}
