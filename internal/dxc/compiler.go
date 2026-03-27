// Package dxc - DirectX Shader Compiler Go bindings
// IDxcCompiler3 interface handling
package dxc

import (
        "fmt"
        "syscall"
        "unsafe"
)

// ============================================================
// DXC BUFFER - Must match C++ struct layout exactly!
// ============================================================

// DxcBuffer - Source buffer for compilation
// C++ struct: { LPCVOID Ptr; SIZE_T Size; UINT Encoding; }
// On 64-bit: 8 + 8 + 4 + 4(padding) = 24 bytes
type DxcBuffer struct {
        Ptr      uintptr // 8 bytes - source pointer
        Size     uintptr // 8 bytes - SIZE_T is pointer-sized!
        Encoding uint32  // 4 bytes - code page (65001 = UTF-8, 0 = auto)
        _        uint32  // 4 bytes - padding for alignment
}

// ============================================================
// COMPILE RESULT
// ============================================================

// CompileResult holds compilation output
type CompileResult struct {
        Success  bool
        Output   []byte
        Errors   []byte
        Warnings []byte
}

// ============================================================
// IDXCRESULT INTERFACE (returned by IDxcCompiler3::Compile)
// ============================================================

// DxcResult wraps IDxcResult interface
type DxcResult struct {
        ptr  unsafe.Pointer
        vtbl uintptr
}

// NewDxcResult wraps an existing result pointer
func NewDxcResult(ptr unsafe.Pointer) *DxcResult {
        if ptr == nil {
                return nil
        }
        return &DxcResult{
                ptr:  ptr,
                vtbl: *(*uintptr)(ptr),
        }
}

// getVTableMethod returns the method pointer at the given vtable index
func (r *DxcResult) getVTableMethod(index int) uintptr {
        if r.vtbl == 0 {
                return 0
        }
        return *(*uintptr)(unsafe.Pointer(r.vtbl + uintptr(index)*unsafe.Sizeof(uintptr(0))))
}

// Close releases the result
func (r *DxcResult) Close() {
        if r.ptr != nil {
                Release(r.ptr)
                r.ptr = nil
        }
}

// GetStatus returns compilation status (0 = success)
func (r *DxcResult) GetStatus() int32 {
        var status int32
        method := r.getVTableMethod(3) // GetStatus is index 3
        syscall.Syscall(method, 2,
                uintptr(r.ptr),
                uintptr(unsafe.Pointer(&status)),
                0)
        return status
}

// GetResult returns the output blob
func (r *DxcResult) GetResult() ([]byte, error) {
        var blob unsafe.Pointer
        method := r.getVTableMethod(4) // GetResult is index 4
        hr, _, _ := syscall.Syscall(method, 2,
                uintptr(r.ptr),
                uintptr(unsafe.Pointer(&blob)),
                0)

        if hr != 0 || blob == nil {
                return nil, fmt.Errorf("GetResult failed: 0x%X", hr)
        }
        defer Release(blob)

        return ReadBlob(blob), nil
}

// GetErrors returns the error/warning blob
func (r *DxcResult) GetErrors() ([]byte, error) {
        var blob unsafe.Pointer
        method := r.getVTableMethod(5) // GetErrors is index 5
        hr, _, _ := syscall.Syscall(method, 2,
                uintptr(r.ptr),
                uintptr(unsafe.Pointer(&blob)),
                0)

        if hr != 0 || blob == nil {
                return nil, fmt.Errorf("GetErrors failed: 0x%X", hr)
        }
        defer Release(blob)

        return ReadBlob(blob), nil
}

// ============================================================
// IDXCCOMPILER3 INTERFACE
// ============================================================

// DxcCompiler wraps IDxcCompiler3 interface
type DxcCompiler struct {
        ptr   unsafe.Pointer
        vtbl  uintptr
        utils *DxcUtils  // Reference to utils for include handler creation
}

