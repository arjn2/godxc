// Package dxc - DirectX Shader Compiler Go bindings
// Wraps internal ffi library for COM operations
package dxc

import (
	"fmt"
	"unsafe"

	"dxgoc/internal/ffi"
)

// ============================================================
// DEBUG
// ============================================================

var DebugMode = false

func debugPrint(format string, args ...interface{}) {
	if DebugMode {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func SetDebug(enabled bool) { DebugMode = enabled }

// ============================================================
// COM INITIALIZATION (delegated to go-ffi)
// ============================================================

// InitCOM initializes COM library (required for DxcCreateInstance)
func InitCOM() error {
	err := ffi.InitCOM()
	if err != nil {
		debugPrint("COM initialization failed: %v", err)
	}
	return err
}

// ============================================================
// TYPE ALIASES (for backward compatibility)
// ============================================================

// GUID is a re-export from go-ffi
type GUID = ffi.GUID

// ============================================================
// WRAPPER FUNCTIONS (for backward compatibility with existing code)
// ============================================================

// GUIDString returns the string representation of a GUID
func GUIDString(g *GUID) string {
	return ffi.GUIDString(g)
}

// QueryInterface calls QueryInterface on a COM object
func QueryInterface(obj unsafe.Pointer, iid *GUID) (unsafe.Pointer, uintptr) {
	return ffi.QueryInterface(obj, iid)
}

// Release calls Release on a COM object
func Release(obj unsafe.Pointer) {
	ffi.Release(obj)
}

// AddRef calls AddRef on a COM object
func AddRef(obj unsafe.Pointer) uint32 {
	return ffi.AddRef(obj)
}

// TryQueryInterface tries to get an interface from an object, returns true if successful
func TryQueryInterface(obj unsafe.Pointer, iid *GUID) (unsafe.Pointer, bool) {
	return ffi.TryQueryInterface(obj, iid)
}

// ============================================================
// HRESULT HELPERS (re-exported from go-ffi for convenience)
// ============================================================

// HRESULT error codes
const (
	S_OK          = ffi.S_OK
	S_FALSE       = ffi.S_FALSE
	E_NOTIMPL     = ffi.E_NOTIMPL
	E_NOINTERFACE = ffi.E_NOINTERFACE
	E_POINTER     = ffi.E_POINTER
	E_ABORT       = ffi.E_ABORT
	E_FAIL        = ffi.E_FAIL
	E_INVALIDARG  = ffi.E_INVALIDARG
	E_OUTOFMEMORY = ffi.E_OUTOFMEMORY
)

// HResultToString converts HRESULT to human-readable string
func HResultToString(hr uintptr) string {
	return ffi.HResultToString(hr)
}

// IsSuccess checks if an HRESULT represents success
func IsSuccess(hr uintptr) bool {
	return ffi.IsSuccess(hr)
}
