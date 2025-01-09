package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dkr290/go-advanced-projects/kv-store/pkg/models"
)

type Store interface {
	Set(key string, value string, req models.JsonRequest) error
	Get(key string, database string) (string, bool)
	Delete(key string, database string)
	Load(filename string) error
	LoadAll(filename string) (map[string]any, error)
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

func (s *KeyValuesStore) Set(key string, value string, req models.JsonRequest) error {
	// Check if the database exists in memory
	db, ok := s.databases[req.Database]
	if !ok {
		db = &sync.Map{}
		s.databases[req.Database] = db
	}

	// Check if the key exists
	if _, ok := db.Load(key); ok {
		return fmt.Errorf("the same key already exists in the database %s", key)
	}
	// Store the key-value pair in memory
	db.Store(key, value)
	// Append the new key-value pair to the file
	file, err := os.OpenFile(req.Database+".jsonl", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open database file for appending: %v", err)
	}
	defer file.Close()
	// Write the key-value pair as a JSON object
	entry := map[string]string{
		"key":   key,
		"value": value,
	}
	entryJSON, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %v", err)
	}
	if _, err := file.WriteString(string(entryJSON) + "\n"); err != nil {
		return fmt.Errorf("failed to write entry to file: %v", err)
	}

	return nil
}

func (s *KeyValuesStore) Get(key string, database string) (string, bool) {
	_ = s.Load(database + ".jsonl")

	db := s.databases[database]
	if d, ok := db.Load(key); ok {
		return d.(string), ok
	}
	return "", false
}

func (s *KeyValuesStore) Delete(key string, database string) {
	// if db, ok := s.databases[database]; ok {
	// 	db.Delete(key)
	// }
	// todo
}

func (s *KeyValuesStore) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	// Create a new sync.Map for this file
	db := &sync.Map{}

	scanner := bufio.NewScanner(file)
	// read the file
	for scanner.Scan() {
		line := scanner.Text()
		entry := make(map[string]string)
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return fmt.Errorf("failed to unmarshal line: %v", err)
		}
		// Extract the key and value
		key, keyExists := entry["key"]
		value, valueExists := entry["value"]
		if !keyExists || !valueExists {
			return fmt.Errorf("missing key or value in entry: %s", line)
		}

		// Store the key-value pair in the sync.Map
		db.Store(key, value)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading file: %v", err)
	}
	// Store the sync.Map in the in-memory database
	s.databases[strings.Split(filename, ".")[0]] = db
	return nil
}

func (s *KeyValuesStore) LoadAll(filename string) (map[string]any, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	result := make(map[string]interface{})
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// Parse each line as a JSON object
		line := scanner.Text()
		entry := make(map[string]string)
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal line: %v", err)
		}

		// Extract key and value and store in the result map
		key, keyExists := entry["key"]
		value, valueExists := entry["value"]
		if !keyExists || !valueExists {
			return nil, fmt.Errorf("missing key or value in entry: %s", line)
		}
		result[key] = value
	}

	// Check for errors during file scanning
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading file %s: %v", filename, err)
	}

	return result, nil
}
