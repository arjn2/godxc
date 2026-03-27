# DXGOC - Go-Based DirectX Shader Compiler

A high-performance Go wrapper for the DirectX Shader Compiler (DXC) using pure FFI (Foreign Function Interface) with direct DLL calls. **No CGO, no C compiler required.**

## Overview

DXGOC is a Go application that interfaces directly with `dxcompiler.dll` via pure Go FFI syscalls. It provides the same functionality as the DXC command-line compiler but with:

- ✅ **Pure Go FFI** - Direct Windows DLL calls via syscalls (no CGO)
- ✅ **High Performance** - 20-25ms shader compilation
- ✅ **Zero C Dependencies** - No external C libraries or compilers
- ✅ **COM Interface** - Direct IDxcCompiler3 integration
- ✅ **DXC-Compatible** - Same command-line syntax as DXC

## Quick Start

### Prerequisites

You need `dxcompiler.dll` from:
- **Vulkan SDK:** https://vulkan.lunarg.com/ (includes DXC)
- **Windows SDK:** https://developer.microsoft.com/en-us/windows/downloads/windows-sdk/
- **DXC Releases:** https://github.com/microsoft/DirectXShaderCompiler/releases

Place `dxcompiler.dll` in your PATH or in the same directory as `dxgoc.exe`.

### Build

**Main Compiler:**
```bash
go build -o bin/dxgoc.exe ./cmd/dxgoc
```

**Both Binaries (dxgoc + hotreload demo):**
```bash
.\build.bat
.\build-hotreload.bat
```

### Compile Shaders

```bash
# Vertex shader (DXC-compatible syntax)
.\dxgoc.exe -i shaders/vertex.hlsl -T vs_6_0 -E main -o shader.dxil

# Pixel shader
.\dxgoc.exe -i shaders/pixel.hlsl -T ps_6_0 -E main -o shader.dxil

# Compute shader
.\dxgoc.exe -i shaders/compute.hlsl -T cs_6_0 -E main -o shader.dxil

# With defines
.\dxgoc.exe -i shader.hlsl -T vs_6_0 -E main -D DEBUG=1 -D QUALITY=2

# With includes
.\dxgoc.exe -i shader.hlsl -T ps_6_0 -I ./include -I ./common
```

### Test Compilation

Sample shaders are included in `shaders/` directory:

```bash
# Compile basic shaders
.\dxgoc.exe -i shaders/basic_vertex.hlsl -T vs_6_0 -E main
.\dxgoc.exe -i shaders/basic_pixel.hlsl -T ps_6_0 -E main
.\dxgoc.exe -i shaders/basic_compute.hlsl -T cs_6_0 -E main

# Compile advanced shaders
.\dxgoc.exe -i shaders/advanced_vertex.hlsl -T vs_6_0 -E main
.\dxgoc.exe -i shaders/advanced_pixel.hlsl -T ps_6_0 -E main
.\dxgoc.exe -i shaders/advanced_compute.hlsl -T cs_6_0 -E main
```

## Architecture

### Core Components

**FFI Library** (`internal/ffi/`)
- Generic Foreign Function Interface for Windows
- Direct syscall wrapper (no CGO marshaling)
- COM object lifecycle management (AddRef, Release, QueryInterface)
- GUID and HRESULT utilities
- ~600 lines of pure Go

**DXC Bindings** (`internal/dxc/`)
- COM wrapper for DirectX Shader Compiler
- Direct `dxcompiler.dll` interface
- Error handling and diagnostics
- All shader profiles supported
- ~400 lines of Go

## How It Works

```
Your Go Application
        ↓
DXGOC (Go Code)
        ↓
FFI Layer (internal/ffi/)
        ↓
Windows Syscalls
        ↓
dxcompiler.dll (DirectX Shader Compiler)
        ↓
DXIL Bytecode
```

DXGOC calls `dxcompiler.dll` **directly via pure Go FFI**, not through CGO or subprocess.

## Hot Reload Demo (Experimental)

`hotreload-demo.exe` provides real-time shader compilation with file watching:

```bash
# Watch a shader file and recompile on save
.\bin\hotreload-demo.exe -shader shaders/basic_vertex.hlsl -watch

# Watch entire directory
.\bin\hotreload-demo.exe -dir shaders/ -watch
```

**Note:** Hot reload functionality is experimental and for demonstration purposes.

## Usage from Go Code

```go
package main

import (
    "fmt"
    "os"
    "dxgoc/internal/dxc"
)

func main() {
    // Create compiler
    compiler, err := dxc.NewCompiler("")
    if err != nil {
        panic(err)
    }
    defer compiler.Close()

    // Read shader source
    source, err := os.ReadFile("shader.hlsl")
    if err != nil {
        panic(err)
    }

    // Compile
    result, err := compiler.Compile(source, []string{
        "-T", "vs_6_0",
        "-E", "main",
    })
    if err != nil {
        panic(err)
    }

    if result.Success {
        fmt.Printf("✓ Compiled: %d bytes\n", len(result.Output))
        os.WriteFile("shader.dxil", result.Output, 0644)
    } else {
        fmt.Printf("✗ Error: %s\n", result.Errors)
    }
}
```

