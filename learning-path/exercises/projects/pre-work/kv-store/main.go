package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Store represents an in-memory key-value store
type Store struct {
	data     map[string]string
	mu       sync.RWMutex
	filename string
}

// Snapshot represents a point-in-time snapshot of the store
type Snapshot struct {
	Version   int               `json:"version"`
	Timestamp string            `json:"timestamp"`
	Data      map[string]string `json:"data"`
}

func main() {
	filename := flag.String("file", "data.json", "Persistence file path")
	autosave := flag.Duration("autosave", 0, "Auto-save interval (e.g., 30s, 1m)")
	flag.Parse()

	store := NewStore(*filename)

	// Load existing data if file exists
	if _, err := os.Stat(*filename); err == nil {
		if err := store.Load(); err != nil {
			fmt.Printf("Warning: could not load data: %v\n", err)
		} else {
			fmt.Printf("Loaded data from %s\n", *filename)
		}
	}

	// Start auto-save if enabled
	if *autosave > 0 {
		go autoSave(store, *autosave)
	}

	// Start REPL
	runREPL(store)
}

// NewStore creates a new key-value store
func NewStore(filename string) *Store {
	return &Store{
		data:     make(map[string]string),
		filename: filename,
	}
}

// Get retrieves a value by key
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// Set stores a key-value pair
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Delete removes a key
func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, existed := s.data[key]
	delete(s.data, key)
	return existed
}

// Exists checks if a key exists
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[key]
	return ok
}

// Keys returns all keys matching the pattern
func (s *Store) Keys(pattern string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string
	for k := range s.data {
		matched, err := filepath.Match(pattern, k)
		if err == nil && matched {
			keys = append(keys, k)
		}
	}
	return keys
}

// Size returns the number of keys in the store
func (s *Store) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// Clear removes all keys
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]string)
}

// Snapshot saves the store to disk
func (s *Store) Snapshot() error {
	s.mu.RLock()
	// Create a copy to avoid holding lock during I/O
	dataCopy := make(map[string]string, len(s.data))
	for k, v := range s.data {
		dataCopy[k] = v
	}
	s.mu.RUnlock()

	snapshot := Snapshot{
		Version:   1,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      dataCopy,
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling snapshot: %w", err)
	}

	// Write to temporary file first, then rename (atomic)
	tmpFile := s.filename + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("writing temp file: %w", err)
	}

	if err := os.Rename(tmpFile, s.filename); err != nil {
		return fmt.Errorf("renaming temp file: %w", err)
	}

	return nil
}

// Load restores the store from disk
func (s *Store) Load() error {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("unmarshaling snapshot: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = snapshot.Data

	return nil
}

// runREPL runs the Read-Eval-Print Loop
func runREPL(store *Store) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("KV Store - Type 'HELP' for commands")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "GET":
			if len(parts) < 2 {
				fmt.Println("Usage: GET <key>")
				continue
			}
			if val, ok := store.Get(parts[1]); ok {
				fmt.Println(val)
			} else {
				fmt.Println("(nil)")
			}

		case "SET":
			if len(parts) < 3 {
				fmt.Println("Usage: SET <key> <value>")
				continue
			}
			value := strings.Join(parts[2:], " ")
			store.Set(parts[1], value)
			fmt.Println("OK")

		case "DELETE", "DEL":
			if len(parts) < 2 {
				fmt.Println("Usage: DELETE <key>")
				continue
			}
			if store.Delete(parts[1]) {
				fmt.Println("OK")
			} else {
				fmt.Println("Key not found")
			}

		case "EXISTS":
			if len(parts) < 2 {
				fmt.Println("Usage: EXISTS <key>")
				continue
			}
			if store.Exists(parts[1]) {
				fmt.Println("1")
			} else {
				fmt.Println("0")
			}

		case "KEYS":
			pattern := "*"
			if len(parts) >= 2 {
				pattern = parts[1]
			}
			keys := store.Keys(pattern)
			for _, key := range keys {
				fmt.Println(key)
			}

		case "SIZE":
			fmt.Println(store.Size())

		case "CLEAR":
			store.Clear()
			fmt.Println("OK")

		case "SNAPSHOT", "SAVE":
			if err := store.Snapshot(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Saved snapshot to %s\n", store.filename)
			}

		case "HELP":
			printHelp()

		case "EXIT", "QUIT":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Printf("Unknown command: %s\n", command)
			fmt.Println("Type 'HELP' for available commands")
		}
	}
}

func printHelp() {
	help := `
Available Commands:
  GET <key>           Get value for key
  SET <key> <value>   Set key to value
  DELETE <key>        Delete key
  EXISTS <key>        Check if key exists (returns 1 or 0)
  KEYS [pattern]      List keys matching pattern (default: *)
  SIZE                Get number of keys
  CLEAR               Remove all keys
  SNAPSHOT            Save to disk
  HELP                Show this help
  EXIT                Exit the program
`
	fmt.Println(help)
}

func autoSave(store *Store, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := store.Snapshot(); err != nil {
			fmt.Printf("\nAuto-save error: %v\n> ", err)
		} else {
			fmt.Printf("\n[Auto-saved to %s]\n> ", store.filename)
		}
	}
}
