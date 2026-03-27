package hotreload

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ============================================================
// FILE WATCHER
// ============================================================

// FileWatcher monitors files for changes
type FileWatcher struct {
	mu        sync.RWMutex
	files     map[string]*watchedFile
	callbacks map[string][]WatchCallback
	stopCh   chan struct{}
	stopped  bool
	interval time.Duration
}

// watchedFile tracks a single file's state
type watchedFile struct {
	path     string
	modTime  time.Time
	size     int64
	checksum uint32
}

// WatchCallback is called when a file changes
type WatchCallback func(event FileEvent) error

// FileEvent describes what changed
type FileEvent struct {
	Path      string
	Type      EventType
	Timestamp time.Time
	OldModTime time.Time
	NewModTime time.Time
}

// EventType describes the type of file change
type EventType int

const (
	EventModified EventType = iota
	EventCreated
	EventDeleted
	EventRenamed
)

// String returns the string representation of EventType
func (et EventType) String() string {
	switch et {
	case EventModified:
		return "modified"
	case EventCreated:
		return "created"
	case EventDeleted:
		return "deleted"
	case EventRenamed:
		return "renamed"
	default:
		return "unknown"
	}
}

// NewFileWatcher creates a new file watcher
func NewFileWatcher(checkInterval time.Duration) *FileWatcher {
	return &FileWatcher{
		files:     make(map[string]*watchedFile),
		callbacks: make(map[string][]WatchCallback),
		stopCh:    make(chan struct{}),
		interval:  checkInterval,
	}
}

// Watch starts watching a file for changes
func (fw *FileWatcher) Watch(path string, callback WatchCallback) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	// Get initial file info
	info, err := getFileInfo(path)
	if err != nil {
		return fmt.Errorf("watch failed: %w", err)
	}

	// Store watched file
	fw.files[path] = info

	// Store callback
	fw.callbacks[path] = append(fw.callbacks[path], callback)

	return nil
}

// Unwatch stops watching a file
func (fw *FileWatcher) Unwatch(path string) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	delete(fw.files, path)
	delete(fw.callbacks, path)
}

// Start begins watching all registered files
func (fw *FileWatcher) Start() {
	go fw.watchLoop()
}

// Stop stops watching all files
func (fw *FileWatcher) Stop() {
	fw.mu.Lock()
	fw.stopped = true
	fw.mu.Unlock()
	close(fw.stopCh)
}

// watchLoop runs the main watch loop
func (fw *FileWatcher) watchLoop() {
	ticker := time.NewTicker(fw.interval)
	defer ticker.Stop()

	for {
		select {
		case <-fw.stopCh:
			return
		case <-ticker.C:
			fw.checkFiles()
		}
	}
}

// checkFiles checks all watched files for changes
func (fw *FileWatcher) checkFiles() {
	fw.mu.RLock()
	files := make(map[string]*watchedFile)
	callbacks := make(map[string][]WatchCallback)
	
	for path, file := range fw.files {
		files[path] = file
		callbacks[path] = fw.callbacks[path]
	}
	fw.mu.RUnlock()

	for path, oldInfo := range files {
		newInfo, err := getFileInfo(path)
		if err != nil {
			// File deleted or inaccessible
			fw.notifyCallbacks(path, callbacks[path], FileEvent{
				Path:       path,
				Type:       EventDeleted,
				Timestamp:  time.Now(),
				OldModTime: oldInfo.modTime,
			})
			fw.Unwatch(path)
			continue
		}

		// Check if file changed
		if newInfo.modTime != oldInfo.modTime || newInfo.size != oldInfo.size {
			fw.notifyCallbacks(path, callbacks[path], FileEvent{
				Path:       path,
				Type:       EventModified,
				Timestamp:  time.Now(),
				OldModTime: oldInfo.modTime,
				NewModTime: newInfo.modTime,
			})

			// Update tracked info
			fw.mu.Lock()
			fw.files[path] = newInfo
			fw.mu.Unlock()
		}
	}
}

// notifyCallbacks calls all callbacks for a file event
func (fw *FileWatcher) notifyCallbacks(path string, cbs []WatchCallback, event FileEvent) {
	for _, cb := range cbs {
		go func(callback WatchCallback) {
			if err := callback(event); err != nil {
				fmt.Printf("[HOTRELOAD] Callback error for %s: %v\n", path, err)
			}
		}(cb)
	}
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

// getFileInfo gets file information
func getFileInfo(path string) (*watchedFile, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &watchedFile{
		path:    path,
		modTime: stat.ModTime(),
		size:    stat.Size(),
	}, nil
}
