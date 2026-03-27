// Package dxc - DirectX Shader Compiler Go bindings
// Diagnostic tools for debugging interface issues
package dxc

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

// DiagnosticResult contains diagnostic information
type DiagnosticResult struct {
	DLLPath         string
	DLLVersion      string
	DxilPath        string
	UtilsIID        string
	UtilsSuccess    bool
	CompilerIID     string
	CompilerSuccess bool
	CompilerCLSID   string
	AllResults      []IIDProbeResult
	VTableEntries   []VTableEntry
}

// IIDProbeResult shows the result of probing a single IID
type IIDProbeResult struct {
	Interface string
	CLSID     string
	IIDName   string
	IID       string
	HRESULT   string
	Success   bool
}

// VTableEntry stores vtable entry for debugging
type VTableEntry struct {
	Index int
	Addr  string
}

// Diagnose runs full diagnostic on DXC setup
func Diagnose(dllPath string) (*DiagnosticResult, error) {
	result := &DiagnosticResult{}

	// Init COM
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

	result.DLLPath = dllPath
	result.DLLVersion, _ = GetDLLVersion(dllPath)

	// Check dxil.dll
	dxilPath := filepath.Join(filepath.Dir(dllPath), "dxil.dll")
	if _, err := os.Stat(dxilPath); err == nil {
		result.DxilPath = dxilPath
	}

	// Load DLL
	dll, err := syscall.LoadDLL(dllPath)
	if err != nil {
		return nil, fmt.Errorf("load dll failed: %v", err)
	}
	defer dll.Release()

	proc, err := dll.FindProc("DxcCreateInstance")
	if err != nil {
		return nil, fmt.Errorf("DxcCreateInstance not found")
	}

	debugPrint("============================================================")
	debugPrint("DXC DIAGNOSTIC")
	debugPrint("============================================================")
	debugPrint("DLL: %s", dllPath)
	debugPrint("Version: %s", result.DLLVersion)
	debugPrint("dxil.dll: %s", result.DxilPath)
	debugPrint("------------------------------------------------------------")

	// ========================================
	// Probe IDxcUtils with CLSID_DxcLibrary
	// ========================================
	debugPrint("Probing IDxcUtils (CLSID_DxcLibrary):")
	for _, iid := range UtilsIIDs {
		var ptr unsafe.Pointer
		hr, _, _ := proc.Call(
			uintptr(unsafe.Pointer(&CLSID_DxcLibrary)),
			uintptr(unsafe.Pointer(&iid.GUID)),
			uintptr(unsafe.Pointer(&ptr)),
		)

		probeResult := IIDProbeResult{
			Interface: iid.Name,
			CLSID:     "CLSID_DxcLibrary",
			IIDName:   iid.Name,
			IID:       GUIDString(&iid.GUID),
			HRESULT:   fmt.Sprintf("0x%08X", hr),
			Success:   hr == 0 && ptr != nil,
		}
		result.AllResults = append(result.AllResults, probeResult)
		debugPrint("  %s: hr=%s, success=%v", iid.Name, probeResult.HRESULT, probeResult.Success)

		if probeResult.Success && !result.UtilsSuccess {
			result.UtilsSuccess = true
			result.UtilsIID = iid.Name

			// Print vtable
			vtbl := *(*uintptr)(ptr)
			debugPrint("  VTable at: 0x%X", vtbl)
			for i := 0; i < 15; i++ {
				addr := *(*uintptr)(unsafe.Pointer(vtbl + uintptr(i)*unsafe.Sizeof(uintptr(0))))
				result.VTableEntries = append(result.VTableEntries, VTableEntry{
					Index: i,
					Addr:  fmt.Sprintf("0x%X", addr),
				})
				debugPrint("    [%2d] 0x%X", i, addr)
			}

			// Try QueryInterface for IDxcCompiler from this object
			debugPrint("  Trying QueryInterface for IDxcCompiler from IDxcUtils...")
			for _, ciid := range CompilerIIDs {
				cptr, chr := QueryInterface(ptr, &ciid.GUID)
				debugPrint("    %s: hr=0x%X, ptr=0x%X", ciid.Name, chr, cptr)
				if chr == 0 && cptr != nil {
					result.CompilerSuccess = true
					result.CompilerIID = ciid.Name
					result.CompilerCLSID = "via QueryInterface"
					Release(cptr)
				}
			}

			Release(ptr)
		}
	}

	// ========================================
	// Probe IDxcCompiler with CLSID_DxcCompiler (CORRECT CLSID!)
	// ========================================
	debugPrint("------------------------------------------------------------")
	debugPrint("Probing IDxcCompiler (CLSID_DxcCompiler - CORRECT!):")
	for _, iid := range CompilerIIDs {
		var ptr unsafe.Pointer
		hr, _, _ := proc.Call(
			uintptr(unsafe.Pointer(&CLSID_DxcCompiler)),
			uintptr(unsafe.Pointer(&iid.GUID)),
			uintptr(unsafe.Pointer(&ptr)),
		)

		probeResult := IIDProbeResult{
			Interface: iid.Name,
			CLSID:     "CLSID_DxcCompiler",
			IIDName:   iid.Name,
			IID:       GUIDString(&iid.GUID),
			HRESULT:   fmt.Sprintf("0x%08X", hr),
			Success:   hr == 0 && ptr != nil,
		}
		result.AllResults = append(result.AllResults, probeResult)
		debugPrint("  %s: hr=%s, success=%v", iid.Name, probeResult.HRESULT, probeResult.Success)

		if probeResult.Success && !result.CompilerSuccess {
			result.CompilerSuccess = true
			result.CompilerIID = iid.Name
			result.CompilerCLSID = "CLSID_DxcCompiler"
			Release(ptr)
		}
	}

	// ========================================
	// Probe IDxcCompiler with CLSID_DxcLibrary (fallback)
	// ========================================
	debugPrint("------------------------------------------------------------")
	debugPrint("Probing IDxcCompiler (CLSID_DxcLibrary - fallback):")
	for _, iid := range CompilerIIDs {
		var ptr unsafe.Pointer
		hr, _, _ := proc.Call(
			uintptr(unsafe.Pointer(&CLSID_DxcLibrary)),
			uintptr(unsafe.Pointer(&iid.GUID)),
			uintptr(unsafe.Pointer(&ptr)),
		)

		probeResult := IIDProbeResult{
			Interface: iid.Name,
			CLSID:     "CLSID_DxcLibrary",
			IIDName:   iid.Name,
			IID:       GUIDString(&iid.GUID),
			HRESULT:   fmt.Sprintf("0x%08X", hr),
			Success:   hr == 0 && ptr != nil,
		}
		result.AllResults = append(result.AllResults, probeResult)
		debugPrint("  %s: hr=%s, success=%v", iid.Name, probeResult.HRESULT, probeResult.Success)

		if probeResult.Success && !result.CompilerSuccess {
			result.CompilerSuccess = true
			result.CompilerIID = iid.Name
			result.CompilerCLSID = "CLSID_DxcLibrary"
			Release(ptr)
		}
	}

	debugPrint("------------------------------------------------------------")
	debugPrint("SUMMARY:")
	debugPrint("  IDxcUtils:    %v (%s)", result.UtilsSuccess, result.UtilsIID)
	debugPrint("  IDxcCompiler: %v (%s via %s)", result.CompilerSuccess, result.CompilerIID, result.CompilerCLSID)
	debugPrint("============================================================")

	return result, nil
}

