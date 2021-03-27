package forms

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestViewHandler(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}

	owner := "owner"
	mock.ExpectQuery("select").
		WithArgs(owner).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).FromCSVString("1000"))

	r := httptest.NewRequest("GET",
		"http://127.0.0.1:8080/view/"+owner, nil)
	w := httptest.NewRecorder()
	ViewHandler(w, r, owner)
	actual := w.Body.String()
	expected := "<h1>owner client account</h1>\n" +
		"\n<div>1000</div>\n" +
		"\n<p>[<a href=\"/deposit/owner\">Deposit</a>][<a href=\"/withdraw/owner\">Withdraw</a>]</p>"

	if !strings.Contains(actual, "1000") {
		t.Errorf("Wrong response\nExpected:\n%s\nActual:\n%s", expected, actual)
	}
}

func TestCreateHandler(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}

	owner := "owner"
	balance := 1000
	mock.ExpectQuery("insert into").
		WithArgs(owner, balance)

	r := httptest.NewRequest("POST",
		"http://127.0.0.1:8080/create/"+owner, nil)
	w := httptest.NewRecorder()
	CreateHandler(w, r, owner)
}

func TestDepositHandler(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}

	owner := "owner"
	balance := 1000
	mock.ExpectQuery("select").
		WithArgs(owner).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).FromCSVString("1000"))
	mock.ExpectQuery("update").
		WithArgs(owner, balance)

	r := httptest.NewRequest("POST",
		"http://127.0.0.1:8080/deposit/"+owner, nil)
	w := httptest.NewRecorder()
	DepositHandler(w, r, owner)
}

func TestWithdrawHandler(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}

	owner := "owner"
	balance := 1000
	mock.ExpectQuery("select").
		WithArgs(owner).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).FromCSVString("1000"))
	mock.ExpectQuery("update").
		WithArgs(owner, balance)

	r := httptest.NewRequest("POST",
		"http://127.0.0.1:8080/withdraw/"+owner, nil)
	w := httptest.NewRecorder()
	WithdrawHandler(w, r, owner)
}
