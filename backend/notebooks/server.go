package notebooks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// NoteStore stores notes
type Note string
type ID string

type NoteStore interface {
	GetAllNotes() []Note
	GetNoteById(id ID) Note
	SaveNote(note Note)
}

type NotebookServer struct {
	Store NoteStore
}

/**************** Routing ******************/

func (n *NotebookServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()
	router.Handle("/notes", http.HandlerFunc(n.notesHandler))
	router.Handle("/notes/", http.HandlerFunc(n.notesHandler))

	router.ServeHTTP(w, r)
}

/**************** Handlers ******************/

func (n *NotebookServer) notesHandler (w http.ResponseWriter, r *http.Request) {
	noteId := strings.TrimPrefix(r.URL.Path, "/notes/")

	switch r.Method {
	case http.MethodGet:
		// Get all notes
		if r.URL.Path == "/notes" {
			n.showAllNotes(w)
			return
		} 

		// Get note by ID
		if strings.TrimSpace(noteId) != "" {
			n.showNoteById(w, ID(noteId))
			return 
		}
	case http.MethodPost:
		bodyBytes := readBody(w, r)
		n.addNote(w, Note(bodyBytes))
	}
}

/**************** Controllers ******************/

func (n *NotebookServer) showAllNotes (w http.ResponseWriter) {

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	notes := n.Store.GetAllNotes()
	_notes := make([]string, 0)

	for _, note := range notes {
		if strings.TrimSpace(string(note)) != "" {
			_notes = append(_notes, strings.TrimSpace(string(note)))
		}
	}


	notesString := strings.Join(_notes, ", ")
	fmt.Fprint(w, notesString)
}

func (n *NotebookServer) showNoteById (w http.ResponseWriter, noteId ID) {
	/** Find note by ID */
	note := n.Store.GetNoteById(noteId)

	if note == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Note not found!")
		return 
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(note))
}

func (n *NotebookServer) addNote (w http.ResponseWriter, note Note) {
	n.Store.SaveNote(note)
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "Note successfully added!")
}

/**************** Utils ******************/

func readBody (w http.ResponseWriter ,r *http.Request) []byte {
	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)	
		return nil
	}

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		errorMessage := "Unable to read body"
		// Client error
		http.Error(w, errorMessage, http.StatusInternalServerError)
		// Server log error
		log.Printf("%s: %s", errorMessage, err.Error())

		return nil 
	}
	defer r.Body.Close()
	return bodyBytes
}
