package main
import (
  "github.com/gorilla/mux"
  "database/sql"
  _"github.com/go-sql-driver/mysql"
  "net/http"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate
type Book struct {
  ID string `json:"id"`
  Name string `json:"name" validate:"required"`
  Author string `json:"author" validate:"required"`
  Published_at string `json:"published_at"`
}
var db *sql.DB
var err error
func main() {
db, err = sql.Open("mysql", "root:@/golang")
  if err != nil {
    panic(err.Error())
  }
  defer db.Close()
  router := mux.NewRouter()
  router.HandleFunc("/api/v1/books", getbooks).Methods("GET")
  router.HandleFunc("/api/v1/books", createbook).Methods("POST")
  router.HandleFunc("/api/v1/books/{id}", getbook).Methods("GET")
  router.HandleFunc("/api/v1/books/{id}", updatebook).Methods("PUT")
  router.HandleFunc("/api/v1/books/{id}", deletebook).Methods("DELETE") 
  http.ListenAndServe(":8000", router)
}
func getbooks(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  var books []Book
  result, err := db.Query("SELECT id, name, author, published_at from books")
  if err != nil {
    panic(err.Error())
  }
  defer result.Close()
  for result.Next() {
    var book Book
    err := result.Scan(&book.ID, &book.Name, &book.Author, &book.Published_at)
    if err != nil {
      panic(err.Error())
    }
    books = append(books, book)
  }
  json.NewEncoder(w).Encode(books)
} 
func createbook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  stmt, err := db.Prepare("INSERT INTO books(name,author) VALUES(?,?)")
  if err != nil {
    panic(err.Error())
  }
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err.Error())
  }
  res := Book{}
  keyVal := make(map[string]string)
  json.Unmarshal(body, &keyVal)
  json.Unmarshal(body, &res)
  name := keyVal["name"]
  author := keyVal["author"]
  v := validator.New()
 err = v.Struct(res)
  for _, e := range err.(validator.ValidationErrors){
    fmt.Println(e)
    panic(err.Error())
  }
  _, err = stmt.Exec(name,author)
  if err != nil {
    panic(err.Error())
  }
  fmt.Fprintf(w, "New book was created")
}

func getbook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  result, err := db.Query("SELECT id, name, author, published_at FROM books WHERE id = ?", params["id"])
  if err != nil {
    panic(err.Error())
  }
  defer result.Close()
  var book Book
  for result.Next() {
    err := result.Scan(&book.ID, &book.Name, &book.Author, &book.Published_at)
    if err != nil {
      panic(err.Error())
    }
  }
  json.NewEncoder(w).Encode(book)
}

func updatebook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  stmt, err := db.Prepare("UPDATE books SET name = ?, author = ? WHERE id = ?")
  if err != nil {
    panic(err.Error())
  }
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err.Error())
  }
  keyVal := make(map[string]string)
  json.Unmarshal(body, &keyVal)
  newName := keyVal["name"]
  newAuthor := keyVal["author"]
  _, err = stmt.Exec(newName, newAuthor, params["id"])
  if err != nil {
    panic(err.Error())
  }
  fmt.Fprintf(w, "book with ID = %s was updated", params["id"])
}

func deletebook(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  stmt, err := db.Prepare("DELETE FROM books WHERE id = ?")
  if err != nil {
    panic(err.Error())
  }
  _, err = stmt.Exec(params["id"])
  if err != nil {
    panic(err.Error())
  }
  fmt.Fprintf(w, "book with ID = %s was deleted", params["id"])
} 
