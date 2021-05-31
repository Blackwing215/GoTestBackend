package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"testbackend/forms"
	"testbackend/queries"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Client struct {
	Owner   string
	Balance int
}

var db *sql.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	var err error
	dburl := os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@tcp(" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + ")/" +
		os.Getenv("DB_NAME")
	db, err = sql.Open(os.Getenv("DB_DRIVER"), dburl)
	//db, err = sql.Open("mysql", config.Login+":"+config.Password+"@/"+config.Db)
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
