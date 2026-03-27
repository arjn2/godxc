package hotreload

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"dxgoc/internal/dxc"
)

// ============================================================
// SHADER HOT RELOAD MANAGER
// ============================================================

// ShaderReloadCallback is called when a shader is reloaded
type ShaderReloadCallback func(event ShaderReloadEvent) error

// ShaderReloadEvent describes a shader reload event
type ShaderReloadEvent struct {
	Path       string      // Full path to shader file
	Name       string      // Shader name
	Success    bool        // Whether compilation succeeded
	Output     []byte      // Compiled output
	Error      error       // Compilation error if any
	Timestamp  time.Time   // When the reload occurred
	Duration   time.Duration // How long compilation took
}

// ShaderConfig holds compilation configuration for a shader
type ShaderConfig struct {
	Path        string   // Path to .hlsl file
	Name        string   // Identifier name
	Profile     string   // Shader profile (vs_6_0, ps_6_0, etc.)
	Entry       string   // Entry point function
	Defines     []string // Preprocessor defines
	Includes    []string // Include paths
	OutputPath  string   // Where to write compiled output
	AutoCompile bool     // Automatically compile on changes
}

// ShaderReloader manages hot reloading of shaders
type ShaderReloader struct {
	mu          sync.RWMutex
	watcher     *FileWatcher
	compiler    *dxc.Compiler
	shaders     map[string]*ShaderConfig
	callbacks   map[string][]ShaderReloadCallback
	isRunning   bool
	stopCh      chan struct{}
	debounce    map[string]*time.Timer
	debounceMs  time.Duration
}

// NewShaderReloader creates a new shader hot reload manager
func NewShaderReloader(dllPath string, debounceMs time.Duration) (*ShaderReloader, error) {
	// Initialize compiler
	compiler, err := dxc.NewCompiler(dllPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create compiler: %w", err)
	}

	return &ShaderReloader{
		watcher:    NewFileWatcher(100 * time.Millisecond),
		compiler:   compiler,
		shaders:    make(map[string]*ShaderConfig),
		callbacks:  make(map[string][]ShaderReloadCallback),
		stopCh:     make(chan struct{}),
		debounce:   make(map[string]*time.Timer),
		debounceMs: debounceMs,
	}, nil
}

// RegisterShader registers a shader for hot reloading
func (sr *ShaderReloader) RegisterShader(config ShaderConfig) error {
	if config.Name == "" {
		return fmt.Errorf("shader name is required")
	}
	if config.Path == "" {
		return fmt.Errorf("shader path is required")
	}
	if config.Profile == "" {
		return fmt.Errorf("shader profile is required")
	}
	if config.Entry == "" {
		config.Entry = "main" // Default entry point
	}

	sr.mu.Lock()
	defer sr.mu.Unlock()

	sr.shaders[config.Name] = &config

	// Start watching if already running
	if sr.isRunning && config.AutoCompile {
		err := sr.watcher.Watch(config.Path, sr.makeWatchCallback(config.Name))
		if err != nil {
			return fmt.Errorf("failed to watch shader: %w", err)
		}
	}

	return nil
}

// UnregisterShader stops watching a shader
func (sr *ShaderReloader) UnregisterShader(name string) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	config, ok := sr.shaders[name]
	if !ok {
		return
	}

	sr.watcher.Unwatch(config.Path)
	delete(sr.shaders, name)
	delete(sr.callbacks, name)

	// Cancel any pending recompilation
	if timer, ok := sr.debounce[name]; ok {
		timer.Stop()
		delete(sr.debounce, name)
	}
}

// OnReload registers a callback for shader reload events
func (sr *ShaderReloader) OnReload(name string, callback ShaderReloadCallback) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if _, ok := sr.shaders[name]; !ok {
		return fmt.Errorf("shader %s not registered", name)
	}

	sr.callbacks[name] = append(sr.callbacks[name], callback)
	return nil
}

// Start starts the hot reload manager
func (sr *ShaderReloader) Start() error {
	sr.mu.Lock()
	if sr.isRunning {
		sr.mu.Unlock()
		return fmt.Errorf("shader reloader already running")
	}

	sr.isRunning = true
	shaders := make(map[string]*ShaderConfig)
	for name, config := range sr.shaders {
		shaders[name] = config
	}
	sr.mu.Unlock()

	// Start watcher
	sr.watcher.Start()

	// Watch all auto-compile shaders
	for name, config := range shaders {
		if config.AutoCompile {
			err := sr.watcher.Watch(config.Path, sr.makeWatchCallback(name))
			if err != nil {
				return fmt.Errorf("failed to watch shader %s: %w", name, err)
			}
		}
	}

	return nil
}

