package notebooks

import (
	"fmt"
	"strconv"
	"sync"
)

// NewInMemoryNotebookStore initializes an empty notebook store.
func NewInMemoryNotebookStore() *InMemoryNotebookStore {
	return &InMemoryNotebookStore{
		map[ID]Note{},
		sync.RWMutex{},
	}
}

// InMemory PlayerStore collects data about notebooks in memory.
type InMemoryNotebookStore struct {
	Notes map[ID]Note
	// A mutex is used to synchronize read/write access to the map
	lock sync.RWMutex
}

func (i *InMemoryNotebookStore) GetAllNotes() []Note {
	i.lock.RLock()	
	defer i.lock.RUnlock()
	notes := make([]Note, len(i.Notes))

	for id, note := range i.Notes {
		if (id != "" && note != "") {
			noteStr := fmt.Sprintf("%s: %s", string(id), string(note))
			notes = append(notes, Note(noteStr))
		}
	}
	return notes
}

// SaveNote will record a new notebook
func (i *InMemoryNotebookStore) SaveNote (note Note) {
	i.lock.Lock()
	defer i.lock.Unlock()
	// BUG: delete operation will cause an ID conflict
	id := strconv.Itoa(len(i.Notes)+1)
	i.Notes[ID(id)] = note 
}

func (i *InMemoryNotebookStore) GetNoteById(id ID) Note {
	i.lock.RLock()	
	defer i.lock.RUnlock()
	return i.Notes[id]
}
