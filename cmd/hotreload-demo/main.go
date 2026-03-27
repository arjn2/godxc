package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"dxgoc/hotreload"
)

func main() {
	shaderPath := flag.String("shader", "", "Path to shader file to monitor")
	shaderDir := flag.String("dir", "", "Directory with shaders to monitor")
	debounce := flag.Duration("debounce", 500*time.Millisecond, "Debounce duration")
	demo := flag.String("demo", "basic", "Demo type: basic, multi, or debug")
	timeout := flag.Duration("timeout", 30*time.Second, "How long to run")
	flag.Parse()

	if *shaderPath == "" && *shaderDir == "" {
		fmt.Println("Shader Hot Reload Demo")
		fmt.Println("=====================\n")
		fmt.Println("Usage:")
		fmt.Println("  hotreload -shader <path>       # Monitor single shader")
		fmt.Println("  hotreload -dir <path>          # Monitor shader directory")
		fmt.Println("  hotreload -demo <type>         # Run demo: basic, multi, debug")
		fmt.Println("\nOptions:")
		fmt.Println("  -debounce <ms>  Debounce delay (default: 500ms)")
		fmt.Println("  -timeout <sec>  Run duration (default: 30s)")
		fmt.Println("\nExamples:")
		fmt.Println("  hotreload -shader shaders/vertex.hlsl")
		fmt.Println("  hotreload -dir shaders -timeout 2m")
		fmt.Println("  hotreload -demo multi -timeout 1m")
		os.Exit(1)
	}

	example := hotreload.NewExample(os.Stdout)

	var err error
	start := time.Now()

	switch *demo {
	case "multi":
		if *shaderDir == "" {
			*shaderDir = "shaders"
		}
		err = runMultiShaderDemo(example, *shaderDir, *debounce, *timeout)

	case "debug":
		if *shaderPath == "" {
			*shaderPath = "shaders/vertex.hlsl"
		}
		err = runDebugDemo(example, *shaderPath, *debounce, *timeout)

	default: // "basic"
		if *shaderPath == "" && *shaderDir != "" {
			// Use first .hlsl file in directory
			files, _ := filepath.Glob(filepath.Join(*shaderDir, "*.hlsl"))
			if len(files) > 0 {
				*shaderPath = files[0]
			}
		}
		if *shaderPath == "" {
			*shaderPath = "shaders/vertex.hlsl"
		}
		err = runBasicDemo(example, *shaderPath, *debounce, *timeout)
	}

	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("\n❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Demo completed in %v\n", elapsed)
}

