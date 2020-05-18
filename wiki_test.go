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
