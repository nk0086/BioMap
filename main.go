package main

import (
	"html/template"
	"log"
	"net/http"

	"example.com/bio-map/oauth"
)

type Page struct {
	Title string
	Body  string
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
	page := Page{
		Title: "BioMap",
		Body:  "This is a test",
	}

	if err := templates["register"].Execute(w, page); err != nil {
		log.Printf("failed to execute template: %v", err)
	}

}