func runBasicDemo(example *hotreload.Example, shaderPath string, debounce, timeout time.Duration) error {
	// Create reloader
	reloader, err := hotreload.NewShaderReloader("", debounce)
	if err != nil {
		return fmt.Errorf("failed to create reloader: %w", err)
	}
	defer reloader.Close()

	// Determine output path
	outputPath := filepath.Join(
		filepath.Dir(shaderPath),
		filepath.Base(shaderPath[:len(shaderPath)-len(filepath.Ext(shaderPath))])+".dxil",
	)

	// Register shader
	err = reloader.RegisterShader(hotreload.ShaderConfig{
		Path:        shaderPath,
		Name:        "MainShader",
		Profile:     "vs_6_0",
		Entry:       "main",
		OutputPath:  outputPath,
		AutoCompile: true,
	})
	if err != nil {
		return fmt.Errorf("failed to register shader: %w", err)
	}

	fmt.Printf("Registered: %s\n", shaderPath)

	// Register callback
	recompileCount := 0
	err = reloader.OnReload("MainShader", func(event hotreload.ShaderReloadEvent) error {
		recompileCount++
		if event.Success {
			fmt.Printf("[#%d] ✓ %s (%v, %d bytes)\n",
				recompileCount,
				event.Timestamp.Format("15:04:05"),
				event.Duration,
				len(event.Output))
		} else {
			fmt.Printf("[#%d] ✗ %s: %v\n",
				recompileCount,
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

	fmt.Printf("Hot reload started. Edit %s to trigger recompilation.\n", filepath.Base(shaderPath))
	fmt.Printf("Running for %v (Ctrl+C to exit)...\n\n", timeout)

	// Wait for timeout
	time.Sleep(timeout)

	fmt.Printf("\nRecompiled %d times\n", recompileCount)
	return nil
}

func runMultiShaderDemo(example *hotreload.Example, shaderDir string, debounce, timeout time.Duration) error {
	// Create reloader
	reloader, err := hotreload.NewShaderReloader("", debounce)
	if err != nil {
		return fmt.Errorf("failed to create reloader: %w", err)
	}
	defer reloader.Close()

	// Find shaders
	files, _ := filepath.Glob(filepath.Join(shaderDir, "*.hlsl"))
	if len(files) == 0 {
		return fmt.Errorf("no .hlsl files found in %s", shaderDir)
	}

	// Register shaders
	count := 0
	for _, path := range files {
		name := filepath.Base(path[:len(path)-5]) // remove .hlsl
		outputPath := filepath.Join(shaderDir, name+".dxil")

		err := reloader.RegisterShader(hotreload.ShaderConfig{
			Path:        path,
			Name:        name,
			Profile:     "vs_6_0",
			Entry:       "main",
			OutputPath:  outputPath,
			AutoCompile: true,
		})
		if err != nil {
			fmt.Printf("Warning: Failed to register %s: %v\n", name, err)
		} else {
			count++
			fmt.Printf("Registered: %s\n", name)
		}

		// Register callback
		shaderName := name
		reloader.OnReload(name, func(event hotreload.ShaderReloadEvent) error {
			if event.Success {
				fmt.Printf("  ✓ %s compiled in %v (%d bytes)\n",
					shaderName, event.Duration, len(event.Output))
			} else {
				fmt.Printf("  ✗ %s failed: %v\n", shaderName, event.Error)
			}
			return nil
		})
	}

	if count == 0 {
		return fmt.Errorf("no shaders registered")
	}

	// Start reloader
	err = reloader.Start()
	if err != nil {
		return fmt.Errorf("failed to start reloader: %w", err)
	}
	defer reloader.Stop()

	fmt.Printf("\nMonitoring %d shaders in %s\n", count, shaderDir)
	fmt.Printf("Running for %v...\n\n", timeout)

	time.Sleep(timeout)

	return nil
}

func runDebugDemo(example *hotreload.Example, shaderPath string, debounce, timeout time.Duration) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create reloader
	reloader, err := hotreload.NewShaderReloader("", debounce)
	if err != nil {
		return fmt.Errorf("failed to create reloader: %w", err)
	}
	defer reloader.Close()

	// Register with debug info
	err = reloader.RegisterShader(hotreload.ShaderConfig{
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

	fmt.Printf("Registered: %s (profile: ps_6_0)\n\n", shaderPath)

	// Detailed callback
	eventCount := 0
	err = reloader.OnReload("DebugShader", func(event hotreload.ShaderReloadEvent) error {
		eventCount++
		fmt.Printf("\n--- Reload Event #%d ---\n", eventCount)
		fmt.Printf("Time:     %s\n", event.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("File:     %s\n", event.Path)
		fmt.Printf("Success:  %v\n", event.Success)
		fmt.Printf("Duration: %v\n", event.Duration)

		if event.Success {
			fmt.Printf("Output:   %d bytes\n", len(event.Output))
		} else {
			fmt.Printf("Error:    %v\n", event.Error)
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

	fmt.Printf("Debug mode active. Edit %s to trigger recompilation.\n", filepath.Base(shaderPath))
	fmt.Printf("Running for %v...\n\n", timeout)

	// Manual recompile after 2 seconds
	time.AfterFunc(2*time.Second, func() {
		fmt.Println("\n[Manual] Recompiling...")
		_ = reloader.Recompile("DebugShader")
	})

	time.Sleep(timeout)

	fmt.Printf("\nTotal events: %d\n", eventCount)
	return nil
}
