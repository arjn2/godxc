// Basic Pixel Shader
// Samples a texture and applies simple lighting

Texture2D DiffuseTexture : register(t0);
SamplerState LinearSampler : register(s0);

cbuffer LightingConstants : register(b0) {
    float3 LightDirection;
    float LightIntensity;
    float3 LightColor;
    float AmbientIntensity;
};

struct PSInput {
    float4 Position : SV_Position;
    float2 UV : TEXCOORD0;
    float3 Normal : NORMAL;
};

float4 main(PSInput input) : SV_Target {
    // Sample diffuse texture
    float4 diffuse = DiffuseTexture.Sample(LinearSampler, input.UV);
    
    // Simple diffuse lighting
    float3 normal = normalize(input.Normal);
    float NdotL = max(dot(normal, -LightDirection), 0.0);
    
    // Combine ambient and diffuse lighting
    float3 ambient = diffuse.rgb * AmbientIntensity;
    float3 direct = diffuse.rgb * NdotL * LightIntensity * LightColor;
    
    return float4(ambient + direct, diffuse.a);
}
