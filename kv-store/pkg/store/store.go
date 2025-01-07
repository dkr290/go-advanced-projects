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
	Delete(key string, database string)
	Save(filename string) error
	Load(filename string) error
	LoadAll(filename string) (map[string]map[string]interface{}, error)
}

type KeyValuesStore struct {
	databases map[string]*sync.Map
}

func init() {
	_ = NewKeyValuesStore()
}

func NewKeyValuesStore() *KeyValuesStore {
	return &KeyValuesStore{
		databases: make(map[string]*sync.Map),
	}
}

func (s *KeyValuesStore) Set(key string, value string, req models.JsonReqiest) {
	if _, ok := s.databases[req.Database]; !ok {
		s.databases[req.Database] = &sync.Map{}
	}
	s.databases[req.Database].Store(key, value)
}

func (s *KeyValuesStore) Get(key string, database string) (string, bool) {
	db, ok := s.databases[database]
	if !ok {
		return "", false
	}
	value, ok := db.Load(key)

	return value.(string), ok
}

func (s *KeyValuesStore) Delete(key string, database string) {
	if db, ok := s.databases[database]; ok {
		db.Delete(key)
	}
}

func (s *KeyValuesStore) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	var regularMap map[string]map[string]interface{}
	if err := decoder.Decode(&regularMap); err != nil {
		return err
	}
	s.databases = make(map[string]*sync.Map)
	for dbName, innerMap := range regularMap {
		syncMap := &sync.Map{}
		for key, value := range innerMap {
			syncMap.Store(key, value)
		}
		s.databases[dbName] = syncMap
	}

	return nil
}

func (s *KeyValuesStore) LoadAll(filename string) (map[string]map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var regularMap map[string]map[string]interface{}
	if err := decoder.Decode(&regularMap); err != nil {
		return nil, err
	}

	return regularMap, nil
}

func (s *KeyValuesStore) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// Convert map[string]*sync.Map to a regular map
	regularMap := make(map[string]map[string]interface{})
	for dbName, syncMap := range s.databases {
		innerMap := make(map[string]interface{})
		syncMap.Range(func(key, value interface{}) bool {
			innerMap[key.(string)] = value
			return true
		})
		regularMap[dbName] = innerMap
	}

	encoder := gob.NewEncoder(file)
	return encoder.Encode(regularMap)
}
