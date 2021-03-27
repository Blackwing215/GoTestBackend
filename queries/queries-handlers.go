package queries

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Client struct {
	Owner   string
	Balance int
}

var (
	validPath = regexp.MustCompile("^/queries/(create|update|withdraw|deposit|view)") ///([a-zA-Z0-9]+)$")
	DB        *sql.DB
)

func MakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	names, ok := r.URL.Query()["name"]
	if !ok {
		fmt.Fprintln(w, "Wrong request!")
		return
	}

	output := "owner: balance"

	for _, owner := range names {
		c := Client{Owner: owner}
		row := DB.QueryRow("select balance from go_test_schema.clients where owner=?", owner)
		err := row.Scan(&c.Balance)
		switch err {
		case sql.ErrNoRows:
			fmt.Fprintln(w, owner+" not found")
		case nil:
			output += "\n" + owner + ":" + fmt.Sprint(c.Balance)
		default:
			fmt.Fprintln(w, "Error while getting "+owner+"'s data")
		}
	}
	fmt.Fprintln(w, output)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	names, ok := r.URL.Query()["name"]
	if !ok {
		fmt.Fprintln(w, "Wrong request!")
		return
	}

	balances, ok := r.URL.Query()["balance"]
	if !ok {
		balances = make([]string, len(names))
	}

	if len(names) != len(balances) {
		fmt.Fprintln(w, "Wrong request! Different ammount of owners and sums")
		return
	}

	output := "/request/view?name="

	for i, owner := range names {
		sum, err := strconv.Atoi(balances[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			continue
		}
		_, err = DB.Exec("insert into go_test_schema.clients (owner, balance) values (?, ?)",
			sum, owner)
		if err != nil {
			fmt.Fprintln(w, "Create error for "+owner+"\n"+err.Error())
			continue
		}
		output += owner + ","
	}

	http.Redirect(w, r, output[:len(output)-1], http.StatusFound)
}

func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	names, ok := r.URL.Query()["name"]
	if !ok {
		fmt.Fprintln(w, "Wrong request!")
		return
	}

	balances, ok := r.URL.Query()["balance"]
	if !ok {
		fmt.Fprintln(w, "Wrong request!")
		return
	}

	if len(names) != len(balances) {
		fmt.Fprintln(w, "Wrong request! Different ammount of owners and sums")
		return
	}

	output := "/request/view?name="

	for i, owner := range names {
		sum, err := strconv.Atoi(balances[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			continue
		}
		_, err = DB.Exec("update go_test_schema.clients set balance=balance-? where owner=?",
			sum, owner)
		if err != nil {
			fmt.Fprintln(w, "Withdraw error for "+owner+"\n"+err.Error())
			continue
		}
		output += owner + ","
	}

	http.Redirect(w, r, output[:len(output)-1], http.StatusFound)
}

func DepositHandler(w http.ResponseWriter, r *http.Request) {
	names, ok := r.URL.Query()["name"]
	if !ok {
		fmt.Fprintln(w, "Wrong request!")
		return
	}

	balances, ok := r.URL.Query()["balance"]
	if !ok {
		fmt.Fprintln(w, "Wrong request!")
		return
	}

	if len(names) != len(balances) {
		fmt.Fprintln(w, "Wrong request! Different ammount of owners and sums")
		return
	}

	output := "/request/view?name="

	for i, owner := range names {
		sum, err := strconv.Atoi(balances[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			continue
		}
		_, err = DB.Exec("update go_test_schema.clients set balance=balance+? where owner=?",
			sum, owner)
		if err != nil {
			fmt.Fprintln(w, "Deposit error for "+owner+"\n"+err.Error())
			continue
		}
		output += owner + ","
	}

	http.Redirect(w, r, output[:len(output)-1], http.StatusFound)
}
