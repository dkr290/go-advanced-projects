package store

import (
	"encoding/gob"
	"os"
	"sync"

	"github.com/dkr290/go-advanced-projects/kv-store/pkg/models"
)

type Store interface {
	Set(key string, value string, req models.JsonReqiest)
	Get(key string, database string) (string, bool)
	Delete(key string)
	Save(filename string) error
	Load(filename string) error
}

type KeyValuesStore struct {
	Databases map[string]*sync.Map
	mutex     sync.Mutex
}

func NewKeyValuesStore() *KeyValuesStore {
	return &KeyValuesStore{}
}

func (s *KeyValuesStore) Set(key string, value string, req models.JsonReqiest) {
	if _, ok := s.Databases[req.Database]; !ok {
		s.Databases[req.Database] = &sync.Map{}
	}
	s.Databases[req.Database].Store(key, value)
}

func (s *KeyValuesStore) Get(key string, database string) (string, bool) {
	if db, ok := s.Databases[database]; ok {
		if value, ok := db.Load(key); ok {
			return value.(string), ok
		}
	}
	return "", false
}

func (s *KeyValuesStore) Delete(key string) {
	s.data.Delete(key)
}

func (s *KeyValuesStore) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	s.data.Range(func(key, value any) bool {
		if err := encoder.Encode(key); err != nil {
			return false
		}
		if err := encoder.Encode(value); err != nil {
			return false
		}
		return true
	})
	return nil
}

func (s *KeyValuesStore) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	for {
		var key string
		var value string
		if err := decoder.Decode(&key); err != nil {
			break
		}
		if err := decoder.Decode(&value); err != nil {
			break
		}
		s.Set(key, value)
	}

	return nil
}
