package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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
	viewHandler(w, req, "/view/wiki")
	resp := w.Result()
	body, err2 := ioutil.ReadAll(resp.Body)
	genericErrorHandler(t, err2)
	expectedHtmlBody := `<h1>wiki</h1>

<p>[<a href="/edit/wiki">edit</a>]</p>

<div>wiki page!</div>`
	contentErrorHandler(t, string(body), expectedHtmlBody)
	if resp.StatusCode != 200 {
		t.Errorf("received %q status when I expected %q", resp.StatusCode, 200)
	}
}

func TestViewHandlerRedirect(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost/view/unknown", nil)
	w := httptest.NewRecorder()
	viewHandler(w, req, "/view/unknown")
	resp := w.Result()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("received status code %v when I expected a %v", resp.StatusCode, http.StatusFound)
	}
	if url, _ := resp.Location(); url.Path != "/edit/unknown" {
		t.Errorf("unexpected path -- got %v when I wanted /edit/unknown", url.Path)
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
	editHandler(w, req, "/edit/newpage")
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

func TestSaveHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost/save/saveme", strings.NewReader("save me!"))
	w := httptest.NewRecorder()
	saveHandler(w, req, "/save/saveme")
	resp := w.Result()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("received status code %v when I expected a %v", resp.StatusCode, http.StatusFound)
	}
	if url, _ := resp.Location(); url.Path != "/view/saveme" {
		t.Errorf("unexpected path -- got %v when I wanted /view/saveme", url.Path)
	}
}

func TestPathTraversalInURLFails(t *testing.T) {
	var (
		editPath = "http://localhost/edit/../../../etc/passwd"
		savePath = "http://localhost/save/newpage/../../../../bin/rootKit"
		viewPath = "http://localhost/view/../../../etc/shadow"
	)
	req := httptest.NewRequest("GET", editPath, nil)
	w := httptest.NewRecorder()
	eh := makeHandler(editHandler)
	eh(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("did not get a 404 response, got %v instead", resp.StatusCode)
	}
	saveReq := httptest.NewRequest("GET", savePath, nil)
	saveW := httptest.NewRecorder()
	handler := makeHandler(saveHandler)
	handler(saveW, saveReq)
	saveResp := saveW.Result()
	if saveResp.StatusCode != http.StatusNotFound {
		t.Errorf("did not get a 404 response, got %v instead", saveResp.StatusCode)
	}
	viewReq := httptest.NewRequest("GET", viewPath, nil)
	viewW := httptest.NewRecorder()
	viewHandle := makeHandler(viewHandler)
	viewHandle(viewW, viewReq)
	viewResp := viewW.Result()
	if viewResp.StatusCode != http.StatusNotFound {
		t.Errorf("did not get a 404 response, got %v instead", viewResp.StatusCode)
	}
}
