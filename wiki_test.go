package main

import (
	"io/ioutil"
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
