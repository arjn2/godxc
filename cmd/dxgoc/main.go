// DXC Pure Go - HLSL to DXIL Compiler
// Argument forwarding edition - all DXC flags supported
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	dxc "dxgoc/internal/dxc"
)

// preprocessShader processes #include directives in shader source code
// This bypasses COM include handler issues by resolving includes in Go
func preprocessShader(source []byte, includePaths []string, processedFiles map[string]bool) ([]byte, error) {
	var result bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader(source))

	// Match: #include "file" or #include <file>
	includeRe := regexp.MustCompile(`^\s*#include\s+[<"](.+?)[>"]`)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if line is an include directive
		if matches := includeRe.FindStringSubmatch(line); matches != nil {
			filename := matches[1]

			// Try to load the include file
			content, err := loadInclude(filename, includePaths)
			if err != nil {
				return nil, fmt.Errorf("cannot resolve include %q: %w", filename, err)
			}

			// Prevent infinite recursion by tracking processed files
			absPath, _ := filepath.Abs(filepath.Join(includePaths[0], filename))
			if processedFiles[absPath] {
				// Skip if already processed
				continue
			}
			processedFiles[absPath] = true

			// Recursively preprocess included files
			processed, err := preprocessShader(content, includePaths, processedFiles)
			if err != nil {
				return nil, err
			}

			result.Write(processed)
		} else {
			// Write non-include lines as-is
			result.WriteString(line + "\n")
		}
	}

	return result.Bytes(), scanner.Err()
}

