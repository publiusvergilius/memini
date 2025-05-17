package notebooks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


func TestGETNotebooks (t *testing.T) {
	t.Run("get from in memory notebook store", func(t *testing.T) {
		store :=  NewInMemoryNotebookStore()
		server := &NotebookServer{store}
		store.Notes["1"] = "teste 1"
		store.Notes["2"] = "teste 2"

		tests := []struct {
			testName           string
			noteId             ID
			expectedHTTPStatus int
			expectedNote       Note
		}{
			{
				testName: "Returns first note",
				noteId: "1",
				expectedHTTPStatus: http.StatusOK,
				expectedNote: "teste 1",
			},
			{
				testName: "Returns second note",
				noteId: "2",
				expectedHTTPStatus: http.StatusOK,
				expectedNote: "teste 2",
			},
			{
				testName: "Returns 404 on misssing note",
				noteId: "3",
				expectedHTTPStatus: http.StatusNotFound,
				expectedNote: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.testName, func(t *testing.T) {
				request := newGetNotebookRequest(tt.noteId)
				response := httptest.NewRecorder()

				server.ServeHTTP(response, request)

				assertStatus(t, response.Code, tt.expectedHTTPStatus)
			})
		}
	})
	

	t.Run("get all in plain text from in memory store", func(t *testing.T) {
		store := NewInMemoryNotebookStore()
		store.Notes["1"] = "teste 1"
		store.Notes["2"] = "teste 2"

		server := &NotebookServer{store}

		request, _ := http.NewRequest("GET", "/notes", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "text/plain")

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		want := "1: teste 1, 2: teste 2"

		assertResponseBody(t, response.Body.String(), want)
	})
}

func TestPOSTNotebook (t *testing.T) {
	
	store := NewInMemoryNotebookStore()
	server := &NotebookServer{store}

	t.Run("it returns accepted on POST a save one", func(t *testing.T) {
		var note Note = "teste 1"
		request, _ := http.NewRequest(http.MethodPost, "/notes", strings.NewReader(string(note)))
		request.Header.Set("Content-Type", "text/plain")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.Notes) != 1 {
			t.Fatalf("got %d notes, want %d", len(store.Notes), 1)
		}

		if store.Notes["1"] != note {
			t.Errorf("did not store correnct note: got %q want %q", store.Notes["1"], note)
		}
	})

	t.Run("it returns accepted on POST a save a second", func(t *testing.T) {
		var note Note = "teste 2"
		request, _ := http.NewRequest(http.MethodPost, "/notes", strings.NewReader(string(note)))
		request.Header.Set("Content-Type", "text/plain")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.Notes) != 2 {
			t.Fatalf("got %d notes, want %d", len(store.Notes), 2)
		}

		if store.Notes["2"] != note {
			t.Errorf("did not store correct note: got %q want %q", store.Notes["2"], note)
		}
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newGetNotebookRequest(id ID) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/notes/%s", id), nil)
	return req
}

func newPostNotebookRequest(_ Note) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/notes", nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string){
	t.Helper()
	if got != want {
		t.Fatalf("response body is wrong, got %q want %q", got, want)
	}
}

func assertNotError (t testing.TB, err error) {
	t.Helper()
	t.Fatalf("was not told to error: %s", err.Error())
}