// NewDxcCompiler creates IDxcCompiler3 using multiple strategies
func NewDxcCompiler(proc *syscall.Proc, utils *DxcUtils) (*DxcCompiler, error) {
        debugPrint("Creating IDxcCompiler...")

        // Strategy 1: Direct DxcCreateInstance with CLSID_DxcCompiler (CORRECT CLSID!)
        debugPrint("  Strategy 1: DxcCreateInstance with CLSID_DxcCompiler")
        for _, iid := range CompilerIIDs {
                debugPrint("    Trying %s: %s", iid.Name, GUIDString(&iid.GUID))

                var ptr unsafe.Pointer
                hr, _, _ := proc.Call(
                        uintptr(unsafe.Pointer(&CLSID_DxcCompiler)), // USE CLSID_DxcCompiler!
                        uintptr(unsafe.Pointer(&iid.GUID)),
                        uintptr(unsafe.Pointer(&ptr)),
                )

                debugPrint("      hr=0x%X, ptr=0x%X", hr, ptr)

                if hr == 0 && ptr != nil {
                        compiler := &DxcCompiler{
                                ptr:   ptr,
                                vtbl:  *(*uintptr)(ptr),
                                utils: utils,
                        }
                        debugPrint("      -> SUCCESS: %s", iid.Name)
                        return compiler, nil
                }
        }

        // Strategy 2: Try CLSID_DxcLibrary (some versions use same CLSID)
        debugPrint("  Strategy 2: DxcCreateInstance with CLSID_DxcLibrary")
        for _, iid := range CompilerIIDs {
                debugPrint("    Trying %s: %s", iid.Name, GUIDString(&iid.GUID))

                var ptr unsafe.Pointer
                hr, _, _ := proc.Call(
                        uintptr(unsafe.Pointer(&CLSID_DxcLibrary)),
                        uintptr(unsafe.Pointer(&iid.GUID)),
                        uintptr(unsafe.Pointer(&ptr)),
                )

                debugPrint("      hr=0x%X, ptr=0x%X", hr, ptr)

                if hr == 0 && ptr != nil {
                        compiler := &DxcCompiler{
                                ptr:   ptr,
                                vtbl:  *(*uintptr)(ptr),
                                utils: utils,
                        }
                        debugPrint("      -> SUCCESS: %s", iid.Name)
                        return compiler, nil
                }
        }

        // Strategy 3: QueryInterface from IDxcUtils
        debugPrint("  Strategy 3: QueryInterface from IDxcUtils")
        if utils != nil && utils.ptr != nil {
                for _, iid := range CompilerIIDs {
                        debugPrint("    Trying %s: %s", iid.Name, GUIDString(&iid.GUID))

                        ptr, hr := utils.QueryInterface(&iid.GUID)
                        debugPrint("      hr=0x%X, ptr=0x%X", hr, ptr)

                        if hr == 0 && ptr != nil {
                                compiler := &DxcCompiler{
                                        ptr:   ptr,
                                        vtbl:  *(*uintptr)(ptr),
                                        utils: utils,
                                }
                                debugPrint("      -> SUCCESS: %s (via QueryInterface)", iid.Name)
                                return compiler, nil
                        }
                }
        }

        return nil, fmt.Errorf("IDxcCompiler creation failed - no matching IID found")
}

// Close releases IDxcCompiler
func (c *DxcCompiler) Close() {
        if c.ptr != nil {
                Release(c.ptr)
                c.ptr = nil
        }
}

// GetPtr returns the underlying pointer
func (c *DxcCompiler) GetPtr() unsafe.Pointer { return c.ptr }

// getVTableMethod returns the method pointer at the given vtable index
func (c *DxcCompiler) getVTableMethod(index int) uintptr {
        if c.vtbl == 0 {
                return 0
        }
        return *(*uintptr)(unsafe.Pointer(c.vtbl + uintptr(index)*unsafe.Sizeof(uintptr(0))))
}

