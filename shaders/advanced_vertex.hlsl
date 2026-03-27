// Advanced Vertex Shader with multiple features
// Demonstrates various HLSL features for compilation testing

cbuffer TransformBuffer : register(b0) {
    float4x4 WorldViewProj;
    float4x4 World;
    float4 CameraPos;
};

cbuffer LightBuffer : register(b1) {
    float3 LightDir;
    float LightIntensity;
    float3 LightColor;
    float Padding;
};

Texture2D<float4> NormalMap : register(t0);
SamplerState LinearSampler : register(s0);

struct VSInput {
    float3 Position : POSITION;
    float3 Normal : NORMAL;
    float3 Tangent : TANGENT;
    float2 UV : TEXCOORD0;
    float4 Color : COLOR;
};

struct VSOutput {
    float4 Position : SV_Position;
    float3 Normal : NORMAL;
    float3 Tangent : TANGENT;
    float2 UV : TEXCOORD0;
    float4 Color : COLOR;
    float3 WorldPos : TEXCOORD1;
};

VSOutput main(VSInput input) {
    VSOutput output;
    
    // Transform to world space
    float4 worldPos = mul(float4(input.Position, 1.0f), World);
    output.WorldPos = worldPos.xyz;
    
    // Project to screen space
    output.Position = mul(float4(input.Position, 1.0f), WorldViewProj);
    
    // Transform normal and tangent
    output.Normal = mul(input.Normal, (float3x3)World);
    output.Tangent = mul(input.Tangent, (float3x3)World);
    
    // Pass through texture coordinates and color
    output.UV = input.UV;
    output.Color = input.Color;
    
    return output;
}
