package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ProgressStore struct {
	Books map[string]int `json:"books"`
}

func progressFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "glance", "progress.json"), nil
}

func LoadProgress() (*ProgressStore, error) {
	path, err := progressFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &ProgressStore{Books: make(map[string]int)}, nil
	}
	if err != nil {
		return nil, err
	}

	store := &ProgressStore{Books: make(map[string]int)}
	if err := json.Unmarshal(data, store); err != nil {
		return nil, err
	}

	if store.Books == nil {
		store.Books = make(map[string]int)
	}

	return store, nil
}

func SaveProgress(store *ProgressStore) error {
	path, err := progressFilePath()
	if err != nil {
		return err
	}

	if store == nil {
		store = &ProgressStore{Books: make(map[string]int)}
	}
	if store.Books == nil {
		store.Books = make(map[string]int)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
