// Package dxc - DirectX Shader Compiler Go bindings
// GUID constants from official dxcapi.h
//
// Source: Microsoft DirectX Shader Compiler (DXC) official header
// All GUIDs extracted from dxcapi.h for accuracy
package dxc

// Note: GUID type is imported from ffi package via com.go type alias

// ============================================================
// CLASS IDs (CLSIDs) - For DxcCreateInstance
// ============================================================

// CLSID_DxcCompiler - Use to create IDxcCompiler, IDxcCompiler2, IDxcCompiler3
// {73E22D93-E6CE-47F3-B5BF-F0664F39C1B0}
var CLSID_DxcCompiler = GUID{
        0x73e22d93,
        0xe6ce,
        0x47f3,
        [8]byte{0xb5, 0xbf, 0xf0, 0x66, 0x4f, 0x39, 0xc1, 0xb0},
}

// CLSID_DxcLibrary - Use to create IDxcLibrary, IDxcUtils
// {6245D6AF-66E0-48FD-80B4-4D271796748C}
var CLSID_DxcLibrary = GUID{
        0x6245d6af,
        0x66e0,
        0x48fd,
        [8]byte{0x80, 0xb4, 0x4d, 0x27, 0x17, 0x96, 0x74, 0x8c},
}

// CLSID_DxcUtils - Same as CLSID_DxcLibrary (alias)
// {6245D6AF-66E0-48FD-80B4-4D271796748C}
var CLSID_DxcUtils = CLSID_DxcLibrary

// CLSID_DxcLinker - Use to create IDxcLinker
// {EF6A8087-B0EA-4D56-9E45-D07E1A8B7806}
var CLSID_DxcLinker = GUID{
        0xef6a8087,
        0xb0ea,
        0x4d56,
        [8]byte{0x9e, 0x45, 0xd0, 0x7e, 0x1a, 0x8b, 0x78, 0x06},
}

// CLSID_DxcValidator - Use to create IDxcValidator, IDxcValidator2
// {8CA3E215-F728-4CF3-8CDD-88AF917587A1}
var CLSID_DxcValidator = GUID{
        0x8ca3e215,
        0xf728,
        0x4cf3,
        [8]byte{0x8c, 0xdd, 0x88, 0xaf, 0x91, 0x75, 0x87, 0xa1},
}

// CLSID_DxcAssembler - Use to create IDxcAssembler
// {D728DB68-F903-4F80-94CD-DCCF76EC7151}
var CLSID_DxcAssembler = GUID{
        0xd728db68,
        0xf903,
        0x4f80,
        [8]byte{0x94, 0xcd, 0xdc, 0xcf, 0x76, 0xec, 0x71, 0x51},
}

// CLSID_DxcContainerReflection - Use to create IDxcContainerReflection
// {B9F54489-55B8-400C-BA3A-1675E4728B91}
var CLSID_DxcContainerReflection = GUID{
        0xb9f54489,
        0x55b8,
        0x400c,
        [8]byte{0xba, 0x3a, 0x16, 0x75, 0xe4, 0x72, 0x8b, 0x91},
}

// CLSID_DxcOptimizer - Use to create IDxcOptimizer
// {AE2CD79F-CC22-453F-9B6B-B124E7A5204C}
var CLSID_DxcOptimizer = GUID{
        0xae2cd79f,
        0xcc22,
        0x453f,
        [8]byte{0x9b, 0x6b, 0xb1, 0x24, 0xe7, 0xa5, 0x20, 0x4c},
}

// CLSID_DxcContainerBuilder - Use to create IDxcContainerBuilder
// {94134294-411F-4574-B4D0-8741E25240D2}
var CLSID_DxcContainerBuilder = GUID{
        0x94134294,
        0x411f,
        0x4574,
        [8]byte{0xb4, 0xd0, 0x87, 0x41, 0xe2, 0x52, 0x40, 0xd2},
}

