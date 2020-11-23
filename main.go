package main
import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
)
var db *sql.DB
func rollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("simple_list.html")
		if err != nil {
			log.Fatal(err)
		}
		books, err := dbGetBooks()
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, books)
	}
}
func addmaxBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("simple_max_price.html")
		if err != nil {
			log.Fatal(err)
		}
		books2, err := dbGetMaxBooks()
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, books2)
	}
}
func addOldBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("simple_old_date.html")
		if err != nil {
			log.Fatal(err)
		}
		books2, err := dbGetOldBooks()
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, books2)
	}
}
func addBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("simple_form.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		bakery_type := r.Form.Get("bakery_type")
		bakery_price := r.Form.Get("bakery_price")
		bakery_production_date := r.Form.Get("bakery_production_date")
		bakery_shelf_life := r.Form.Get("bakery_shelf_life")
		err := dbAddBook(bakery_type, bakery_price, bakery_production_date, bakery_shelf_life)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
		fmt.Println(port)
	}
	return ":" + port
}
func main() {
	err := dbConnect()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", rollHandler)
	http.HandleFunc("/add", addBookHandler)
	http.HandleFunc("/max", addmaxBookHandler)
	http.HandleFunc("/old", addOldBookHandler)
	log.Fatal(http.ListenAndServe(GetPort(), nil))
}


type Book struct{
	Bakery_type string
	Bakery_price string
	Bakery_production_date string
	Bakery_shelf_life string

}
type Book2 struct{

	Bakery_type_max string
	Bakery_price_max string
	Bakery_production_date_max string
	Bakery_shelf_life_max string
}
const (
	DB_USER = "postgres"
	DB_PASSWORD = "k1204"
	DB_NAME = "lab"
)
func dbConnect() error {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME))
	if err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS bakeryProducts ( bakery_type text,bakery_price numeric(7,2), bakery_production_date date, bakery_shelf_life date)"); err != nil {
		return err
	}
	return nil
}
func dbAddBook(bakery_type, bakery_price, bakery_production_date, bakery_shelf_life string) error {
	sqlstmt := "INSERT INTO bakeryProducts VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(sqlstmt, bakery_type, bakery_price, bakery_production_date, bakery_shelf_life)
	if err != nil {
		return err
	}
	return nil
}
func dbGetBooks() ([]Book, error) {
	var books []Book
	stmt, err := db.Prepare("SELECT bakery_type, bakery_price, bakery_production_date, bakery_shelf_life FROM bakeryProducts")
	if err != nil {
		return books, err
	}
	res, err := stmt.Query()
	if err != nil {
		return books, err
	}
	var tempBook Book
	for res.Next() {
		err = res.Scan(&tempBook.Bakery_type, &tempBook.Bakery_price, &tempBook.Bakery_production_date, &tempBook.Bakery_shelf_life)
		if err != nil {
			return books, err
		}
		books = append(books, tempBook)
	}
	return books, err
}
func dbGetMaxBooks() ([]Book2, error) {
	var books2 []Book2
	stmt, err := db.Prepare("SELECT bakery_type, bakery_price, bakery_production_date, bakery_shelf_life FROM bakeryProducts where  bakery_price = (SELECT max(bakery_price) FROM bakeryProducts)")
	if err != nil {
		return books2, err
	}
	res, err := stmt.Query()
	if err != nil {
		return books2, err
	}
	var tempBook2 Book2
	for res.Next() {
		err = res.Scan(&tempBook2.Bakery_type_max, &tempBook2.Bakery_price_max, &tempBook2.Bakery_production_date_max, &tempBook2.Bakery_shelf_life_max)
		if err != nil {
			return books2, err
		}
		books2 = append(books2, tempBook2)
	}
	return books2, err
}
func dbGetOldBooks() ([]Book2, error) {
	var books2 []Book2
	stmt, err := db.Prepare("select bakery_type, bakery_price, bakery_production_date, bakery_shelf_life FROM bakeryProducts where bakery_shelf_life < now()")
	if err != nil {
		return books2, err
	}
	res, err := stmt.Query()
	if err != nil {
		return books2, err
	}
	var tempBook2 Book2
	for res.Next() {
		err = res.Scan(&tempBook2.Bakery_type_max, &tempBook2.Bakery_price_max, &tempBook2.Bakery_production_date_max, &tempBook2.Bakery_shelf_life_max)
		if err != nil {
			return books2, err
		}
		books2 = append(books2, tempBook2)
	}
	return books2, err
}