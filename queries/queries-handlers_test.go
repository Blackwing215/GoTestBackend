package queries

import (
	"net/http/httptest"
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
		"http://127.0.0.1:8080/queries/view?name="+owner, nil)
	w := httptest.NewRecorder()
	ViewHandler(w, r)

}

func TestCreateHandler(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}

	owner := "owner"
	balance := "1000"
	mock.ExpectQuery("insert into").
		WithArgs(owner, balance)

	r := httptest.NewRequest("POST",
		"http://127.0.0.1:8080/queries/create?name="+owner+"&balance="+balance, nil)
	w := httptest.NewRecorder()
	CreateHandler(w, r)
}

func TestDepositHandler(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}

	owner := "owner"
	balance := "1000"
	mock.ExpectQuery("select").
		WithArgs(owner).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).FromCSVString("1000"))
	mock.ExpectQuery("update").
		WithArgs(owner, balance)

	r := httptest.NewRequest("POST",
		"http://127.0.0.1:8080/queries/deposit?name="+owner+"&balance="+balance, nil)
	w := httptest.NewRecorder()
	DepositHandler(w, r)
}

func TestWithdrawHandler(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	DB, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}

	owner := "owner"
	balance := "1000"
	mock.ExpectQuery("select").
		WithArgs(owner).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).FromCSVString("1000"))
	mock.ExpectQuery("update").
		WithArgs(owner, balance)

	r := httptest.NewRequest("POST",
		"http://127.0.0.1:8080/queries/withdraw?name="+owner+"&balance="+balance, nil)
	w := httptest.NewRecorder()
	WithdrawHandler(w, r)
}