// PrintDiagnostic prints diagnostic results to stdout
func PrintDiagnostic(result *DiagnosticResult) {
	fmt.Println("============================================================")
	fmt.Println("DXC DIAGNOSTIC REPORT")
	fmt.Println("============================================================")
	fmt.Printf("DLL: %s\n", result.DLLPath)
	fmt.Printf("Version: %s\n", result.DLLVersion)
	fmt.Printf("dxil.dll: %s\n", result.DxilPath)
	fmt.Println("------------------------------------------------------------")
	fmt.Println("IID Probing Results:")
	fmt.Println("------------------------------------------------------------")

	for _, r := range result.AllResults {
		status := "FAILED"
		if r.Success {
			status = "OK"
		}
		fmt.Printf("  %-20s %-20s %-10s %s (%s)\n", r.Interface, r.CLSID, status, r.IID, r.HRESULT)
	}

	if len(result.VTableEntries) > 0 {
		fmt.Println("------------------------------------------------------------")
		fmt.Println("IDxcUtils VTable:")
		for _, e := range result.VTableEntries {
			fmt.Printf("  [%2d] %s\n", e.Index, e.Addr)
		}
	}

	fmt.Println("------------------------------------------------------------")
	fmt.Println("Summary:")
	fmt.Printf("  IDxcUtils:    %v (%s)\n", result.UtilsSuccess, result.UtilsIID)
	fmt.Printf("  IDxcCompiler: %v (%s via %s)\n", result.CompilerSuccess, result.CompilerIID, result.CompilerCLSID)
	fmt.Println("============================================================")

	// Provide analysis
	if result.UtilsSuccess && !result.CompilerSuccess {
		fmt.Println()
		fmt.Println("ANALYSIS:")
		fmt.Println("  IDxcUtils is available but IDxcCompiler is NOT.")
		fmt.Println()
		fmt.Println("  CLSIDs tried:")
		fmt.Println("    CLSID_DxcCompiler: {73E22D93-E6CE-47F3-B5BF-F0664F39C1B0}")
		fmt.Println("    CLSID_DxcLibrary:  {6245D6AF-66E0-48FD-80B4-4D271796748C}")
		fmt.Println()
		fmt.Println("  If all CLSIDs fail for IDxcCompiler:")
		fmt.Println("  - This DXC build may not expose programmatic compilation")
		fmt.Println("  - Try using dxc.exe command-line instead")
		fmt.Println("  - Check if you have the official Microsoft DXC")
	}
}
