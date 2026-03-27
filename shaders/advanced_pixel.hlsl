// Advanced Pixel Shader with material and lighting
// Demonstrates physically-based rendering concepts

cbuffer MaterialBuffer : register(b0) {
    float3 Albedo;
    float Metallic;
    float3 Normal;
    float Roughness;
};

cbuffer LightBuffer : register(b1) {
    float3 LightPos;
    float LightRange;
    float3 LightColor;
    float LightIntensity;
    float3 CameraPos;
    float Padding;
};

Texture2D<float4> AlbedoMap : register(t0);
Texture2D<float4> NormalMap : register(t1);
Texture2D<float4> RoughnessMap : register(t2);
SamplerState LinearSampler : register(s0);

struct PSInput {
    float4 Position : SV_Position;
    float3 Normal : NORMAL;
    float3 Tangent : TANGENT;
    float2 UV : TEXCOORD0;
    float4 Color : COLOR;
    float3 WorldPos : TEXCOORD1;
};

// Simple PBR calculation
float3 CalculateLighting(float3 N, float3 L, float3 V, float3 LightCol, float3 Albedo) {
    float NdotL = max(dot(N, L), 0.0f);
    float NdotV = max(dot(N, V), 0.0f);
    
    // Simple diffuse
    float3 diffuse = Albedo * LightCol * NdotL;
    
    // Simple specular
    float3 H = normalize(L + V);
    float NdotH = max(dot(N, H), 0.0f);
    float3 specular = LightCol * pow(NdotH, 32.0f) * 0.5f;
    
    return diffuse + specular;
}

float4 main(PSInput input) : SV_Target {
    // Sample textures
    float3 albedo = AlbedoMap.Sample(LinearSampler, input.UV).rgb;
    float3 normalMap = NormalMap.Sample(LinearSampler, input.UV).rgb;
    float roughness = RoughnessMap.Sample(LinearSampler, input.UV).r;
    
    // Reconstruct normal from map
    float3 N = normalize(normalMap * 2.0f - 1.0f);
    
    // Calculate lighting direction
    float3 L = normalize(LightPos - input.WorldPos);
    float3 V = normalize(CameraPos - input.WorldPos);
    
    // Calculate lighting
    float3 lighting = CalculateLighting(N, L, V, float3(1.0f, 1.0f, 1.0f), albedo);
    
    // Add ambient
    float3 ambient = albedo * 0.1f;
    
    // Final color
    float3 finalColor = ambient + lighting;
    
    return float4(finalColor, 1.0f);
}
