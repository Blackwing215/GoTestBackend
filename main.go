package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"text/template"

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
	validPath = regexp.MustCompile("^/(create|update|withdraw|deposit|view)/([a-zA-Z0-9]+)$")
	db        *sql.DB
	config    Configuration
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

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

var templates = template.Must(template.ParseFiles("create.html", "deposit.html", "withdraw.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, c Client) {
	err := templates.ExecuteTemplate(w, tmpl+".html", c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	var err error
	//db, err = sql.Open("mysql", "root:rootroot@/go_test_schema")
	db, err = sql.Open("mysql", config.Login+":"+config.Password+"@/"+config.Db)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/create/", makeHandler(createHandler))
	http.HandleFunc("/update/", makeHandler(updateHandler))
	http.HandleFunc("/withdraw/", makeHandler(withdrawHandler))
	http.HandleFunc("/deposit/", makeHandler(depositHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
	db.Close()
}

func viewHandler(w http.ResponseWriter, r *http.Request, owner string) {
	c := Client{Owner: owner}
	row := db.QueryRow("select balance from go_test_schema.clients where owner=?", owner)
	err := row.Scan(&c.Balance)
	switch err {
	case sql.ErrNoRows:
		http.Redirect(w, r, "/create/"+owner, http.StatusFound)
	case nil:
		renderTemplate(w, "view", c)
	default:
		log.Fatal(err)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request, owner string) {
	c := Client{owner, 0}
	renderTemplate(w, "create", c)
}

func updateHandler(w http.ResponseWriter, r *http.Request, owner string) {
	sum := r.FormValue("sum")
	submit := r.FormValue("operation")
	i, err := strconv.Atoi(sum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch submit {
	case "Deposit":
		_, err = db.Exec("update go_test_schema.clients set balance=balance+? where owner=?",
			i, owner)
	case "Withdraw":
		_, err = db.Exec("update go_test_schema.clients set balance=balance-? where owner=?",
			i, owner)
	case "Create":
		_, err = db.Exec("insert into go_test_schema.clients (owner, balance) values (?, ?)",
			owner, i)
	}
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/view/"+owner, http.StatusFound)
}

func withdrawHandler(w http.ResponseWriter, r *http.Request, owner string) {
	row := db.QueryRow("select balance from go_test_schema.clients where owner=?", owner)
	c := Client{Owner: owner}
	err := row.Scan(&c.Balance)
	switch err {
	case sql.ErrNoRows:
		http.Redirect(w, r, "/create/"+owner, http.StatusFound)
	case nil:
		renderTemplate(w, "withdraw", c)
	default:
		log.Fatal(err)
	}
}

func depositHandler(w http.ResponseWriter, r *http.Request, owner string) {
	row := db.QueryRow("select balance from go_test_schema.clients where owner=?", owner)
	c := Client{Owner: owner}
	err := row.Scan(&c.Balance)
	switch err {
	case sql.ErrNoRows:
		http.Redirect(w, r, "/create/"+owner, http.StatusFound)
	case nil:
		renderTemplate(w, "deposit", c)
	default:
		log.Fatal(err)
	}
}