## Performance

Compilation speed: **20-25ms per shader** (average across all profiles)

See `DXC_PERFORMANCE_BENCHMARK.md` for detailed benchmarks and analysis.

## Project Structure

```
.
├── build.bat                    # Build main compiler
├── build-hotreload.bat          # Build hot reload demo
├── LICENSE                      # MIT License
├── NOTICES                      # Attribution
├── README.md                    # This file
├── cmd/
│   ├── dxgoc/                  # Main compiler CLI
│   └── hotreload-demo/         # Hot reload demo (experimental)
├── internal/
│   ├── ffi/                    # Pure Go FFI layer
│   └── dxc/                    # DXC bindings
├── hotreload/                  # Experimental hot reload library
├── shaders/                    # Sample shader files
│   ├── basic_vertex.hlsl
│   ├── basic_pixel.hlsl
│   ├── basic_compute.hlsl
│   ├── advanced_vertex.hlsl
│   ├── advanced_pixel.hlsl
│   └── advanced_compute.hlsl
└── bin/
    ├── dxgoc.exe              # Main compiler binary
    └── hotreload-demo.exe     # Hot reload demo binary
```

## Building

### Quick Build

```bash
# Build with batch script (handles environment setup)
.\build.bat

# Build hot reload demo
.\build-hotreload.bat
```

### Manual Build

```bash
# Main compiler
go build -o bin/dxgoc.exe ./cmd/dxgoc

# Hot reload demo
go build -o bin/hotreload-demo.exe ./cmd/hotreload-demo
```

### Optimized Build

```bash
# Strip debug info and reduce binary size
go build -ldflags="-s -w" -trimpath -buildvcs=false -o bin/dxgoc.exe ./cmd/dxgoc
```

## DLL Path Detection

DXGOC automatically detects `dxcompiler.dll` from multiple locations:

1. **Same directory as executable** (highest priority)
   ```bash
   # Place dxcompiler.dll next to dxgoc.exe
   bin/dxgoc.exe
   bin/dxcompiler.dll
   ```

2. **PATH environment variable**
   ```bash
   # DXC's bin directory should be in your PATH
   # Check with: where dxc
   where dxc
   ```

3. **Common installation paths**
   ```
   C:\Program Files\Microsoft\DirectXShaderCompiler\bin\
   C:\Program Files (x86)\Microsoft Visual Studio\2022\*\DXC\
   ```

4. **Custom path via flag**
   ```bash
   # Explicitly specify DXC location
   .\dxgoc.exe -dll "C:\path\to\dxcompiler.dll" -i shader.hlsl -T vs_6_0
   ```

### Setup DXC Path

**Option A: Add to PATH (Recommended)**
```powershell
# Vulkan SDK (usually already in PATH)
$env:PATH += ";C:\VulkanSDK\1.x.x\Bin"

# Or manually via System Properties
setx PATH "%PATH%;C:\path\to\dxc\bin"
```

**Option B: Copy to project**
```bash
# Copy dxcompiler.dll to bin directory
copy "C:\VulkanSDK\1.x.x\Bin\dxcompiler.dll" .\bin\
copy "C:\VulkanSDK\1.x.x\Bin\dxil.dll" .\bin\
```

**Option C: Environment variable (Optional)**
```powershell
# Set DXC_PATH environment variable
$env:DXC_PATH = "C:\VulkanSDK\1.x.x\Bin"
```

## Requirements

- Windows 10+ (x64, x86, or ARM64)
- Go 1.18 or later
- `dxcompiler.dll` (from Vulkan SDK, Windows SDK, or DXC release)
- Optional: `dxil.dll` (DXIL validator)

## License

**DXGOC Code:** MIT License (see `LICENSE`)

**DirectX Shader Compiler:** Multiple licenses (see `NOTICES`)
- LLVM License
- MIT License
- Microsoft License

## Supported Shader Profiles

- `vs_6_0` - Vertex Shader
- `ps_6_0` - Pixel Shader
- `cs_6_0` - Compute Shader
- `gs_6_0` - Geometry Shader
- `hs_6_0` - Hull Shader
- `ds_6_0` - Domain Shader

## Key Differences from CGO Approach

| Aspect | DXGOC (Pure FFI) | CGO |
|--------|------------------|-----|
| **Compilation** | No C compiler needed | Requires C compiler |
| **Speed** | 3-5x faster | Slower (thread marshaling) |
| **DLL Loading** | Dynamic (syscalls) | Static linking |
| **Deployment** | Single .exe + DLL | Complex setup |
| **Platform** | Windows only | Cross-platform |

## Troubleshooting

**"dxcompiler.dll not found"**
- See **DLL Path Detection** section above
- Install Vulkan SDK or Windows SDK
- Add DXC directory to PATH
- Place `dxcompiler.dll` in same directory as `dxgoc.exe`
- Or use `-dll "path\to\dxcompiler.dll"` flag

**Compilation fails**
- Check HLSL syntax
- Verify shader profile is supported
- Check include paths are correct
- Review error messages from compiler

## Contributing

Improvements and bug reports welcome!

---

**DXGOC** - Production-grade Go shader compilation.
