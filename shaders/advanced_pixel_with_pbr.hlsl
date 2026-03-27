// Advanced Pixel Shader using PBR headers
// Include paths: -I shaders/headers

#include "common.hlsl"
#include "pbr.hlsl"
#include "lighting.hlsl"

cbuffer MaterialBuffer : register(b0) {
    float4 Albedo;
    float Metallic;
    float Roughness;
};

cbuffer LightBuffer : register(b1) {
    DirectionalLight MainLight;
};

Texture2D AlbedoTexture : register(t0);
Texture2D NormalTexture : register(t1);
Texture2D RoughnessTexture : register(t2);
SamplerState MainSampler : register(s0);

float4 main(PixelInput input) : SV_Target {
    // Sample textures
    float3 albedo = AlbedoTexture.Sample(MainSampler, input.UV).rgb;
    float3 normal = normalize(NormalTexture.Sample(MainSampler, input.UV).rgb * 2.0 - 1.0);
    float roughness = RoughnessTexture.Sample(MainSampler, input.UV).r;
    
    // Calculate lighting
    float3 lighting = CalculateDirectionalLight(normal, MainLight);
    
    // Combine with albedo
    float3 finalColor = albedo * lighting;
    
    return float4(finalColor, 1.0);
}