// CLSID_DxcPdbUtils - Use to create IDxcPdbUtils, IDxcPdbUtils2
// {54621DFB-F2CE-457E-AE8C-EC355FAEEC7C}
var CLSID_DxcPdbUtils = GUID{
        0x54621dfb,
        0xf2ce,
        0x457e,
        [8]byte{0xae, 0x8c, 0xec, 0x35, 0x5f, 0xae, 0xec, 0x7c},
}

// CLSID_DxcCompilerArgs - Use to create IDxcCompilerArgs
// {3E56AE82-224D-470F-A1A1-FE3016EE9F9D}
var CLSID_DxcCompilerArgs = GUID{
        0x3e56ae82,
        0x224d,
        0x470f,
        [8]byte{0xa1, 0xa1, 0xfe, 0x30, 0x16, 0xee, 0x9f, 0x9d},
}

// ============================================================
// INTERFACE IDs (IIDs) - Blob Interfaces
// ============================================================

// IID_IDxcBlob - Basic blob interface (alias of ID3D10Blob, ID3DBlob)
// {8BA5FB08-5195-40E2-AC58-0D989C3A0102}
var IID_IDxcBlob = GUID{
        0x8ba5fb08,
        0x5195,
        0x40e2,
        [8]byte{0xac, 0x58, 0x0d, 0x98, 0x9c, 0x3a, 0x01, 0x02},
}

// IID_IDxcBlobEncoding - Blob with known encoding
// {7241D424-2646-4191-97C0-98E96E42FC68}
var IID_IDxcBlobEncoding = GUID{
        0x7241d424,
        0x2646,
        0x4191,
        [8]byte{0x97, 0xc0, 0x98, 0xe9, 0x6e, 0x42, 0xfc, 0x68},
}

// IID_IDxcBlobWide - Blob containing null-terminated wide string (UTF-16 on Windows)
// {A3F84EAB-0FAA-497E-A39C-EE6ED60B2D84}
var IID_IDxcBlobWide = GUID{
        0xa3f84eab,
        0x0faa,
        0x497e,
        [8]byte{0xa3, 0x9c, 0xee, 0x6e, 0xd6, 0x0b, 0x2d, 0x84},
}

// IID_IDxcBlobUtf8 - Blob containing UTF-8 encoded string
// {3DA636C9-BA71-4024-A301-30CBF125305B}
var IID_IDxcBlobUtf8 = GUID{
        0x3da636c9,
        0xba71,
        0x4024,
        [8]byte{0xa3, 0x01, 0x30, 0xcb, 0xf1, 0x25, 0x30, 0x5b},
}

// ============================================================
// INTERFACE IDs (IIDs) - Handler Interfaces
// ============================================================

// IID_IDxcIncludeHandler - Include file handler interface
// {7F61FC7D-950D-467F-B3E3-3C02FB49187C}
var IID_IDxcIncludeHandler = GUID{
        0x7f61fc7d,
        0x950d,
        0x467f,
        [8]byte{0xb3, 0xe3, 0x3c, 0x02, 0xfb, 0x49, 0x18, 0x7c},
}

// IID_IDxcCompilerArgs - Compiler arguments container
// {73EFFE2A-70DC-45F8-9690-EFF64C02429D}
var IID_IDxcCompilerArgs = GUID{
        0x73effe2a,
        0x70dc,
        0x45f8,
        [8]byte{0x96, 0x90, 0xef, 0xf6, 0x4c, 0x02, 0x42, 0x9d},
}

// ============================================================
// INTERFACE IDs (IIDs) - Library/Utils Interfaces
// ============================================================

// IID_IDxcLibrary - Legacy library interface (deprecated, use IDxcUtils)
// {E5204DC7-D18C-4C3C-BDFB-851673980FE7}
var IID_IDxcLibrary = GUID{
        0xe5204dc7,
        0xd18c,
        0x4c3c,
        [8]byte{0xbd, 0xfb, 0x85, 0x16, 0x73, 0x98, 0x0f, 0xe7},
}

