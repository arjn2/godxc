// Argument parsing and validation for DXC Pure Go compiler
package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ============================================================
// ARGUMENT STRUCTURE
// ============================================================

// Args holds parsed command-line arguments
type Args struct {
	// Global options (not forwarded to DXC)
	Version bool
	Check   bool
	Diag    bool
	Profiles bool
	Debug   bool
	DLLPath string

	// DXC arguments (forwarded directly)
	DXCArgs []string

	// Parsed DXC arguments for convenience
	InputFile string
	OutputFile string
	TargetProfile string
	EntryPoint string
	IncludePaths []string
	Defines []string
	IsSpirV bool
}

// ============================================================
// ARGUMENT PARSING
// ============================================================

// ParseArgs parses command-line arguments into global options and DXC args
func ParseArgs(osArgs []string) (*Args, error) {
	args := &Args{
		DXCArgs: []string{},
		IncludePaths: []string{},
		Defines: []string{},
	}

	// Separate global options from DXC arguments
	globalOptions := []string{}
	dxcStartIdx := -1

	for i := 1; i < len(osArgs); i++ {
		arg := osArgs[i]

		// Check if this is a global option
		if isGlobalOption(arg) {
			globalOptions = append(globalOptions, arg)
			
			// Handle options that take a value
			if arg == "-dll" && i+1 < len(osArgs) {
				i++
				globalOptions = append(globalOptions, osArgs[i])
			}
		} else {
			// First non-global option marks the start of DXC args
			dxcStartIdx = i
			break
		}
	}

	// Parse global options
	if err := parseGlobalOptions(globalOptions, args); err != nil {
		return nil, err
	}

	// Collect DXC arguments
	if dxcStartIdx > 0 {
		args.DXCArgs = osArgs[dxcStartIdx:]
	}

	// Parse and validate DXC arguments
	if len(args.DXCArgs) > 0 {
		if err := parseDXCArgs(args); err != nil {
			return nil, err
		}
	}

	return args, nil
}

// ============================================================
// GLOBAL OPTIONS PARSING
// ============================================================

// isGlobalOption checks if an argument is one of our global options
func isGlobalOption(arg string) bool {
	return arg == "-version" ||
		arg == "-check" ||
		arg == "-diag" ||
		arg == "-profiles" ||
		arg == "-debug" ||
		arg == "-dll"
}

// parseGlobalOptions parses global command-line options
func parseGlobalOptions(globalArgs []string, args *Args) error {
	for i := 0; i < len(globalArgs); i++ {
		opt := globalArgs[i]

		switch opt {
		case "-version":
			args.Version = true
		case "-check":
			args.Check = true
		case "-diag":
			args.Diag = true
		case "-profiles":
			args.Profiles = true
		case "-debug":
			args.Debug = true
		case "-dll":
			if i+1 >= len(globalArgs) {
				return fmt.Errorf("-dll requires a path argument")
			}
			i++
			args.DLLPath = globalArgs[i]
		}
	}

	return nil
}

// ============================================================
// DXC ARGUMENTS PARSING
// ============================================================

// parseDXCArgs parses DXC arguments and extracts useful information
func parseDXCArgs(args *Args) error {
	for i := 0; i < len(args.DXCArgs); i++ {
		arg := args.DXCArgs[i]

		switch arg {
		case "-i":
			if i+1 >= len(args.DXCArgs) {
				return fmt.Errorf("-i requires an input file")
			}
			i++
			args.InputFile = args.DXCArgs[i]

		case "-Fo":
			if i+1 >= len(args.DXCArgs) {
				return fmt.Errorf("-Fo requires an output file")
			}
			i++
			args.OutputFile = args.DXCArgs[i]

		case "-T":
			if i+1 >= len(args.DXCArgs) {
				return fmt.Errorf("-T requires a target profile")
			}
			i++
			args.TargetProfile = args.DXCArgs[i]

		case "-E":
			if i+1 >= len(args.DXCArgs) {
				return fmt.Errorf("-E requires an entry point")
			}
			i++
			args.EntryPoint = args.DXCArgs[i]

		case "-I":
			if i+1 >= len(args.DXCArgs) {
				return fmt.Errorf("-I requires an include path")
			}
			i++
			args.IncludePaths = append(args.IncludePaths, args.DXCArgs[i])

		case "-D":
			if i+1 >= len(args.DXCArgs) {
				return fmt.Errorf("-D requires a definition")
			}
			i++
			args.Defines = append(args.Defines, args.DXCArgs[i])

		case "-spirv":
			args.IsSpirV = true
		}
	}

	return nil
}

// ============================================================
// VALIDATION
// ============================================================

