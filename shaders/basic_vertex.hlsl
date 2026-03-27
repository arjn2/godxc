// Basic Vertex Shader
// Transforms vertex position and passes through UV coordinates
// MODIFIED: Adding animation test - HOT RELOAD TEST #1

cbuffer Constants : register(b0) {
    float4x4 WorldViewProj;
    float Time;  // NEW: Added time parameter for animation testing
};

struct VSInput {
    float3 Position : POSITION;
    float2 UV : TEXCOORD0;
    float3 Normal : NORMAL;
};

struct VSOutput {
    float4 Position : SV_Position;
    float2 UV : TEXCOORD0;
    float3 Normal : NORMAL;
};

VSOutput main(VSInput input) {
    VSOutput output;
    
    // Transform position with enhanced animation
    float3 animatedPos = input.Position + float3(sin(Time) * 0.05, cos(Time) * 0.1, sin(Time * 2.0) * 0.05);
    output.Position = mul(WorldViewProj, float4(animatedPos, 1.0));
    
    // Pass through UV and normal
    output.UV = input.UV;
    output.Normal = input.Normal;
    
    return output;
}