// IID_IDxcUtils - Modern utility interface (replaces IDxcLibrary)
// {4605C4CB-2019-492A-ADA4-65F20BB7D67F}
var IID_IDxcUtils = GUID{
        0x4605c4cb,
        0x2019,
        0x492a,
        [8]byte{0xad, 0xa4, 0x65, 0xf2, 0x0b, 0xb7, 0xd6, 0x7f},
}

// ============================================================
// INTERFACE IDs (IIDs) - Compiler Interfaces
// ============================================================

// IID_IDxcCompiler - Original compiler interface (deprecated)
// {8C210BF3-011F-4422-8D70-6F9ACB8DB617}
var IID_IDxcCompiler = GUID{
        0x8c210bf3,
        0x011f,
        0x4422,
        [8]byte{0x8d, 0x70, 0x6f, 0x9a, 0xcb, 0x8d, 0xb6, 0x17},
}

// IID_IDxcCompiler2 - Extended compiler interface (deprecated)
// {A005A9D9-B8BB-4594-B5C9-0E633BEC4D37}
var IID_IDxcCompiler2 = GUID{
        0xa005a9d9,
        0xb8bb,
        0x4594,
        [8]byte{0xb5, 0xc9, 0x0e, 0x63, 0x3b, 0xec, 0x4d, 0x37},
}

// IID_IDxcCompiler3 - Modern compiler interface (RECOMMENDED)
// {228B4687-5A6A-4730-900C-9702B2203F54}
var IID_IDxcCompiler3 = GUID{
        0x228b4687,
        0x5a6a,
        0x4730,
        [8]byte{0x90, 0x0c, 0x97, 0x02, 0xb2, 0x20, 0x3f, 0x54},
}

// ============================================================
// INTERFACE IDs (IIDs) - Result Interfaces
// ============================================================

// IID_IDxcOperationResult - Legacy result interface
// {CEDB484A-D4E9-445A-B991-CA21CA157DC2}
var IID_IDxcOperationResult = GUID{
        0xcedb484a,
        0xd4e9,
        0x445a,
        [8]byte{0xb9, 0x91, 0xca, 0x21, 0xca, 0x15, 0x7d, 0xc2},
}

// IID_IDxcResult - Modern result interface (replaces IDxcOperationResult)
// {58346CDA-DDE7-4497-9461-6F87AF5E0659}
var IID_IDxcResult = GUID{
        0x58346cda,
        0xdde7,
        0x4497,
        [8]byte{0x94, 0x61, 0x6f, 0x87, 0xaf, 0x5e, 0x06, 0x59},
}

// IID_IDxcExtraOutputs - Additional outputs from DXC operations
// {319B37A2-A5C2-494A-A5DE-4801B2FAF989}
var IID_IDxcExtraOutputs = GUID{
        0x319b37a2,
        0xa5c2,
        0x494a,
        [8]byte{0xa5, 0xde, 0x48, 0x01, 0xb2, 0xfa, 0xf9, 0x89},
}

// ============================================================
// INTERFACE IDs (IIDs) - Validator Interfaces
// ============================================================

// IID_IDxcValidator - Shader validator interface
// {A6E82BD2-1FD7-4826-9811-2857E797F49A}
var IID_IDxcValidator = GUID{
        0xa6e82bd2,
        0x1fd7,
        0x4826,
        [8]byte{0x98, 0x11, 0x28, 0x57, 0xe7, 0x97, 0xf4, 0x9a},
}

// IID_IDxcValidator2 - Extended validator interface
// {458E1FD1-B1B2-4750-A6E1-9C10F03BED92}
var IID_IDxcValidator2 = GUID{
        0x458e1fd1,
        0xb1b2,
        0x4750,
        [8]byte{0xa6, 0xe1, 0x9c, 0x10, 0xf0, 0x3b, 0xed, 0x92},
}

