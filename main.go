package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testbackend/forms"
	"testbackend/queries"

	_ "github.com/go-sql-driver/mysql"
)

type Configuration struct {
	Login    string
	Password string
	Db       string
}

type Client struct {
	Owner   string
	Balance int
}

var (
	db     *sql.DB
	config Configuration
)

func init() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Unable to read configuration file!")
		log.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Unable to read configuration file!")
		log.Fatal(err)
	}
}

func main() {
	var err error
	db, err = sql.Open("mysql", config.Login+":"+config.Password+"@/"+config.Db)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	forms.DB = db
	http.HandleFunc("/view/", forms.MakeHandler(forms.ViewHandler))
	http.HandleFunc("/create/", forms.MakeHandler(forms.CreateHandler))
	http.HandleFunc("/update/", forms.MakeHandler(forms.UpdateHandler))
	http.HandleFunc("/deposit/", forms.MakeHandler(forms.DepositHandler))
	http.HandleFunc("/withdraw/", forms.MakeHandler(forms.WithdrawHandler))

	queries.DB = db
	http.HandleFunc("/queries/view", queries.MakeHandler(queries.ViewHandler))
	http.HandleFunc("/queries/create", queries.MakeHandler(queries.CreateHandler))
	http.HandleFunc("/queries/deposit", queries.MakeHandler(queries.DepositHandler))
	http.HandleFunc("/queries/withdraw", queries.MakeHandler(queries.WithdrawHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
	db.Close()
}
