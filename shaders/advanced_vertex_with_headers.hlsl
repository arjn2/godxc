// Advanced Vertex Shader using common headers
// Include paths: -I shaders/headers

#include "common.hlsl"

cbuffer TransformBuffer : register(b0) {
    float4x4 WorldViewProj;
    float4x4 World;
    float4x4 NormalMatrix;
};

VertexOutput main(VertexInput input) {
    VertexOutput output;
    
    // Transform position to clip space
    output.Position = mul(float4(input.Position, 1.0), WorldViewProj);
    
    // Transform normal to world space
    output.Normal = normalize(mul(input.Normal, (float3x3)NormalMatrix));
    
    // Calculate world position
    output.WorldPos = mul(float4(input.Position, 1.0), World).xyz;
    
    // Pass through UV coordinates
    output.UV = input.UV;
    
    return output;
}