// ============================================================
// INTERFACE IDs (IIDs) - Linker Interfaces
// ============================================================

// IID_IDxcLinker - Shader linker interface
// {F1B5BE2A-62DD-4327-A1C2-42AC1E1E78E6}
var IID_IDxcLinker = GUID{
        0xf1b5be2a,
        0x62dd,
        0x4327,
        [8]byte{0xa1, 0xc2, 0x42, 0xac, 0x1e, 0x1e, 0x78, 0xe6},
}

// ============================================================
// INTERFACE IDs (IIDs) - Container Interfaces
// ============================================================

// IID_IDxcContainerBuilder - Container builder interface
// {334B1F50-2292-4B35-99A1-25588D8C17FE}
var IID_IDxcContainerBuilder = GUID{
        0x334b1f50,
        0x2292,
        0x4b35,
        [8]byte{0x99, 0xa1, 0x25, 0x58, 0x8d, 0x8c, 0x17, 0xfe},
}

// IID_IDxcAssembler - Assembler interface
// {091F7A26-1C1F-4948-904B-E6E3A8A771D5}
var IID_IDxcAssembler = GUID{
        0x091f7a26,
        0x1c1f,
        0x4948,
        [8]byte{0x90, 0x4b, 0xe6, 0xe3, 0xa8, 0xa7, 0x71, 0xd5},
}

// IID_IDxcContainerReflection - Container reflection interface
// {D2C21B26-8350-4BDC-976A-331CE6F4C54C}
var IID_IDxcContainerReflection = GUID{
        0xd2c21b26,
        0x8350,
        0x4bdc,
        [8]byte{0x97, 0x6a, 0x33, 0x1c, 0xe6, 0xf4, 0xc5, 0x4c},
}

// ============================================================
// INTERFACE IDs (IIDs) - Optimizer Interfaces
// ============================================================

// IID_IDxcOptimizerPass - Optimizer pass interface
// {AE2CD79F-CC22-453F-9B6B-B124E7A5204C}
var IID_IDxcOptimizerPass = GUID{
        0xae2cd79f,
        0xcc22,
        0x453f,
        [8]byte{0x9b, 0x6b, 0xb1, 0x24, 0xe7, 0xa5, 0x20, 0x4c},
}

// IID_IDxcOptimizer - Optimizer interface
// {25740E2E-9CBA-401B-9119-4FB42F39F270}
var IID_IDxcOptimizer = GUID{
        0x25740e2e,
        0x9cba,
        0x401b,
        [8]byte{0x91, 0x19, 0x4f, 0xb4, 0x2f, 0x39, 0xf2, 0x70},
}

// ============================================================
// INTERFACE IDs (IIDs) - Version Info Interfaces
// ============================================================

// IID_IDxcVersionInfo - Version information interface
// {B04F5B50-2059-4F12-A8FF-A1E0CDE1CC7E}
var IID_IDxcVersionInfo = GUID{
        0xb04f5b50,
        0x2059,
        0x4f12,
        [8]byte{0xa8, 0xff, 0xa1, 0xe0, 0xcd, 0xe1, 0xcc, 0x7e},
}

// IID_IDxcVersionInfo2 - Extended version information interface
// {FB6904C4-42F0-4B62-9C46-983AF7DA7C83}
var IID_IDxcVersionInfo2 = GUID{
        0xfb6904c4,
        0x42f0,
        0x4b62,
        [8]byte{0x9c, 0x46, 0x98, 0x3a, 0xf7, 0xda, 0x7c, 0x83},
}

// IID_IDxcVersionInfo3 - Version information with custom version string
// {5E13E843-9D25-473C-9AD2-03B2D0B44B1E}
var IID_IDxcVersionInfo3 = GUID{
        0x5e13e843,
        0x9d25,
        0x473c,
        [8]byte{0x9a, 0xd2, 0x03, 0xb2, 0xd0, 0xb4, 0x4b, 0x1e},
}