// loadInclude searches for an include file in the provided include paths
func loadInclude(filename string, includePaths []string) ([]byte, error) {
	// Try each include path
	for _, path := range includePaths {
		fullPath := filepath.Join(path, filename)
		if content, err := ioutil.ReadFile(fullPath); err == nil {
			return content, nil
		}
	}

	// Try relative to current directory
	if content, err := ioutil.ReadFile(filename); err == nil {
		return content, nil
	}

	return nil, fmt.Errorf("file not found: %s", filename)
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	// Parse arguments
	args, err := ParseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Handle special actions
	if args.IsAction() {
		action := args.GetAction()
		switch action {
		case "version":
			fmt.Println("dxgoc version 1.0.0 (Pure Go, No CGO)")
			fmt.Println("DirectX Shader Compiler (DXC) wrapper")
			os.Exit(0)
		case "check":
			fmt.Println("Checking DXC availability...")
			if err := dxc.InitCOM(); err != nil {
				fmt.Printf("COM init failed: %v\n", err)
				os.Exit(1)
			}
			compiler, err := dxc.NewCompiler("")
			if err != nil {
				fmt.Printf("DXC not available: %v\n", err)
				os.Exit(1)
			}
			compiler.Close()
			fmt.Println("✓ DXC is available and working")
			os.Exit(0)
		case "profiles":
			fmt.Println("Supported shader profiles:")
			fmt.Println("  Vertex:   vs_6_0, vs_6_1, vs_6_2, vs_6_3, vs_6_4, vs_6_5, vs_6_6, vs_6_7, vs_6_8")
			fmt.Println("  Pixel:    ps_6_0, ps_6_1, ps_6_2, ps_6_3, ps_6_4, ps_6_5, ps_6_6, ps_6_7, ps_6_8")
			fmt.Println("  Compute:  cs_6_0, cs_6_1, cs_6_2, cs_6_3, cs_6_4, cs_6_5, cs_6_6, cs_6_7, cs_6_8")
			fmt.Println("  Geometry: gs_6_0, gs_6_1, gs_6_2, gs_6_3, gs_6_4, gs_6_5, gs_6_6, gs_6_7, gs_6_8")
			fmt.Println("  Hull:     hs_6_0, hs_6_1, hs_6_2, hs_6_3, hs_6_4, hs_6_5, hs_6_6, hs_6_7, hs_6_8")
			fmt.Println("  Domain:   ds_6_0, ds_6_1, ds_6_2, ds_6_3, ds_6_4, ds_6_5, ds_6_6, ds_6_7, ds_6_8")
			os.Exit(0)
		case "diagnostic":
			fmt.Println("Diagnostic Information:")
			if err := dxc.InitCOM(); err != nil {
				fmt.Printf("  COM: FAILED - %v\n", err)
			} else {
				fmt.Println("  COM: OK")
			}
			compiler, err := dxc.NewCompiler("")
			if err != nil {
				fmt.Printf("  DXC: FAILED - %v\n", err)
			} else {
				fmt.Println("  DXC: OK")
				compiler.Close()
			}
			os.Exit(0)
		}
	}

	// Validate arguments
	if err := args.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Show usage if needed
	if args.ShouldShowUsage() {
		showUsage()
		os.Exit(1)
	}

	// Enable debug if requested
	if args.Debug {
		dxc.SetDebug(true)
	}

	// Print parsed arguments
	args.PrintParsedInfo()

	// Initialize COM
	if err := dxc.InitCOM(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: COM initialization failed: %v\n", err)
		os.Exit(1)
	}

	// Create compiler
	fmt.Println("\nCreating DXC compiler instance...")
	compiler, err := dxc.NewCompiler("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create DXC compiler: %v\n", err)
		os.Exit(1)
	}
	defer compiler.Close()
	fmt.Println("✓ Compiler created successfully")

	// Read shader source
	fmt.Printf("\nReading shader: %s\n", args.InputFile)
	source, err := ioutil.ReadFile(args.InputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to read shader: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Read %d bytes\n", len(source))

	// Preprocess shader to handle #include directives
	// This avoids COM include handler complexity
	if len(args.IncludePaths) > 0 {
		fmt.Println("\nPreprocessing shader includes...")
		allPaths := args.IncludePaths
		// Also add shader directory to include paths
		shaderDir := filepath.Dir(args.InputFile)
		allPaths = append([]string{shaderDir}, allPaths...)

		processed, err := preprocessShader(source, allPaths, make(map[string]bool))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error preprocessing shader: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("  Original: %d bytes\n", len(source))
		fmt.Printf("  Processed: %d bytes (includes expanded)\n", len(processed))
		source = processed
	}

	// Convert arguments to DXC format
	dxcArgs := args.ConvertToDXCArgs()
	fmt.Printf("\nForwarding %d arguments to DXC...\n", len(dxcArgs))

	// Show converted arguments
	fmt.Println("  DXC Arguments:")
	for i, arg := range dxcArgs {
		fmt.Printf("    [%d] %s\n", i, arg)
	}

	// Compile
	fmt.Println("\nCompiling shader...")
	result, err := compiler.CompileRaw(source, dxcArgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Compilation failed: %v\n", err)
		os.Exit(1)
	}

	if !result.Success {
		fmt.Fprintf(os.Stderr, "Error: Compilation failed\n")
		if len(result.Errors) > 0 {
			fmt.Fprintf(os.Stderr, "Errors:\n%s\n", string(result.Errors))
		}
		os.Exit(1)
	}

	fmt.Println("✓ Compilation successful")
	fmt.Printf("✓ Output: %d bytes\n", len(result.Output))

	// Write output
	outputPath := args.GetOutputFile()
	fmt.Printf("\nWriting output: %s\n", outputPath)
	if err := ioutil.WriteFile(outputPath, result.Output, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to write output: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Output written successfully")

	// Show error messages if any
	if len(result.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n%s\n", string(result.Warnings))
	}

	fmt.Println("\n✓ Done!")
}

func showUsage() {
	fmt.Println(`dxgoc - DirectX Shader Compiler (DXC) Pure Go Wrapper

USAGE:
  dxgoc [global-options] -i <file> [dxc-options]

GLOBAL OPTIONS:
  -version    Show version information
  -check      Check if DXC is available
  -profiles   Show supported shader profiles
  -diag       Show diagnostic information
  -debug      Enable debug output
  -dll <path> Use specific DXC DLL (default: auto-detect)

DXC OPTIONS (forwarded directly to compiler):
  -i <file>   Input HLSL file (required)
  -Fo <file>  Output file (default: input.dxil or input.spv)
  -T <prof>   Target profile (e.g., vs_6_0, ps_6_0, cs_6_0)
  -E <entry>  Entry point function (default: main)
  -I <path>   Include path (can specify multiple times)
  -D <define> Define preprocessor symbol (can specify multiple times)
  -spirv      Compile to SPIR-V instead of DXIL

EXAMPLES:
  # Compile vertex shader to DXIL
  dxgoc -i shader.hlsl -T vs_6_0 -E main

  # Compile to SPIR-V (Vulkan)
  dxgoc -i shader.hlsl -T vs_6_0 -E main -spirv

  # With includes and defines
  dxgoc -i shader.hlsl -T vs_6_0 -E main -I ./headers -D DEBUG=1

  # Specify output file
  dxgoc -i shader.hlsl -T vs_6_0 -Fo output.dxil

NOTES:
  - Includes (#include) are processed in Go before compilation
  - No CGO required - pure Go implementation
  - All DXC arguments are supported via forwarding
`)
}
