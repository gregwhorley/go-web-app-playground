package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWikiSave(t *testing.T) {
	page := &Page{
		Title: "PageTitle",
		Body:  []byte("this is the page body"),
	}
	err := page.save()
	if err != nil {
		t.Errorf("test failed -- %q", err)
	}
	content, err2 := ioutil.ReadFile(page.Title + ".txt")
	if err2 != nil {
		t.Errorf("test failed -- %q", err2)
	}
	stringContent := string(content)
	pageContent := string(page.Body)
	if stringContent != pageContent {
		t.Errorf("test failed -- %q", err2)
	}
}

func TestWikiLoadPage(t *testing.T) {
	page := &Page{
		Title: "LoadPage",
		Body:  []byte("testing load page"),
	}
	page.save()
	page2, err := loadPage(page.Title)
	if err != nil {
		t.Errorf("test failed -- %q", err)
	}
	expectedContent := string(page.Body)
	loadPageContent := string(page2.Body)
	if loadPageContent != expectedContent {
		t.Errorf("test failed -- content does not match\n %q \n %q", expectedContent, loadPageContent)
	}
}

func TestViewHandler(t *testing.T) {
	page := &Page{
		Title: "wiki",
		Body:  []byte("wiki page!"),
	}
	page.save()
	req, err := http.NewRequest("GET", "http://localhost/view/wiki", nil)
	if err != nil {
		t.Errorf("test failed -- %q", err)
	}
	w := httptest.NewRecorder()
	viewHandler(w, req)
	resp := w.Result()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		t.Errorf("test failed -- %q", err2)
	}
	expectedHtmlBody := fmt.Sprintf("<h1>%v</h1><div>%v</div>", page.Title, string(page.Body))
	if string(body) != expectedHtmlBody {
		t.Errorf("content does not match \n %v \n %v", string(body), expectedHtmlBody)
	}
	if resp.StatusCode != 200 {
		t.Errorf("received %q status when I expected %q", resp.StatusCode, 200)
	}
}

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		t.Errorf("test failed -- %q", err)
	}
	w := httptest.NewRecorder()
	rootHandler(w, req)
	resp := w.Result()
	body, reqErr := ioutil.ReadAll(resp.Body)
	if reqErr != nil {
		t.Errorf("test failed -- %q", err)
	}
	expectedBody := "Nothing to see here!"
	if string(body) != expectedBody {
		t.Errorf("content does not match \n %v \n %v", string(body), expectedBody)
	}
}