// ============================================================
// INTERFACE IDs (IIDs) - PDB Interfaces
// ============================================================

// IID_IDxcPdbUtils - PDB utilities (deprecated)
// {E6C9647E-9D6A-4C3B-B94C-524B5A6C343D}
var IID_IDxcPdbUtils = GUID{
        0xe6c9647e,
        0x9d6a,
        0x4c3b,
        [8]byte{0xb9, 0x4c, 0x52, 0x4b, 0x5a, 0x6c, 0x34, 0x3d},
}

// IID_IDxcPdbUtils2 - Modern PDB utilities
// {4315D938-F369-4F93-95A2-252017CC3807}
var IID_IDxcPdbUtils2 = GUID{
        0x4315d938,
        0xf369,
        0x4f93,
        [8]byte{0x95, 0xa2, 0x25, 0x20, 0x17, 0xcc, 0x38, 0x07},
}

// ============================================================
// DXC CONSTANTS - Code Pages
// ============================================================

const (
        // DXC_CP_UTF8 - UTF-8 code page (65001)
        DXC_CP_UTF8 = 65001
        // DXC_CP_UTF16 - UTF-16 code page (1200)
        DXC_CP_UTF16 = 1200
        // DXC_CP_UTF32 - UTF-32 code page (12000)
        DXC_CP_UTF32 = 12000
        // DXC_CP_ACP - ANSI code page, binary, or auto-detect with BOM
        DXC_CP_ACP = 0
        // DXC_CP_WIDE - Wide character encoding (UTF-16 on Windows, UTF-32 on Linux)
        DXC_CP_WIDE = DXC_CP_UTF16
)

// ============================================================
// DXC CONSTANTS - Hash Flags
// ============================================================

const (
        // DXC_HASHFLAG_INCLUDES_SOURCE - Hash computed with source info (-Zss)
        DXC_HASHFLAG_INCLUDES_SOURCE = 1
)

// ============================================================
// DXC CONSTANTS - FourCC Part Identifiers
// ============================================================

const (
        // DXC_PART_PDB - PDB part
        DXC_PART_PDB = 0x42444C49 // 'ILDB'
        // DXC_PART_PDB_NAME - PDB name part
        DXC_PART_PDB_NAME = 0x4E444C49 // 'ILDN'
        // DXC_PART_PRIVATE_DATA - Private data part
        DXC_PART_PRIVATE_DATA = 0x56495250 // 'PRIV'
        // DXC_PART_ROOT_SIGNATURE - Root signature part
        DXC_PART_ROOT_SIGNATURE = 0x30535452 // 'RTS0'
        // DXC_PART_DXIL - DXIL part
        DXC_PART_DXIL = 0x4C495844 // 'DXIL'
        // DXC_PART_REFLECTION_DATA - Reflection data part
        DXC_PART_REFLECTION_DATA = 0x54415453 // 'STAT'
        // DXC_PART_SHADER_HASH - Shader hash part
        DXC_PART_SHADER_HASH = 0x48534148 // 'HASH'
        // DXC_PART_INPUT_SIGNATURE - Input signature part
        DXC_PART_INPUT_SIGNATURE = 0x31475349 // 'ISG1'
        // DXC_PART_OUTPUT_SIGNATURE - Output signature part
        DXC_PART_OUTPUT_SIGNATURE = 0x3147534F // 'OSG1'
        // DXC_PART_PATCH_CONSTANT_SIGNATURE - Patch constant signature part
        DXC_PART_PATCH_CONSTANT_SIGNATURE = 0x31475350 // 'PSG1'
)

// ============================================================
// DXC CONSTANTS - Output Kinds
// ============================================================