// Compile compiles HLSL source to DXIL/SPIR-V using IDxcCompiler3::Compile
// Args are passed directly to DXC compiler (e.g., -T vs_6_0 -E main -spirv)
func (c *DxcCompiler) Compile(source []byte, args []string) (*CompileResult, error) {
        debugPrint("Compile: args=%v", args)
        debugPrint("Source length: %d", len(source))

        // Safety check
        if len(source) == 0 {
                return nil, fmt.Errorf("source is empty")
        }

        // Include handler is NULL - not needed!
        // The -I flags in args are handled by DXC's internal command-line processor.
        // This is the same mechanism used by command-line dxc.exe.
        // No COM interface CreateDefaultIncludeHandler needed (unknown VTable index, crashes).
        var includeHandler unsafe.Pointer = nil

        // Keep references to UTF16 strings to prevent GC
        utf16Args := make([][]uint16, len(args))
        argPtrs := make([]uintptr, len(args))
        for i, arg := range args {
                utf16Args[i] = syscall.StringToUTF16(arg)
                argPtrs[i] = uintptr(unsafe.Pointer(&utf16Args[i][0]))
                debugPrint("  Arg[%d]: %q", i, arg)
        }

        // Create source buffer
        buf := DxcBuffer{
                Ptr:      uintptr(unsafe.Pointer(&source[0])),
                Size:     uintptr(len(source)),
                Encoding: 65001, // CP_UTF8
        }
        debugPrint("DxcBuffer: Ptr=0x%X, Size=%d, Encoding=%d", buf.Ptr, buf.Size, buf.Encoding)
        debugPrint("IID_IDxcResult: %s", GUIDString(&IID_IDxcResult))

        // Get Compile method (index 3)
        compileMethod := c.getVTableMethod(3)
        debugPrint("Compile method at: 0x%X", compileMethod)

        // Call IDxcCompiler3::Compile
        // HRESULT Compile(
        //   [in]  const DxcBuffer *pSource,
        //   [in]  LPCWSTR *pArguments,
        //   [in]  UINT32 argCount,
        //   [in]  IDxcIncludeHandler *pIncludeHandler,
        //   [in]  REFIID riid,              <-- REQUIRED!
        //   [out] LPVOID *ppResult
        // );
        // Note: syscall.Syscall9(trap, nargs, a1, a2, a3, a4, a5, a6, a7, a8, a9)
        var resultPtr unsafe.Pointer
        hr, _, _ := syscall.Syscall9(
                compileMethod,
                7, // nargs - 7 arguments
                uintptr(c.ptr),                              // a1: this
                uintptr(unsafe.Pointer(&buf)),               // a2: pSource
                uintptr(unsafe.Pointer(&argPtrs[0])),        // a3: pArguments
                uintptr(len(args)),                          // a4: argCount
                uintptr(includeHandler),                     // a5: pIncludeHandler
                uintptr(unsafe.Pointer(&IID_IDxcResult)),    // a6: riid
                uintptr(unsafe.Pointer(&resultPtr)),         // a7: ppResult
                0, 0,                                        // a8, a9: unused
        )

        debugPrint("Compile: hr=0x%X, result=0x%X", hr, resultPtr)

        if hr != 0 {
                return nil, fmt.Errorf("compile failed: 0x%X", hr)
        }

        if resultPtr == nil {
                return nil, fmt.Errorf("compile returned null result")
        }

        // Parse result
        result := NewDxcResult(resultPtr)
        defer result.Close()

        compileResult := &CompileResult{
                Success: result.GetStatus() == 0,
        }

        // Get output
        if output, err := result.GetResult(); err == nil {
                compileResult.Output = output
        }

        // Get errors
        if errors, err := result.GetErrors(); err == nil {
                compileResult.Errors = errors
        }

        return compileResult, nil
}

// CompileSimple compiles HLSL and returns output
func (c *DxcCompiler) CompileSimple(source, entry, profile string) ([]byte, error) {
        result, err := c.Compile([]byte(source), []string{"-T", profile, "-E", entry})
        if err != nil {
                return nil, err
        }
        if !result.Success {
                return nil, fmt.Errorf("compilation failed:\n%s", string(result.Errors))
        }
        return result.Output, nil
}

// CompileToSpirv compiles HLSL to SPIR-V
func (c *DxcCompiler) CompileToSpirv(source, entry, profile string) ([]byte, error) {
        result, err := c.Compile([]byte(source), []string{"-T", profile, "-E", entry, "-spirv"})
        if err != nil {
                return nil, err
        }
        if !result.Success {
                return nil, fmt.Errorf("compilation failed:\n%s", string(result.Errors))
        }
        return result.Output, nil
}

// CompileRaw compiles HLSL with raw arguments (forwarded directly to DXC)
// This is the most flexible API - accepts any DXC command-line arguments
func (c *DxcCompiler) CompileRaw(source []byte, args []string) (*CompileResult, error) {
        result, err := c.Compile(source, args)
        if err != nil {
                return nil, err
        }
        return result, nil
}

// ============================================================
// BLOB READING
// ============================================================

// ReadBlob reads content from IDxcBlob
func ReadBlob(blob unsafe.Pointer) []byte {
        if blob == nil {
                return nil
        }

        vtbl := *(*uintptr)(blob)
        ptrSize := unsafe.Sizeof(uintptr(0))

        // GetBufferPointer (method 3)
        getPtr := *(*uintptr)(unsafe.Pointer(vtbl + 3*ptrSize))
        // GetBufferSize (method 4)
        getSize := *(*uintptr)(unsafe.Pointer(vtbl + 4*ptrSize))

        p, _, _ := syscall.Syscall(getPtr, 1, uintptr(blob), 0, 0)
        s, _, _ := syscall.Syscall(getSize, 1, uintptr(blob), 0, 0)

        if p == 0 || s == 0 {
                return nil
        }

        size := int(s)
        data := make([]byte, size)

        // Copy data safely
        src := unsafe.Slice((*byte)(unsafe.Pointer(p)), size)
        copy(data, src)

        return data
}
