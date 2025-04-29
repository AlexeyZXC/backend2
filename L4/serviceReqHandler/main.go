package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"serviceReqHandler/red"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// version = "0.0.0"

	db         *sql.DB
	measurable = red.MeasurableHandler
	router     = mux.NewRouter()
	web        = http.Server{
		Handler: router,
		Addr:    ":8081",
	}
)

func init() {
	router.HandleFunc("/entities", measurable(ListEntitiesHandler)).Methods(http.MethodGet)
	router.HandleFunc("/entity", measurable(AddEntityHandler)).Methods(http.MethodPost)

	var err error
	db, err = sql.Open("mysql", "root:test@tcp(mysql:3306)/test")
	// db, err = sql.Open("mysql", "root:test@tcp(acl:3306)/test")

	if err != nil {
		panic(err)
	}
	log.Println("Connected to database ok")
}

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9091", nil); err != http.ErrServerClosed {
			panic(fmt.Errorf("error on listen and serve: %v", err))
		}
	}()

	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}

// curl -d "token=admin_secret_token&data=abc&id=1" http://localhost:8081/entity

const sqlInsertEntity = `INSERT INTO entities(id, data) VALUES (?, ?)`

func AddEntityHandler(w http.ResponseWriter, r *http.Request) {
	defer log.Printf("AddEntityHandler end\n")
	log.Printf("AddEntityHandler start\n")

	res, err := http.Get(fmt.Sprintf("http://acl:8082/identity?token=%s", r.FormValue("token")))
	log.Printf("AddEntityHandler get done")

	switch {
	case err != nil:
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(fmt.Sprintf("AddEntityHandler: get error: %v \n", err)))
		return
	case res.StatusCode != http.StatusOK:
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("AddEntityHandler: get res.StatusCode: %v \n", res.StatusCode)))
		return
	}
	res.Body.Close()
	_, err = db.Exec(sqlInsertEntity, r.FormValue("id"), r.FormValue("data"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("AddEntityHandler: Exec err: %v \n", err)))
		return
	}
}

const sqlSelectEntities = `
SELECT id, data FROM entities
`

type ListEntityItemResponse struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

// curl http://localhost:8081/entities
func ListEntitiesHandler(w http.ResponseWriter, r *http.Request) {
	defer log.Printf("ListEntitiesHandler end\n")
	log.Printf("ListEntitiesHandler start\n")

	rr, err := db.Query(sqlSelectEntities)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("ListEntitiesHandler: Query err: %v \n", err)))
		return
	}
	defer rr.Close()
	var ii = []*ListEntityItemResponse{}
	for rr.Next() {
		i := &ListEntityItemResponse{}
		err = rr.Scan(&i.Id, &i.Data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("ListEntitiesHandler: scan err: %v \n", err)))
			return
		}
		ii = append(ii, i)
	}
	bb, err := json.Marshal(ii)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("ListEntitiesHandler: Marshal err: %v \n", err)))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("ListEntitiesHandler: w.write err: %v \n", err)))
		return
	}
}
