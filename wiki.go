package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// Page : plaintext wiki page
type Page struct {
	Title string
	Body  []byte
}

var (
	// html template caching will immediately halt execution if files not found
	templates = template.Must(template.ParseFiles("edit.html", "view.html"))
	validPath = regexp.MustCompile("^/(edit|save|view)/([A-Za-z0-9]+)$")
)

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Nothing to see here!")
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, loadErr := loadPage(title)
	if loadErr != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, loadErr := loadPage(title)
	if loadErr != nil {
		p = &Page{
			Title: title,
			Body:  nil,
		}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{
		Title: title,
		Body:  []byte(body),
	}
	saveErr := p.save()
	if saveErr != nil {
		http.Error(w, saveErr.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		match := validPath.FindStringSubmatch(request.URL.Path)
		if match == nil {
			http.NotFound(writer, request)
			return
		}
		fn(writer, request, match[2])
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