// Stop stops the hot reload manager
func (sr *ShaderReloader) Stop() {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if !sr.isRunning {
		return
	}

	sr.isRunning = false
	sr.watcher.Stop()

	// Cancel pending recompilations
	for _, timer := range sr.debounce {
		timer.Stop()
	}
	sr.debounce = make(map[string]*time.Timer)

	close(sr.stopCh)
}

// Recompile manually recompiles a shader
func (sr *ShaderReloader) Recompile(name string) error {
	sr.mu.RLock()
	config, ok := sr.shaders[name]
	if !ok {
		sr.mu.RUnlock()
		return fmt.Errorf("shader %s not registered", name)
	}
	sr.mu.RUnlock()

	return sr.compileShader(name, config)
}

// compileShader compiles a shader and notifies callbacks
func (sr *ShaderReloader) compileShader(name string, config *ShaderConfig) error {
	startTime := time.Now()

	// Read shader file
	source, err := os.ReadFile(config.Path)
	if err != nil {
		sr.notifyCallbacks(name, ShaderReloadEvent{
			Path:      config.Path,
			Name:      name,
			Success:   false,
			Error:     err,
			Timestamp: startTime,
			Duration:  time.Since(startTime),
		})
		return err
	}

	// Build arguments
	args := []string{
		"-T", config.Profile,
		"-E", config.Entry,
	}

	// Add defines
	for _, def := range config.Defines {
		args = append(args, "-D", def)
	}

	// Add includes
	for _, inc := range config.Includes {
		args = append(args, "-I", inc)
	}

	// Compile
	result, err := sr.compiler.Compile(source, args)
	duration := time.Since(startTime)

	if err != nil {
		sr.notifyCallbacks(name, ShaderReloadEvent{
			Path:      config.Path,
			Name:      name,
			Success:   false,
			Error:     err,
			Timestamp: startTime,
			Duration:  duration,
		})
		return err
	}

	// Check compilation status
	success := result.Success
	if !success {
		err = fmt.Errorf("compilation failed: %s", string(result.Errors))
	}

	event := ShaderReloadEvent{
		Path:      config.Path,
		Name:      name,
		Success:   success,
		Output:    result.Output,
		Error:     err,
		Timestamp: startTime,
		Duration:  duration,
	}

	// Write output if successful
	if success && config.OutputPath != "" {
		// Create output directory if needed
		outDir := filepath.Dir(config.OutputPath)
		os.MkdirAll(outDir, 0755)

		if writeErr := os.WriteFile(config.OutputPath, result.Output, 0644); writeErr != nil {
			event.Error = fmt.Errorf("failed to write output: %w", writeErr)
			event.Success = false
		}
	}

	sr.notifyCallbacks(name, event)
	return nil
}

// makeWatchCallback creates a watch callback for a shader
func (sr *ShaderReloader) makeWatchCallback(name string) WatchCallback {
	return func(event FileEvent) error {
		if event.Type != EventModified {
			return nil
		}

		sr.mu.RLock()
		config, ok := sr.shaders[name]
		debounceMs := sr.debounceMs
		sr.mu.RUnlock()

		if !ok {
			return nil
		}

		// Debounce: cancel existing timer and create new one
		sr.mu.Lock()
		if existingTimer, ok := sr.debounce[name]; ok {
			existingTimer.Stop()
		}

		timer := time.AfterFunc(debounceMs, func() {
			sr.mu.Lock()
			delete(sr.debounce, name)
			sr.mu.Unlock()

			_ = sr.compileShader(name, config)
		})

		sr.debounce[name] = timer
		sr.mu.Unlock()

		return nil
	}
}

// notifyCallbacks calls all registered callbacks for a shader
func (sr *ShaderReloader) notifyCallbacks(name string, event ShaderReloadEvent) {
	sr.mu.RLock()
	callbacks := sr.callbacks[name]
	sr.mu.RUnlock()

	for _, cb := range callbacks {
		go func(callback ShaderReloadCallback) {
			if err := callback(event); err != nil {
				fmt.Printf("[HOTRELOAD] Callback error for %s: %v\n", name, err)
			}
		}(cb)
	}
}

// Close closes the shader reloader and releases resources
func (sr *ShaderReloader) Close() {
	sr.Stop()
	if sr.compiler != nil {
		sr.compiler.Close()
	}
}
