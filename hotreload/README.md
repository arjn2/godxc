# Shader Hot Reloading

This package provides automatic shader hot reloading functionality for dxgoc. It monitors shader files for changes and automatically recompiles them using the DXC compiler.

## Features

- **File Monitoring** - Automatically detects shader file changes
- **Debouncing** - Prevents excessive recompilation on rapid file changes
- **Multi-Shader** - Monitor and compile multiple shaders simultaneously
- **Callbacks** - Get notified when compilation completes
- **Manual Recompile** - Manually trigger recompilation anytime
- **Error Handling** - Graceful error handling and reporting
- **Configuration** - Per-shader compilation settings (profile, entry point, defines, etc.)

## Architecture

```
FileWatcher (lower level)
    ↓
ShaderReloader (high level)
    ↓
Your Application
```

### FileWatcher
- Monitors file system for changes
- Detects modifications at specified intervals
- Notifies registered callbacks

### ShaderReloader
- Manages shader compilation pipeline
- Handles debouncing and retry logic
- Coordinates with DXC compiler
- Manages callbacks and output

## Usage

### Basic Example

```go
package main

import (
	"time"
	"dxgoc/hotreload"
)

func main() {
	// Create reloader with 500ms debounce
	reloader, err := hotreload.NewShaderReloader("", 500*time.Millisecond)
	if err != nil {
		panic(err)
	}
	defer reloader.Close()

	// Register shader
	err = reloader.RegisterShader(hotreload.ShaderConfig{
		Path:        "shaders/vertex.hlsl",
		Name:        "VertexShader",
		Profile:     "vs_6_0",
		Entry:       "main",
		OutputPath:  "output/vertex.dxil",
		AutoCompile: true,
	})
	if err != nil {
		panic(err)
	}

	// Register callback
	err = reloader.OnReload("VertexShader", func(event hotreload.ShaderReloadEvent) error {
		if event.Success {
			println("✓ Recompiled:", event.Name, event.Duration)
		} else {
			println("✗ Failed:", event.Name, event.Error)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Start monitoring
	reloader.Start()
	defer reloader.Stop()

	// Keep running
	time.Sleep(1 * time.Minute)
}
```

### Multi-Shader Setup

```go
reloader, _ := hotreload.NewShaderReloader("", 300*time.Millisecond)
defer reloader.Close()

// Register multiple shaders
for _, shader := range []struct {
	path    string
	name    string
	profile string
}{
	{"shaders/vertex.hlsl", "VS", "vs_6_0"},
	{"shaders/pixel.hlsl", "PS", "ps_6_0"},
	{"shaders/compute.hlsl", "CS", "cs_6_0"},
} {
	reloader.RegisterShader(hotreload.ShaderConfig{
		Path:        shader.path,
		Name:        shader.name,
		Profile:     shader.profile,
		Entry:       "main",
		OutputPath:  fmt.Sprintf("output/%s.dxil", shader.name),
		AutoCompile: true,
	})
}

reloader.Start()
defer reloader.Stop()
```

### With Defines and Includes

```go
reloader.RegisterShader(hotreload.ShaderConfig{
	Path:       "shaders/shader.hlsl",
	Name:       "MyShader",
	Profile:    "ps_6_0",
	Entry:      "PSMain",
	Defines:    []string{"ENABLE_PBR=1", "MAX_LIGHTS=8"},
	Includes:   []string{"shaders/common", "shaders/headers"},
	OutputPath: "output/shader.dxil",
})
```

## Configuration

### ShaderConfig

```go
type ShaderConfig struct {
	Path        string   // Path to .hlsl file (required)
	Name        string   // Identifier name (required)
	Profile     string   // Shader profile: vs_6_0, ps_6_0, cs_6_0, etc. (required)
	Entry       string   // Entry point function (default: "main")
	Defines     []string // Preprocessor defines (-D flag)
	Includes    []string // Include paths (-I flag)
	OutputPath  string   // Where to write compiled output
	AutoCompile bool     // Auto-compile on file changes
}
```

### ShaderReloadEvent

```go
type ShaderReloadEvent struct {
	Path      string        // Full path to shader file
	Name      string        // Shader name
	Success   bool          // Compilation success
	Output    []byte        // Compiled bytecode
	Error     error         // Error if compilation failed
	Timestamp time.Time     // When compilation occurred
	Duration  time.Duration // How long compilation took
}
```

## API Reference

### ShaderReloader

#### NewShaderReloader(dllPath string, debounceMs time.Duration)
Creates a new shader reloader instance.
- `dllPath`: Path to dxcompiler.dll (empty = auto-detect)
- `debounceMs`: Milliseconds to wait before recompiling after file change

#### RegisterShader(config ShaderConfig)
Registers a shader for monitoring.

#### UnregisterShader(name string)
Stops monitoring a shader.

#### OnReload(name string, callback ShaderReloadCallback)
Registers a callback for reload events.

#### Start()
Starts file monitoring and hot reload.

#### Stop()
Stops hot reload (can be restarted).

#### Recompile(name string)
Manually recompile a shader immediately.

#### Close()
Closes reloader and releases resources.

## Performance Considerations

- **Debouncing**: Default 300-500ms debounce prevents rapid recompilation
- **File Watching**: Polls at 100ms intervals (configurable)
- **Threading**: Callbacks run on goroutines (non-blocking)
- **Memory**: One compiler instance shared across all shaders

## Error Handling

Compilation errors are reported in `ShaderReloadEvent.Error`. Callbacks should handle errors gracefully:

```go
reloader.OnReload("MyShader", func(event hotreload.ShaderReloadEvent) error {
	if !event.Success {
		log.Printf("Compilation error: %v", event.Error)
		// Handle error - perhaps notify UI or fallback to previous version
		return nil // Don't propagate error
	}
	
	// Process successful compilation
	return nil
})
```

## Limitations

- File watching is polling-based (not using Windows file system API)
- Requires DXC to be installed for compilation
- No support for shader dependencies/includes yet
- Single DXC compiler instance (sequential compilation)

## Future Enhancements

- [ ] Windows file system API for instant notifications
- [ ] Parallel shader compilation (multiple compiler instances)
- [ ] Shader dependency tracking
- [ ] Compiled output caching
- [ ] Performance profiling
- [ ] Integration with graphics API for runtime updates
