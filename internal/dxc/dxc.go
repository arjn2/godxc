// Package dxc - DirectX Shader Compiler Go bindings
// Main entry point for DXC compilation
package dxc

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

// ============================================================
// COMPILER - Main DXC wrapper
// ============================================================

// Compiler wraps DXC functionality
type Compiler struct {
	dll      *syscall.DLL
	dxilDLL  *syscall.DLL
	proc     *syscall.Proc
	utils    *DxcUtils
	compiler *DxcCompiler
	dllPath  string
	dxilPath string
}

// NewCompiler creates a new DXC compiler instance
func NewCompiler(dllPath string) (*Compiler, error) {
	// Initialize COM
	if err := InitCOM(); err != nil {
		return nil, err
	}

	var err error
	if dllPath == "" {
		dllPath, err = FindDLL("dxcompiler.dll")
		if err != nil {
			return nil, err
		}
	}

	debugPrint("============================================================")
	debugPrint("DXC Go Compiler")
	debugPrint("============================================================")
	debugPrint("Loading: %s", dllPath)

	// Load dxcompiler.dll
	dll, err := syscall.LoadDLL(dllPath)
	if err != nil {
		return nil, fmt.Errorf("load dll failed: %v", err)
	}

	// Get DxcCreateInstance function
	proc, err := dll.FindProc("DxcCreateInstance")
	if err != nil {
		dll.Release()
		return nil, fmt.Errorf("DxcCreateInstance not found in dll")
	}

	debugPrint("DxcCreateInstance: 0x%X", proc.Addr())

	c := &Compiler{
		dll:     dll,
		proc:    proc,
		dllPath: dllPath,
	}

	// Load dxil.dll (required for DXIL validation)
	dxilPath := filepath.Join(filepath.Dir(dllPath), "dxil.dll")
	if _, err := os.Stat(dxilPath); err == nil {
		c.dxilDLL, _ = syscall.LoadDLL(dxilPath)
		c.dxilPath = dxilPath
		debugPrint("dxil.dll: %s", dxilPath)
	}

	// ========================================
	// STEP 1: Create IDxcUtils
	// ========================================
	debugPrint("------------------------------------------------------------")
	debugPrint("Creating IDxcUtils...")
	c.utils, err = NewDxcUtils(proc)
	if err != nil {
		dll.Release()
		return nil, fmt.Errorf("IDxcUtils: %v", err)
	}
	debugPrint("IDxcUtils: OK")
	c.utils.PrintVTable()

	// ========================================
	// STEP 2: Create IDxcCompiler3
	// ========================================
	debugPrint("------------------------------------------------------------")
	debugPrint("Creating IDxcCompiler3...")
	c.compiler, err = NewDxcCompiler(proc, c.utils)
	if err != nil {
		debugPrint("IDxcCompiler3: FAILED - %v", err)
		debugPrint("")
		debugPrint("WARNING: This DXC build does not expose IDxcCompiler3!")
		debugPrint("  Possible causes:")
		debugPrint("  - Vulkan SDK DXC may require different IIDs")
		debugPrint("  - This DXC build may be command-line only")
		debugPrint("")
		c.utils.Close()
		dll.Release()
		return nil, fmt.Errorf("IDxcCompiler3 not available - this DXC build may not support programmatic compilation")
	}
	debugPrint("IDxcCompiler3: OK")

	debugPrint("------------------------------------------------------------")
	debugPrint("Compiler ready!")
	debugPrint("============================================================")

	return c, nil
}

// Close releases all resources
func (c *Compiler) Close() {
	if c.compiler != nil {
		c.compiler.Close()
	}
	if c.utils != nil {
		c.utils.Close()
	}
	if c.dxilDLL != nil {
		c.dxilDLL.Release()
	}
	if c.dll != nil {
		c.dll.Release()
	}
}

// GetDLLPath returns the path to dxcompiler.dll
func (c *Compiler) GetDLLPath() string { return c.dllPath }

// GetDxilPath returns the path to dxil.dll
func (c *Compiler) GetDxilPath() string { return c.dxilPath }

// GetVersion returns DXC version string
func (c *Compiler) GetVersion() string {
	v, _ := GetDLLVersion(c.dllPath)
	return v
}