// Validate checks that required arguments are present
func (a *Args) Validate() error {
	// If no DXC args provided, that's ok - we'll show usage
	if len(a.DXCArgs) == 0 {
		return nil
	}

	// Input file is required
	if a.InputFile == "" {
		return fmt.Errorf("input file required (-i <file>)")
	}

	return nil
}

// ============================================================
// UTILITY FUNCTIONS
// ============================================================

// HasInputFile checks if input file was provided
func (a *Args) HasInputFile() bool {
	return a.InputFile != ""
}

// ConvertToDXCArgs converts our dxgoc arguments to native DXC format
// DXC expects: dxc.exe <input-file> [options]
// Not: dxc.exe -i <input-file> [options]
func (a *Args) ConvertToDXCArgs() []string {
	result := []string{}

	// Start with input file as first positional argument (convert to absolute path)
	if a.InputFile != "" {
		inputPath := a.InputFile
		if !filepath.IsAbs(inputPath) {
			absPath, err := filepath.Abs(inputPath)
			if err == nil {
				inputPath = absPath
			}
		}
		result = append(result, inputPath)
	}

	// Always add shader directory as first include path (so headers in same dir can be found)
	if a.InputFile != "" {
		shaderDir := filepath.Dir(a.InputFile)
		if !filepath.IsAbs(shaderDir) {
			absPath, err := filepath.Abs(shaderDir)
			if err == nil {
				shaderDir = absPath
			}
		}
		result = append(result, "-I")
		result = append(result, shaderDir)
	}

	// Add all other DXC arguments except -i <file>
	for i := 0; i < len(a.DXCArgs); i++ {
		arg := a.DXCArgs[i]

		// Skip -i and its argument
		if arg == "-i" {
			if i+1 < len(a.DXCArgs) {
				i++ // skip the filename
			}
			continue
		}

		// Convert relative include paths to absolute paths
		if arg == "-I" {
			result = append(result, arg)
			if i+1 < len(a.DXCArgs) {
				i++
				includePath := a.DXCArgs[i]
				// Convert to absolute path if it's relative
				if !filepath.IsAbs(includePath) {
					absPath, err := filepath.Abs(includePath)
					if err == nil {
						includePath = absPath
					}
				}
				result = append(result, includePath)
			}
			continue
		}

		result = append(result, arg)
	}

	return result
}

// GetOutputFile returns the output file path or generates one
func (a *Args) GetOutputFile() string {
	if a.OutputFile != "" {
		return a.OutputFile
	}

	// Generate default based on input file
	if a.InputFile == "" {
		return "output.dxil"
	}

	// Remove extension from input file
	idx := strings.LastIndex(a.InputFile, ".")
	base := a.InputFile
	if idx > 0 {
		base = a.InputFile[:idx]
	}

	// Determine extension based on target
	ext := ".dxil"
	if a.IsSpirV {
		ext = ".spv"
	}

	return base + ext
}

// PrintParsedInfo prints information about parsed arguments
func (a *Args) PrintParsedInfo() {
	if len(a.DXCArgs) == 0 {
		return
	}

	fmt.Println("\nParsed DXC Arguments:")
	if a.InputFile != "" {
		fmt.Printf("  Input:    %s\n", a.InputFile)
	}
	if a.TargetProfile != "" {
		fmt.Printf("  Target:   %s\n", a.TargetProfile)
	}
	if a.EntryPoint != "" {
		fmt.Printf("  Entry:    %s\n", a.EntryPoint)
	}
	if a.OutputFile != "" {
		fmt.Printf("  Output:   %s\n", a.OutputFile)
	}
	if len(a.IncludePaths) > 0 {
		fmt.Printf("  Includes: %s\n", strings.Join(a.IncludePaths, ", "))
	}
	if len(a.Defines) > 0 {
		fmt.Printf("  Defines:  %s\n", strings.Join(a.Defines, ", "))
	}
	if a.IsSpirV {
		fmt.Println("  Mode:     SPIR-V")
	}
}

// ============================================================
// COMMAND-LINE INFO
// ============================================================

// ShouldShowUsage determines if usage should be shown
func (a *Args) ShouldShowUsage() bool {
	// Show usage if no args provided
	if len(a.DXCArgs) == 0 {
		return true
	}

	// Show usage if input file is missing
	if a.InputFile == "" {
		return true
	}

	return false
}

// IsAction checks if this is a special action (not compilation)
func (a *Args) IsAction() bool {
	return a.Version || a.Check || a.Diag || a.Profiles
}

// GetAction returns a string describing the requested action
func (a *Args) GetAction() string {
	switch {
	case a.Version:
		return "version"
	case a.Check:
		return "check"
	case a.Diag:
		return "diagnostic"
	case a.Profiles:
		return "profiles"
	default:
		return "compile"
	}
}
