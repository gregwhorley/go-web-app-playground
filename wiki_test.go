package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func genericErrorHandler(t *testing.T, e error) {
	if e != nil {
		t.Errorf("test failed -- %q", e)
	}
}

func contentErrorHandler(t *testing.T, s0 string, s1 string) {
	if s0 != s1 {
		t.Errorf("test failed -- content does not match\n %q \n %q", s0, s1)
	}
}

func TestWikiSave(t *testing.T) {
	page := &Page{
		Title: "PageTitle",
		Body:  []byte("this is the page body"),
	}
	err := page.save()
	genericErrorHandler(t, err)
	content, err2 := ioutil.ReadFile(page.Title + ".txt")
	genericErrorHandler(t, err2)
	stringContent := string(content)
	pageContent := string(page.Body)
	contentErrorHandler(t, stringContent, pageContent)
}

func TestWikiLoadPage(t *testing.T) {
	page := &Page{
		Title: "LoadPage",
		Body:  []byte("testing load page"),
	}
	page.save()
	page2, err := loadPage(page.Title)
	genericErrorHandler(t, err)
	expectedContent := string(page.Body)
	loadPageContent := string(page2.Body)
	contentErrorHandler(t, expectedContent, loadPageContent)
}

func TestViewHandler(t *testing.T) {
	page := &Page{
		Title: "wiki",
		Body:  []byte("wiki page!"),
	}
	page.save()
	req, err := http.NewRequest("GET", "http://localhost/view/wiki", nil)
	genericErrorHandler(t, err)
	w := httptest.NewRecorder()
	viewHandler(w, req)
	resp := w.Result()
	body, err2 := ioutil.ReadAll(resp.Body)
	genericErrorHandler(t, err2)
	expectedHtmlBody := fmt.Sprintf("<h1>%s</h1><div>%s</div>", page.Title, string(page.Body))
	contentErrorHandler(t, string(body), expectedHtmlBody)
	if resp.StatusCode != 200 {
		t.Errorf("received %q status when I expected %q", resp.StatusCode, 200)
	}
}

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/", nil)
	genericErrorHandler(t, err)
	w := httptest.NewRecorder()
	rootHandler(w, req)
	resp := w.Result()
	body, reqErr := ioutil.ReadAll(resp.Body)
	genericErrorHandler(t, reqErr)
	expectedBody := "Nothing to see here!"
	contentErrorHandler(t, expectedBody, string(body))
}

func TestEditHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/edit/newpage", nil)
	genericErrorHandler(t, err)
	w := httptest.NewRecorder()
	editHandler(w, req)
	resp := w.Result()
	body, reqErr := ioutil.ReadAll(resp.Body)
	genericErrorHandler(t, reqErr)
	expectedBody :=
		`<h1>Editing newpage</h1>

<form action="/save/newpage" method="POST">
    <div><textarea name="body" rows="20" cols="80"></textarea></div>
    <div><input type="submit" value="Save"></div>
</form>`
	contentErrorHandler(t, string(body), expectedBody)
}
