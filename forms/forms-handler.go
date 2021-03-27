package forms

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Client struct {
	Owner   string
	Balance int
}

var (
	validPath = regexp.MustCompile("^/(create|update|withdraw|deposit|view)/([a-zA-Z0-9]+)$")
	DB        *sql.DB
)

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

var templates = template.Must(template.ParseFiles(
	"create.html",
	"deposit.html",
	"withdraw.html",
	"view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, c Client) {
	err := templates.ExecuteTemplate(w, tmpl+".html", c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request, owner string) {
	c := Client{Owner: owner}
	row := DB.QueryRow("select balance from go_test_schema.clients where owner=?", owner)
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

func CreateHandler(w http.ResponseWriter, r *http.Request, owner string) {
	c := Client{owner, 0}
	renderTemplate(w, "create", c)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request, owner string) {
	sum := r.FormValue("sum")
	submit := r.FormValue("operation")
	i, err := strconv.Atoi(sum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch submit {
	case "Deposit":
		_, err = DB.Exec("update go_test_schema.clients set balance=balance+? where owner=?",
			i, owner)
	case "Withdraw":
		_, err = DB.Exec("update go_test_schema.clients set balance=balance-? where owner=?",
			i, owner)
	case "Create":
		_, err = DB.Exec("insert into go_test_schema.clients (owner, balance) values (?, ?)",
			owner, i)
	}
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/view/"+owner, http.StatusFound)
}

func WithdrawHandler(w http.ResponseWriter, r *http.Request, owner string) {
	row := DB.QueryRow("select balance from go_test_schema.clients where owner=?", owner)
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

func DepositHandler(w http.ResponseWriter, r *http.Request, owner string) {
	row := DB.QueryRow("select balance from go_test_schema.clients where owner=?", owner)
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
