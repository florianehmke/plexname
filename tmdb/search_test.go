package tmdb_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imba3r/grabber/core/tmdb"
)

func TestSearch_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("../../tests/fixtures/tmdb-search.json")
		if err != nil {
			t.Error(err)
		} else {
			w.Write(f)
		}
	}))
	defer ts.Close()

	s := tmdb.NewService(ts.URL, "apiKey")
	r, err := s.Search("Test", -1, -1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if r.TotalResults != 14 {
		t.Errorf("Expected 14 results, got %d", r.TotalResults)
	}
	if len(r.Results) != 14 {
		t.Errorf("Expected 14 results, got %d", len(r.Results))
	}
}

func TestSearch_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := ioutil.ReadFile("../../tests/fixtures/tmdb-error.json")
		if err != nil {
			t.Error(err)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(f)
		}
	}))
	defer ts.Close()

	s := tmdb.NewService(ts.URL, "apiKey")
	_, err := s.Search("Test", -1, -1)
	if err == nil {
		t.Errorf("Expected an error")
	}
	if err.Error() != "The resource you requested could not be found." {
		t.Errorf("Expected a different error message.")
	}
}
