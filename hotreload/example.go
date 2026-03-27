package hotreload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ============================================================
// DEMO/EXAMPLE USAGE
// ============================================================

// Example demonstrates shader hot reloading usage
type Example struct {
	outputWriter io.Writer
}

// NewExample creates a new example
func NewExample(out io.Writer) *Example {
	if out == nil {
		out = os.Stdout
	}
	return &Example{outputWriter: out}
}

// logf logs a message
func (e *Example) logf(format string, args ...interface{}) {
	fmt.Fprintf(e.outputWriter, format+"\n", args...)
}

// RunBasicExample demonstrates basic hot reload setup
func (e *Example) RunBasicExample(shaderPath string) error {
	e.logf("=== Shader Hot Reload Example ===\n")

	// Create reloader with 500ms debounce
	reloader, err := NewShaderReloader("", 500*time.Millisecond)
	if err != nil {
		return fmt.Errorf("failed to create reloader: %w", err)
	}
	defer reloader.Close()

	// Get output directory
	outputDir := filepath.Dir(shaderPath)
	outputFile := filepath.Join(outputDir, "output.dxil")

	// Register shader
	err = reloader.RegisterShader(ShaderConfig{
		Path:        shaderPath,
		Name:        "MainShader",
		Profile:     "vs_6_0",
		Entry:       "main",
		OutputPath:  outputFile,
		AutoCompile: true,
	})
	if err != nil {
		return fmt.Errorf("failed to register shader: %w", err)
	}

	e.logf("Registered shader: %s", shaderPath)

	// Register reload callback
	err = reloader.OnReload("MainShader", func(event ShaderReloadEvent) error {
		if event.Success {
			e.logf("[%s] ✓ Compilation successful (%v, %d bytes)",
				event.Timestamp.Format("15:04:05"),
				event.Duration,
				len(event.Output))
		} else {
			e.logf("[%s] ✗ Compilation failed: %v",
				event.Timestamp.Format("15:04:05"),
				event.Error)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to register callback: %w", err)
	}

	// Start reloader
	err = reloader.Start()
	if err != nil {
		return fmt.Errorf("failed to start reloader: %w", err)
	}
	defer reloader.Stop()

	e.logf("Hot reload started. Monitoring for changes...")
	e.logf("Edit %s to trigger recompilation.\n", shaderPath)

	// Keep running for 30 seconds or until interrupted
	time.Sleep(30 * time.Second)

	return nil
}

// RunMultiShaderExample demonstrates monitoring multiple shaders
func (e *Example) RunMultiShaderExample(shaderDir string) error {
	e.logf("=== Multi-Shader Hot Reload Example ===\n")

	reloader, err := NewShaderReloader("", 300*time.Millisecond)
	if err != nil {
		return fmt.Errorf("failed to create reloader: %w", err)
	}
	defer reloader.Close()

	// Shader configurations
	configs := []ShaderConfig{
		{
			Path:       filepath.Join(shaderDir, "vertex.hlsl"),
			Name:       "VertexShader",
			Profile:    "vs_6_0",
			Entry:      "main",
			OutputPath: filepath.Join(shaderDir, "vertex.dxil"),
		},
		{
			Path:       filepath.Join(shaderDir, "pixel.hlsl"),
			Name:       "PixelShader",
			Profile:    "ps_6_0",
			Entry:      "main",
			OutputPath: filepath.Join(shaderDir, "pixel.dxil"),
		},
		{
			Path:       filepath.Join(shaderDir, "compute.hlsl"),
			Name:       "ComputeShader",
			Profile:    "cs_6_0",
			Entry:      "main",
			OutputPath: filepath.Join(shaderDir, "compute.dxil"),
		},
	}

	// Register all shaders
	for _, config := range configs {
		config.AutoCompile = true
		if err := reloader.RegisterShader(config); err != nil {
			e.logf("Warning: Failed to register %s: %v", config.Name, err)
			continue
		}
		e.logf("Registered: %s", config.Name)
	}

	// Register callbacks for each shader
	for _, config := range configs {
		name := config.Name
		err := reloader.OnReload(name, func(event ShaderReloadEvent) error {
			status := "✓"
			if !event.Success {
				status = "✗"
			}
			e.logf("[%s] %s %s (%v)",
				event.Timestamp.Format("15:04:05"),
				status,
				event.Name,
				event.Duration)
			return nil
		})
		if err != nil {
			e.logf("Warning: Failed to register callback for %s", name)
		}
	}

	// Start reloader
	if err := reloader.Start(); err != nil {
		return fmt.Errorf("failed to start reloader: %w", err)
	}
	defer reloader.Stop()

	e.logf("\nMonitoring %d shaders...\n", len(configs))

	// Keep running
	time.Sleep(1 * time.Minute)

	return nil
}

// RunDebugExample demonstrates debug/verbose output
func (e *Example) RunDebugExample(shaderPath string) error {
	e.logf("=== Debug Example ===\n")

	reloader, err := NewShaderReloader("", 200*time.Millisecond)
	if err != nil {
		return fmt.Errorf("failed to create reloader: %w", err)
	}
	defer reloader.Close()

	err = reloader.RegisterShader(ShaderConfig{
		Path:        shaderPath,
		Name:        "DebugShader",
		Profile:     "ps_6_0",
		Entry:       "main",
		Defines:     []string{"DEBUG=1"},
		AutoCompile: true,
	})
	if err != nil {
		return fmt.Errorf("failed to register shader: %w", err)
	}

	// Detailed callback with debug output
	err = reloader.OnReload("DebugShader", func(event ShaderReloadEvent) error {
		e.logf("\n--- Reload Event ---")
		e.logf("Time:     %s", event.Timestamp.Format("2006-01-02 15:04:05"))
		e.logf("File:     %s", event.Path)
		e.logf("Success:  %v", event.Success)
		e.logf("Duration: %v", event.Duration)
		if event.Success {
			e.logf("Output:   %d bytes", len(event.Output))
		} else {
			e.logf("Error:    %v", event.Error)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to register callback: %w", err)
	}

	if err := reloader.Start(); err != nil {
		return fmt.Errorf("failed to start reloader: %w", err)
	}
	defer reloader.Stop()

	e.logf("Debug mode active. Watching %s\n", shaderPath)

	// Manual recompilation example
	time.Sleep(2 * time.Second)
	e.logf("\nManually recompiling...")
	_ = reloader.Recompile("DebugShader")

	time.Sleep(10 * time.Second)

	return nil
}
