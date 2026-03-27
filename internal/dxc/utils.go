// Package dxc - DirectX Shader Compiler Go bindings
// IDxcUtils interface handling
package dxc

import (
	"fmt"
	"syscall"
	"unsafe"
)

// ============================================================
// IDXCUTILS INTERFACE
// ============================================================

// DxcUtils wraps IDxcUtils interface
type DxcUtils struct {
	ptr  unsafe.Pointer
	vtbl uintptr
}

// NewDxcUtils creates IDxcUtils from DxcCreateInstance
func NewDxcUtils(proc *syscall.Proc) (*DxcUtils, error) {
	debugPrint("Creating IDxcUtils...")

	for _, iid := range UtilsIIDs {
		debugPrint("  Trying %s: %s", iid.Name, GUIDString(&iid.GUID))

		var ptr unsafe.Pointer
		hr, _, _ := proc.Call(
			uintptr(unsafe.Pointer(&CLSID_DxcLibrary)),
			uintptr(unsafe.Pointer(&iid.GUID)),
			uintptr(unsafe.Pointer(&ptr)),
		)

		debugPrint("    hr=0x%X, ptr=0x%X", hr, ptr)

		if hr == 0 && ptr != nil {
			utils := &DxcUtils{
				ptr:  ptr,
				vtbl: *(*uintptr)(ptr),
			}
			debugPrint("    -> SUCCESS: %s", iid.Name)
			return utils, nil
		}
	}

	return nil, fmt.Errorf("IDxcUtils creation failed - no matching IID")
}

// Close releases IDxcUtils
func (u *DxcUtils) Close() {
	if u.ptr != nil {
		Release(u.ptr)
		u.ptr = nil
	}
}

// GetPtr returns the underlying pointer
func (u *DxcUtils) GetPtr() unsafe.Pointer { return u.ptr }

// getVTableMethod returns the method pointer at the given vtable index
func (u *DxcUtils) getVTableMethod(index int) uintptr {
	if u.vtbl == 0 {
		return 0
	}
	return *(*uintptr)(unsafe.Pointer(u.vtbl + uintptr(index)*unsafe.Sizeof(uintptr(0))))
}

// QueryInterface queries for another interface from this object
func (u *DxcUtils) QueryInterface(iid *GUID) (unsafe.Pointer, uintptr) {
	return QueryInterface(u.ptr, iid)
}

// PrintVTable prints vtable entries for debugging
func (u *DxcUtils) PrintVTable() {
	if u.vtbl == 0 {
		debugPrint("  VTable: nil")
		return
	}
	debugPrint("  VTable at: 0x%X", u.vtbl)
	for i := 0; i < 15; i++ {
		addr := u.getVTableMethod(i)
		debugPrint("    [%2d] 0x%X", i, addr)
	}
}

// CreateBlobFromFile creates a blob from file content
func (u *DxcUtils) CreateBlobFromFile(filePath string) (unsafe.Pointer, error) {
	path, _ := syscall.UTF16PtrFromString(filePath)

	method := u.getVTableMethod(4) // CreateBlobFromFile is index 4
	var blob unsafe.Pointer
	hr, _, _ := syscall.Syscall6(
		method,
		4,
		uintptr(u.ptr),
		uintptr(unsafe.Pointer(path)),
		0, // codePage (0 = auto-detect)
		uintptr(unsafe.Pointer(&blob)),
		0, 0,
	)

	if hr != 0 {
		return nil, fmt.Errorf("CreateBlobFromFile failed: 0x%X", hr)
	}
	return blob, nil
}

// CreateBlobWithEncoding creates a blob from memory
func (u *DxcUtils) CreateBlobWithEncoding(data []byte, encoding uint32) (unsafe.Pointer, error) {
	method := u.getVTableMethod(5) // CreateBlobWithEncodingOnHeapCopy is index 5
	var blob unsafe.Pointer
	hr, _, _ := syscall.Syscall6(
		method,
		5,
		uintptr(u.ptr),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(encoding),
		uintptr(unsafe.Pointer(&blob)),
		0,
	)

	if hr != 0 {
		return nil, fmt.Errorf("CreateBlobWithEncoding failed: 0x%X", hr)
	}
	return blob, nil
}

// CreateDefaultIncludeHandler creates default include handler
// Note: This is currently disabled as it causes crashes with some DXC versions
// The -I flags should work on their own through DXC's command-line processing
func (u *DxcUtils) CreateDefaultIncludeHandler() (unsafe.Pointer, error) {
	// Disabled for now - returning error to skip include handler usage
	return nil, fmt.Errorf("CreateDefaultIncludeHandler not implemented")
}
