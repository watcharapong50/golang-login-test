package main

import "net/http"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "golang.org/x/crypto/bcrypt"

var db *sql.DB
var err error

func main() {
    db,err = sql.Open("mysql","root:@/testing")
    if err != nil {
         panic(err.Error())
    }
     defer db.Close()
     err = db.Ping()
     if err != nil {
          panic(err.Error())
    }
    http.HandleFunc("/signup",signupPage)
    http.HandleFunc("/login",loginPage)
    http.HandleFunc("/", homePage)
    http.ListenAndServe(":8080", nil)

}
func homePage(res http.ResponseWriter, req *http.Request) {
    http.ServeFile(res, req,"index.html")
}
func loginPage(res http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        http.ServeFile(res, req,"login.html")
        return
    }
    username := req.FormValue("username")
    password := req.FormValue("password")
	   var dbName string
	   var dbPass string

	   err := db.QueryRow("SELECT username,password FROM users WHERE username=? ", username).Scan(&dbName,&dbPass)

  	if err != nil {
  		http.Redirect(res, req,"/login",301)
  		return
  	}

  	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password))
  	if err != nil {
  		http.Redirect(res, req, "/login",301)
  		return
  	}
	   res.Write([]byte("Hello : " + dbName))
}

func signupPage(res http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        http.ServeFile(res, req, "signup.html")
        return
    }
   username := req.FormValue("username")
   password := req.FormValue("password")
   var user string
   err := db.QueryRow("SELECT username FROM users WHERE username=?",username).Scan(&user)
   switch {
      case err == sql.ErrNoRows:
          hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
          if err != nil {
              http.Error(res, "Server error, unable to create your account.", 500)
              return
          }
        _, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
        if err != nil {
            http.Error(res, "Server error, unable to create your account.", 500)
            return
        }
        res.Write([]byte("User created!"))
        return
      case err != nil:
          http.Error(res, "Server error, unable to create your account.", 500)
          return
    default:
        http.Redirect(res, req,"/", 301)
    }
}
