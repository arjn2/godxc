// Package ffi - Generic Foreign Function Interface for Windows COM
// Provides low-level syscall bindings and COM object manipulation
package ffi

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

// ============================================================
// TYPES
// ============================================================

// GUID represents a Windows GUID/UUID (16 bytes)
type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// ============================================================
// COM INITIALIZATION
// ============================================================

var (
	ole32              = syscall.NewLazyDLL("ole32.dll")
	coInitializeEx     = ole32.NewProc("CoInitializeEx")
	coUninitialize     = ole32.NewProc("CoUninitialize")
	coCreateInstance   = ole32.NewProc("CoCreateInstance")
	coGetClassObject   = ole32.NewProc("CoGetClassObject")

	comInitOnce sync.Once
	comInitErr  error
)

const (
	COINIT_MULTITHREADED = 0
	COINIT_APARTMENTTHREADED = 2
)

// InitCOM initializes the COM library for the current thread
// Must be called before using any COM objects
// Safe to call multiple times - only initializes once per thread
func InitCOM() error {
	comInitOnce.Do(func() {
		ret, _, _ := coInitializeEx.Call(0, COINIT_MULTITHREADED)
		if ret != S_OK && ret != S_FALSE {
			comInitErr = fmt.Errorf("CoInitializeEx failed: 0x%X", ret)
		}
	})
	return comInitErr
}

// UninitCOM uninitializes COM for the current thread
func UninitCOM() {
	coUninitialize.Call()
}

// ============================================================
// HRESULT CONSTANTS
// ============================================================

const (
	S_OK          uintptr = 0x00000000 // Success
	S_FALSE       uintptr = 0x00000001 // Success but false
	E_NOTIMPL     uintptr = 0x80004001 // Not implemented
	E_NOINTERFACE uintptr = 0x80004002 // Interface not supported
	E_POINTER     uintptr = 0x80004003 // Invalid pointer
	E_ABORT       uintptr = 0x80004004 // Operation aborted
	E_FAIL        uintptr = 0x80004005 // Generic failure
	E_INVALIDARG  uintptr = 0x80070057 // Invalid argument
	E_OUTOFMEMORY uintptr = 0x8007000E // Out of memory
)

// ============================================================
// HRESULT HELPERS
// ============================================================

// IsSuccess checks if an HRESULT indicates success (S_OK or S_FALSE)
func IsSuccess(hr uintptr) bool {
	return hr == S_OK || hr == S_FALSE
}

// IsFailed checks if an HRESULT indicates failure
func IsFailed(hr uintptr) bool {
	return !IsSuccess(hr)
}

// HResultToString converts an HRESULT code to a human-readable string
func HResultToString(hr uintptr) string {
	switch hr {
	case S_OK:
		return "S_OK"
	case S_FALSE:
		return "S_FALSE"
	case E_NOTIMPL:
		return "E_NOTIMPL"
	case E_NOINTERFACE:
		return "E_NOINTERFACE"
	case E_POINTER:
		return "E_POINTER"
	case E_ABORT:
		return "E_ABORT"
	case E_FAIL:
		return "E_FAIL"
	case E_INVALIDARG:
		return "E_INVALIDARG"
	case E_OUTOFMEMORY:
		return "E_OUTOFMEMORY"
	default:
		return fmt.Sprintf("HRESULT(0x%X)", hr)
	}
}

// ============================================================
// COM OBJECT OPERATIONS
// ============================================================

// QueryInterface queries a COM object for a specific interface
// Returns the interface pointer and HRESULT code
func QueryInterface(obj unsafe.Pointer, iid *GUID) (unsafe.Pointer, uintptr) {
	if obj == nil {
		return nil, E_POINTER
	}

	// Get vtable pointer from object
	vtbl := *(*uintptr)(obj)

	// Get QueryInterface function pointer (method 0)
	queryAddr := *(*uintptr)(unsafe.Pointer(vtbl))

	var ptr unsafe.Pointer
	ret, _, _ := syscall.Syscall(queryAddr, 3,
		uintptr(obj),
		uintptr(unsafe.Pointer(iid)),
		uintptr(unsafe.Pointer(&ptr)))

	return ptr, ret
}

