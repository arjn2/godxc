package ffi

import (
	"fmt"
	"sync"
	"syscall"
)

// ============================================================
// DLL MANAGER
// ============================================================

// DllManager manages loading and caching of DLLs
type DllManager struct {
	mu    sync.RWMutex
	dlls  map[string]*syscall.LazyDLL
	procs map[string]map[string]*syscall.LazyProc
}

// NewDllManager creates a new DLL manager
func NewDllManager() *DllManager {
	return &DllManager{
		dlls:  make(map[string]*syscall.LazyDLL),
		procs: make(map[string]map[string]*syscall.LazyProc),
	}
}

// GetDll gets or loads a DLL by name
func (dm *DllManager) GetDll(name string) (*syscall.LazyDLL, error) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dll, ok := dm.dlls[name]; ok {
		return dll, nil
	}

	dll := syscall.NewLazyDLL(name)
	dm.dlls[name] = dll
	dm.procs[name] = make(map[string]*syscall.LazyProc)

	return dll, nil
}

// GetProc gets a procedure from a DLL
func (dm *DllManager) GetProc(dllName string, procName string) (*syscall.LazyProc, error) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dll, ok := dm.dlls[dllName]
	if !ok {
		return nil, fmt.Errorf("DLL not loaded: %s", dllName)
	}

	procMap, ok := dm.procs[dllName]
	if !ok {
		return nil, fmt.Errorf("procedure map not found for DLL: %s", dllName)
	}

	if proc, ok := procMap[procName]; ok {
		return proc, nil
	}

	proc := dll.NewProc(procName)
	procMap[procName] = proc

	return proc, nil
}

// GetOrLoadProc gets or loads a procedure from a DLL, loading the DLL if necessary
func (dm *DllManager) GetOrLoadProc(dllName string, procName string) (*syscall.LazyProc, error) {
	if _, err := dm.GetDll(dllName); err != nil {
		return nil, err
	}

	return dm.GetProc(dllName, procName)
}

// Clear clears all cached DLLs
func (dm *DllManager) Clear() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.dlls = make(map[string]*syscall.LazyDLL)
	dm.procs = make(map[string]map[string]*syscall.LazyProc)
}

// ============================================================
// GLOBAL DLL MANAGER
// ============================================================

var globalDllManager = NewDllManager()

// GetGlobalDllManager returns the global DLL manager instance
func GetGlobalDllManager() *DllManager {
	return globalDllManager
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

// LoadDllProc loads a DLL and gets a specific procedure
func LoadDllProc(dllName string, procName string) (*syscall.LazyProc, error) {
	return globalDllManager.GetOrLoadProc(dllName, procName)
}

// CallDllProc calls a procedure from a DLL with up to 6 arguments
func CallDllProc(dllName string, procName string, args ...uintptr) (uintptr, error) {
	proc, err := LoadDllProc(dllName, procName)
	if err != nil {
		return 0, err
	}

	if len(args) == 0 {
		ret, _, _ := proc.Call()
		return ret, nil
	}

	// Use appropriate syscall based on argument count
	switch len(args) {
	case 1:
		ret, _, _ := syscall.Syscall(proc.Addr(), 1, args[0], 0, 0)
		return ret, nil
	case 2:
		ret, _, _ := syscall.Syscall(proc.Addr(), 2, args[0], args[1], 0)
		return ret, nil
	case 3:
		ret, _, _ := syscall.Syscall(proc.Addr(), 3, args[0], args[1], args[2])
		return ret, nil
	case 4:
		ret, _, _ := syscall.Syscall6(proc.Addr(), 4, args[0], args[1], args[2], args[3], 0, 0)
		return ret, nil
	case 5:
		ret, _, _ := syscall.Syscall6(proc.Addr(), 5, args[0], args[1], args[2], args[3], args[4], 0)
		return ret, nil
	case 6:
		ret, _, _ := syscall.Syscall6(proc.Addr(), 6, args[0], args[1], args[2], args[3], args[4], args[5])
		return ret, nil
	default:
		return 0, fmt.Errorf("too many arguments (max 6, got %d)", len(args))
	}
}