// Compile compiles HLSL source code with raw arguments
func (c *Compiler) Compile(source []byte, args []string) (*CompileResult, error) {
	if c.compiler == nil {
		return nil, fmt.Errorf("compiler not initialized")
	}
	return c.compiler.Compile(source, args)
}

// CompileSimple compiles HLSL and returns output
func (c *Compiler) CompileSimple(source, entry, profile string) ([]byte, error) {
	if c.compiler == nil {
		return nil, fmt.Errorf("compiler not initialized")
	}
	return c.compiler.CompileSimple(source, entry, profile)
}

// CompileToSpirv compiles HLSL to SPIR-V
func (c *Compiler) CompileToSpirv(source, entry, profile string) ([]byte, error) {
	if c.compiler == nil {
		return nil, fmt.Errorf("compiler not initialized")
	}
	return c.compiler.CompileToSpirv(source, entry, profile)
}

// CompileRaw compiles HLSL with raw arguments (forwarded directly to DXC)
func (c *Compiler) CompileRaw(source []byte, args []string) (*CompileResult, error) {
	if c.compiler == nil {
		return nil, fmt.Errorf("compiler not initialized")
	}
	return c.compiler.CompileRaw(source, args)
}

// ============================================================
// DLL UTILITIES
// ============================================================

// CheckDLL checks if dxcompiler.dll is available
func CheckDLL() (bool, string) {
	p, err := FindDLL("dxcompiler.dll")
	if err != nil {
		return false, err.Error()
	}
	return true, fmt.Sprintf("found: %s", p)
}

// FindDLL locates a DLL in common locations
func FindDLL(name string) (string, error) {
	var paths []string

	// Check alongside dxc.exe if in PATH
	if p, err := exec.LookPath("dxc.exe"); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(p), name))
	}

	// Check common relative paths
	paths = append(paths,
		name,
		filepath.Join("bin", name),
		filepath.Join("lib", name),
	)

	// Check alongside current executable
	var exePath [syscall.MAX_PATH]uint16
	k32 := syscall.NewLazyDLL("kernel32.dll")
	gmf := k32.NewProc("GetModuleFileNameW")
	ret, _, _ := gmf.Call(0, uintptr(unsafe.Pointer(&exePath[0])), syscall.MAX_PATH)
	if ret > 0 {
		exeDir := filepath.Dir(syscall.UTF16ToString(exePath[:ret]))
		paths = append(paths, filepath.Join(exeDir, name))
	}

	// Check standard DXC installation
	paths = append(paths,
		filepath.Join("C:\\Program Files\\Microsoft\\DirectXShaderCompiler\\bin", name),
	)

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("%s not found in PATH or common locations", name)
}

// GetDLLVersion extracts version info from DLL
func GetDLLVersion(dllPath string) (string, error) {
	if dllPath == "" {
		var err error
		dllPath, err = FindDLL("dxcompiler.dll")
		if err != nil {
			return "", err
		}
	}

	p, _ := syscall.UTF16PtrFromString(dllPath)
	buf := make([]byte, 0x1000)

	ver := syscall.NewLazyDLL("version.dll")
	ret, _, _ := ver.NewProc("GetFileVersionInfoW").Call(
		uintptr(unsafe.Pointer(p)),
		0,
		uintptr(len(buf)),
		uintptr(unsafe.Pointer(&buf[0])),
	)

	if ret == 0 {
		return "", fmt.Errorf("GetFileVersionInfo failed")
	}

	q, _ := syscall.UTF16PtrFromString("\\")
	var fi *byte
	var fl uint32

	ver.NewProc("VerQueryValueW").Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(q)),
		uintptr(unsafe.Pointer(&fi)),
		uintptr(unsafe.Pointer(&fl)),
	)

	if fi != nil && binary.LittleEndian.Uint32(buf[0:4]) == 0xFEEF04BD {
		return fmt.Sprintf("%d.%d.%d.%d",
			binary.LittleEndian.Uint32(buf[8:12]),
			binary.LittleEndian.Uint32(buf[12:16]),
			binary.LittleEndian.Uint32(buf[16:20]),
			binary.LittleEndian.Uint32(buf[20:24])), nil
	}

	return "unknown", nil
}