// TryQueryInterface attempts to query an interface, returns (ptr, success)
func TryQueryInterface(obj unsafe.Pointer, iid *GUID) (unsafe.Pointer, bool) {
	ptr, hr := QueryInterface(obj, iid)
	return ptr, IsSuccess(hr)
}

// AddRef increments the reference count of a COM object
// Returns the new reference count
func AddRef(obj unsafe.Pointer) uint32 {
	if obj == nil {
		return 0
	}

	vtbl := *(*uintptr)(obj)
	// AddRef is method 1
	addrefAddr := *(*uintptr)(unsafe.Pointer(vtbl + 1*unsafe.Sizeof(uintptr(0))))

	ret, _, _ := syscall.Syscall(addrefAddr, 1, uintptr(obj), 0, 0)
	return uint32(ret)
}

// Release decrements the reference count of a COM object
// When reference count reaches 0, the object is deleted
func Release(obj unsafe.Pointer) {
	if obj == nil {
		return
	}

	vtbl := *(*uintptr)(obj)
	// Release is method 2
	releaseAddr := *(*uintptr)(unsafe.Pointer(vtbl + 2*unsafe.Sizeof(uintptr(0))))

	syscall.Syscall(releaseAddr, 1, uintptr(obj), 0, 0)
}

// ============================================================
// GUID UTILITIES
// ============================================================

// GUIDString returns the string representation of a GUID in canonical format
// Format: {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}
func GUIDString(g *GUID) string {
	if g == nil {
		return "{00000000-0000-0000-0000-000000000000}"
	}

	return fmt.Sprintf("{%08X-%04X-%04X-%02X%02X-%02X%02X%02X%02X%02X%02X}",
		g.Data1,
		g.Data2,
		g.Data3,
		g.Data4[0], g.Data4[1],
		g.Data4[2], g.Data4[3], g.Data4[4], g.Data4[5], g.Data4[6], g.Data4[7])
}

// GUIDEqual compares two GUIDs for equality
func GUIDEqual(g1, g2 *GUID) bool {
	if g1 == nil || g2 == nil {
		return g1 == g2
	}
	return g1.Data1 == g2.Data1 &&
		g1.Data2 == g2.Data2 &&
		g1.Data3 == g2.Data3 &&
		g1.Data4 == g2.Data4
}

// ============================================================
// STRING CONVERSION
// ============================================================

// UTF16BytesToString converts UTF-16 bytes (pointer and length) to Go string
func UTF16BytesToString(ptr unsafe.Pointer, length int) string {
	if ptr == nil || length == 0 {
		return ""
	}

	// Cast to UTF-16 array
	utf16Slice := (*[1 << 30]uint16)(ptr)[:length:length]
	return syscall.UTF16ToString(utf16Slice)
}

// StringToUTF16Ptr converts a Go string to a UTF-16 pointer
// The returned pointer points to dynamically allocated memory that must not be freed by Go
func StringToUTF16Ptr(s string) (unsafe.Pointer, error) {
	utf16Ptr, err := syscall.UTF16PtrFromString(s)
	return unsafe.Pointer(utf16Ptr), err
}

// ============================================================
// DLL LOADING
// ============================================================

// LoadDLL loads a dynamic library by name
// Uses lazy loading - DLL is loaded on first use
func LoadDLL(name string) (*syscall.LazyDLL, error) {
	dll := syscall.NewLazyDLL(name)
	// Lazy load, will fail when trying to use if DLL not found
	return dll, nil
}

// LoadProc gets a procedure from a loaded DLL
func LoadProc(dll *syscall.LazyDLL, name string) *syscall.LazyProc {
	return dll.NewProc(name)
}
