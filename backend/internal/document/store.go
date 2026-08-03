package document

import "sync"

type Store interface {
	Save(doc *Document) error
	Get(id string) (*Document, error)
	List() ([]*Document, error)
}

type InMemoryStore struct {
	docs map[string]*Document
	mu   sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{docs: make(map[string]*Document)}
}

func (s *InMemoryStore) Save(doc *Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.docs[doc.ID] = doc
	return nil
}

func (s *InMemoryStore) Get(id string) (*Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.docs[id]
	if !ok {
		return nil, nil
	}
	return doc, nil
}

func (s *InMemoryStore) List() ([]*Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*Document, 0, len(s.docs))
	for _, doc := range s.docs {
		list = append(list, doc)
	}
	return list, nil
}