const (
        DXC_OUT_NONE             = 0
        DXC_OUT_OBJECT           = 1  // IDxcBlob - Shader or library object
        DXC_OUT_ERRORS           = 2  // IDxcBlobUtf8 or IDxcBlobWide
        DXC_OUT_PDB              = 3  // IDxcBlob
        DXC_OUT_SHADER_HASH      = 4  // IDxcBlob - DxcShaderHash
        DXC_OUT_DISASSEMBLY      = 5  // IDxcBlobUtf8 or IDxcBlobWide
        DXC_OUT_HLSL             = 6  // IDxcBlobUtf8 or IDxcBlobWide
        DXC_OUT_TEXT             = 7  // IDxcBlobUtf8 or IDxcBlobWide
        DXC_OUT_REFLECTION       = 8  // IDxcBlob - RDAT part
        DXC_OUT_ROOT_SIGNATURE   = 9  // IDxcBlob
        DXC_OUT_EXTRA_OUTPUTS    = 10 // IDxcExtraOutputs
        DXC_OUT_REMARKS          = 11 // IDxcBlobUtf8 or IDxcBlobWide
        DXC_OUT_TIME_REPORT      = 12 // IDxcBlobUtf8 or IDxcBlobWide
        DXC_OUT_TIME_TRACE       = 13 // IDxcBlobUtf8 or IDxcBlobWide
        DXC_OUT_LAST             = DXC_OUT_TIME_TRACE
        DXC_OUT_NUM_ENUMS        = DXC_OUT_LAST + 1
        DXC_OUT_FORCE_DWORD      = 0xFFFFFFFF
)

// ============================================================
// DXC CONSTANTS - Validator Flags
// ============================================================

const (
        DxcValidatorFlags_Default         = 0
        DxcValidatorFlags_InPlaceEdit     = 1 // Validator can update shader blob in-place
        DxcValidatorFlags_RootSignatureOnly = 2
        DxcValidatorFlags_ModuleOnly      = 4
        DxcValidatorFlags_ValidMask       = 7
)

// ============================================================
// DXC CONSTANTS - Version Info Flags
// ============================================================

const (
        DxcVersionInfoFlags_None     = 0
        DxcVersionInfoFlags_Debug    = 1 // Matches VS_FF_DEBUG
        DxcVersionInfoFlags_Internal = 2 // Internal Validator (non-signing)
)

// ============================================================
// IID COLLECTIONS - For Iteration/Probing
// ============================================================

// IIDEntry stores IID with its name for debugging
type IIDEntry struct {
        Name string
        GUID GUID
}

// UtilsIIDs - IDxcUtils/IDxcLibrary interface IDs (in order of preference)
var UtilsIIDs = []IIDEntry{
        {"IDxcUtils", IID_IDxcUtils},
        {"IDxcLibrary", IID_IDxcLibrary},
}

// CompilerIIDs - IDxcCompiler interface IDs (in order of preference)
var CompilerIIDs = []IIDEntry{
        {"IDxcCompiler3", IID_IDxcCompiler3},
        {"IDxcCompiler2", IID_IDxcCompiler2},
        {"IDxcCompiler", IID_IDxcCompiler},
}

// ResultIIDs - IDxcResult interface IDs (in order of preference)
var ResultIIDs = []IIDEntry{
        {"IDxcResult", IID_IDxcResult},
        {"IDxcOperationResult", IID_IDxcOperationResult},
}

// BlobIIDs - IDxcBlob interface IDs
var BlobIIDs = []IIDEntry{
        {"IDxcBlob", IID_IDxcBlob},
        {"IDxcBlobEncoding", IID_IDxcBlobEncoding},
        {"IDxcBlobUtf8", IID_IDxcBlobUtf8},
        {"IDxcBlobWide", IID_IDxcBlobWide},
}

// ValidatorIIDs - IDxcValidator interface IDs
var ValidatorIIDs = []IIDEntry{
        {"IDxcValidator2", IID_IDxcValidator2},
        {"IDxcValidator", IID_IDxcValidator},
}
