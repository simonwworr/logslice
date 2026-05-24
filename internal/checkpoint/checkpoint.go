package checkpoint

import (
	"encoding/json"
	"os"
	"time"
)

// State holds the persisted progress of a log slicing run.
type State struct {
	// File is the path of the log file being processed.
	File string `json:"file"`
	// Offset is the byte offset of the last successfully processed line.
	Offset int64 `json:"offset"`
	// LastTimestamp is the timestamp of the last processed log line.
	LastTimestamp time.Time `json:"last_timestamp"`
	// LinesRead is the total number of lines read so far.
	LinesRead int64 `json:"lines_read"`
	// SavedAt is when this checkpoint was written.
	SavedAt time.Time `json:"saved_at"`
}

// Save writes the checkpoint state to the given path as JSON.
func Save(path string, s State) error {
	s.SavedAt = time.Now()
	f, err := os.CreateTemp("", "checkpoint-*.tmp")
	if err != nil {
		return err
	}
	tmpName := f.Name()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		f.Close()
		os.Remove(tmpName)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}

// Load reads a checkpoint state from the given path.
// Returns a zero State and no error if the file does not exist.
func Load(path string) (State, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return State{}, nil
	}
	if err != nil {
		return State{}, err
	}
	defer f.Close()
	var s State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return State{}, err
	}
	return s, nil
}

// Remove deletes the checkpoint file. It is a no-op if the file does not exist.
func Remove(path string) error {
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
