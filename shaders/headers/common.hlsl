// Common shader definitions and utilities

#ifndef COMMON_HLSL
#define COMMON_HLSL

// ============================================================
// CONSTANTS
// ============================================================

#define PI 3.14159265359
#define TWO_PI 6.28318530718
#define HALF_PI 1.57079632679

// ============================================================
// COMMON STRUCTURES
// ============================================================

// Standard vertex input structure
struct VertexInput {
    float3 Position : POSITION;
    float2 UV : TEXCOORD0;
    float3 Normal : NORMAL;
    float3 Tangent : TANGENT;
};

// Standard vertex output structure
struct VertexOutput {
    float4 Position : SV_Position;
    float2 UV : TEXCOORD0;
    float3 Normal : TEXCOORD1;
    float3 WorldPos : TEXCOORD2;
};

// Pixel shader input (matches vertex output)
struct PixelInput {
    float4 Position : SV_Position;
    float2 UV : TEXCOORD0;
    float3 Normal : TEXCOORD1;
    float3 WorldPos : TEXCOORD2;
};

// ============================================================
// MATH UTILITIES
// ============================================================

float3 Normalize(float3 v) {
    return normalize(v);
}

float Length(float3 v) {
    return length(v);
}

float3 Reflect(float3 incident, float3 normal) {
    return reflect(incident, normal);
}

float3 Refract(float3 incident, float3 normal, float ior) {
    return refract(incident, normal, ior);
}

#endif
