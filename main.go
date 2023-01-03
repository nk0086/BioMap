package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"example.com/bio-map/oauth"
)

type Page struct {
	Title     string
	Body      string
	Organisms []*Organism
}

var templates = make(map[string]*template.Template)

func main() {
	port := "8080"

	// メインページ：マップへの生物分布の表示、検索
	templates["index"] = loadTemplate("index")
	http.HandleFunc("/", handleIndex)

	// 生物の情報を登録
	templates["register"] = loadTemplate("register")
	http.HandleFunc("/register", handleRegister)

	// 生物の情報を一覧表示
	templates["list"] = loadTemplate("list")
	http.HandleFunc("/list", handleList)

	http.HandleFunc("/login", oauth.LoginHandler)
	http.HandleFunc("/about", oauth.CallbackHandler)
	http.ListenAndServe(":"+port, nil)
}

func loadTemplate(name string) *template.Template {
	tmpl, err := template.ParseFiles(
		"template/"+name+".html",
		"template/_header.html",
		"template/_footer.html",
	)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	return tmpl
}

func handleList(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDatabase("database.db")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer db.Close()
	organisms, err := selectAllFromTable(db)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	page := Page{
		Title:     "BioMap",
		Body:      "This is a test",
		Organisms: organisms,
	}

	if err := templates["list"].Execute(w, page); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "BioMap",
		Body:  "This is a test",
	}

	if err := templates["index"].Execute(w, page); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	//生物情報を登録する
	if r.Method == "POST" {
		if err := r.ParseMultipartForm(1024); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		name := r.FormValue("name")
		latitude, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
		longitude, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
		if err != nil {
			log.Printf("failed to parse float: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer file.Close()

		image, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		db, err := connectToDatabase("database.db")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		defer db.Close()

		_, err = insertIntoTable(db, name, image, latitude, longitude)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/register", http.StatusFound)

	} else {
		page := Page{
			Title: "BioMap",
			Body:  "This is a test",
		}

		if err := templates["register"].Execute(w, page); err != nil {
			log.Printf("failed to execute template: %v", err)
		}
	}

}
