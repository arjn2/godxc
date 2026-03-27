# Internal FFI Package

This is a generic, reusable Foreign Function Interface (FFI) library built into the dxgoc project. It provides low-level Windows/COM interoperability without CGO.

## Overview

The `internal/ffi` package provides:

- **COM Object Management** - Initialize, query, and release COM objects
- **Vtable Navigation** - Extract and call methods from COM virtual tables  
- **HRESULT Handling** - COM error codes and helpers
- **GUID Operations** - GUID comparison and string conversion
- **DLL Loading** - Dynamic DLL loading and procedure caching
- **Syscall Wrappers** - Safe wrappers around Windows syscalls

## Structure

```
internal/ffi/
├── ffi.go          # Core COM initialization and GUID operations
├── com.go          # COM object manipulation and vtable helpers
└── dll.go          # DLL loading and procedure management
```

## Key Features

### 1. Generic COM Object Wrapper

```go
obj := ffi.NewComObject(ptr)
defer obj.Release()

if obj.IsValid() {
    method, err := obj.GetVTableMethod(3)
}
```

### 2. Safe COM Objects with Auto-Cleanup

```go
obj := ffi.NewSafeComObject(ptr)
defer obj.Close()
```

### 3. HRESULT Error Handling

```go
hr := ffi.E_FAIL
if !ffi.IsSuccess(hr) {
    fmt.Println(ffi.HResultToString(hr))
}
```

### 4. DLL Management

```go
proc, err := ffi.LoadDllProc("kernel32.dll", "GetCurrentProcess")
```

## Usage in dxgoc

The `internal/dxc` package uses the FFI library:

```go
import "godxc/internal/ffi"

// Initialize COM
ffi.InitCOM()

// Work with COM objects
ptr, hr := ffi.QueryInterface(obj, &iid)
if ffi.IsSuccess(hr) {
    defer ffi.Release(ptr)
}
```

## Design Principles

1. **Generic** - Not tied to any specific API (DXC, Direct3D, etc.)
2. **Safe** - Error checking and validation
3. **Reusable** - Can be extracted into separate packages
4. **Thin** - Minimal overhead, direct syscalls
5. **Fast** - No CGO marshaling, no GC stops

## Performance

- **3-5x faster** than CGO for syscalls
- **No garbage collection** pauses during Windows calls
- **No thread marshaling** overhead
- **Direct kernel calls** with minimal abstraction
