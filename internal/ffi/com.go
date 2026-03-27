package ffi

import (
	"fmt"
	"syscall"
	"unsafe"
)

// ============================================================
// COM BLOB OPERATIONS
// ============================================================

// BlobReader is a generic interface for reading binary data from COM blobs
type BlobReader interface {
	GetBufferPointer() unsafe.Pointer
	GetBufferSize() uintptr
	GetBufferBytes() []byte
}

// ============================================================
// VTABLE HELPER
// ============================================================

// GetVTableMethod retrieves a function pointer from a COM object's virtual table
// obj: COM object pointer
// methodIndex: zero-based index in the vtable
// Returns the function pointer and error
func GetVTableMethod(obj unsafe.Pointer, methodIndex int) (uintptr, error) {
	if obj == nil {
		return 0, fmt.Errorf("object pointer is nil")
	}

	vtbl := *(*uintptr)(obj)
	if vtbl == 0 {
		return 0, fmt.Errorf("vtable pointer is nil")
	}

	methodAddr := *(*uintptr)(unsafe.Pointer(vtbl + uintptr(methodIndex)*unsafe.Sizeof(uintptr(0))))
	if methodAddr == 0 {
		return 0, fmt.Errorf("method at index %d is nil", methodIndex)
	}

	return methodAddr, nil
}

// ============================================================
// SYSCALL WRAPPERS
// ============================================================

// CallSyscall1 calls a function with 1 argument (besides the object pointer)
func CallSyscall1(funcPtr uintptr, arg1 uintptr) (uintptr, error) {
	if funcPtr == 0 {
		return 0, fmt.Errorf("function pointer is nil")
	}
	ret, _, _ := syscall.Syscall(funcPtr, 1, arg1, 0, 0)
	return ret, nil
}

// CallSyscall2 calls a function with 2 arguments
func CallSyscall2(funcPtr uintptr, arg1, arg2 uintptr) (uintptr, error) {
	if funcPtr == 0 {
		return 0, fmt.Errorf("function pointer is nil")
	}
	ret, _, _ := syscall.Syscall(funcPtr, 2, arg1, arg2, 0)
	return ret, nil
}

// CallSyscall3 calls a function with 3 arguments
func CallSyscall3(funcPtr uintptr, arg1, arg2, arg3 uintptr) (uintptr, error) {
	if funcPtr == 0 {
		return 0, fmt.Errorf("function pointer is nil")
	}
	ret, _, _ := syscall.Syscall(funcPtr, 3, arg1, arg2, arg3)
	return ret, nil
}

// CallSyscall6 calls a function with up to 6 arguments
func CallSyscall6(funcPtr uintptr, arg1, arg2, arg3, arg4, arg5, arg6 uintptr) (uintptr, error) {
	if funcPtr == 0 {
		return 0, fmt.Errorf("function pointer is nil")
	}
	ret, _, _ := syscall.Syscall6(funcPtr, 6, arg1, arg2, arg3, arg4, arg5, arg6)
	return ret, nil
}

// CallSyscall9 calls a function with up to 9 arguments
func CallSyscall9(funcPtr uintptr, args ...uintptr) (uintptr, error) {
	if funcPtr == 0 {
		return 0, fmt.Errorf("function pointer is nil")
	}
	if len(args) > 9 {
		return 0, fmt.Errorf("too many arguments (max 9, got %d)", len(args))
	}

	// Pad with zeros
	paddedArgs := make([]uintptr, 9)
	copy(paddedArgs, args)

	ret, _, _ := syscall.Syscall9(funcPtr, uintptr(len(args)),
		paddedArgs[0], paddedArgs[1], paddedArgs[2],
		paddedArgs[3], paddedArgs[4], paddedArgs[5],
		paddedArgs[6], paddedArgs[7], paddedArgs[8])
	return ret, nil
}

// ============================================================
// COM OBJECT WRAPPER
// ============================================================

// ComObject is a base type for COM object wrappers
type ComObject struct {
	Ptr  unsafe.Pointer // Pointer to the COM object
	Vtbl uintptr        // Pointer to the virtual table
}

// NewComObject creates a new COM object wrapper
func NewComObject(ptr unsafe.Pointer) *ComObject {
	if ptr == nil {
		return nil
	}
	vtbl := *(*uintptr)(ptr)
	return &ComObject{
		Ptr:  ptr,
		Vtbl: vtbl,
	}
}

// IsValid checks if the COM object is valid (non-nil)
func (co *ComObject) IsValid() bool {
	return co != nil && co.Ptr != nil && co.Vtbl != 0
}

// GetVTableMethod gets a method from this object's vtable
func (co *ComObject) GetVTableMethod(index int) (uintptr, error) {
	return GetVTableMethod(co.Ptr, index)
}

// Release releases this COM object
func (co *ComObject) Release() {
	if co != nil && co.Ptr != nil {
		Release(co.Ptr)
		co.Ptr = nil
		co.Vtbl = 0
	}
}

// QueryInterface queries for an interface on this object
func (co *ComObject) QueryInterface(iid *GUID) (unsafe.Pointer, uintptr) {
	if co == nil || co.Ptr == nil {
		return nil, E_POINTER
	}
	return QueryInterface(co.Ptr, iid)
}

// AddRef increments reference count
func (co *ComObject) AddRef() uint32 {
	if co == nil || co.Ptr == nil {
		return 0
	}
	return AddRef(co.Ptr)
}

// ============================================================
// SAFE COM WRAPPER
// ============================================================

// SafeComObject is a COM object wrapper with automatic cleanup using defer
// Use with 'defer obj.Close()' to ensure Release is called
type SafeComObject struct {
	*ComObject
}

// NewSafeComObject wraps a COM object for automatic cleanup
func NewSafeComObject(ptr unsafe.Pointer) *SafeComObject {
	return &SafeComObject{NewComObject(ptr)}
}

// Close releases the COM object
func (sco *SafeComObject) Close() error {
	if sco != nil && sco.ComObject != nil {
		sco.Release()
	}
	return nil
}
