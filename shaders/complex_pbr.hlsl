// Complex PBR Shader - Production Quality
// Physically Based Rendering with multiple light types

#include "headers/common.hlsl"
#include "headers/lighting.hlsl"
#include "headers/pbr.hlsl"

cbuffer SceneConstants : register(b0) {
    float4x4 ViewProj;
    float4x4 World;
    float4x4 ViewInverse;
    float3 CameraPosition;
    float Padding0;
    float3 AmbientLight;
    float Padding1;
};

cbuffer LightConstants : register(b1) {
    float3 MainLightDirection;
    float MainLightIntensity;
    float3 MainLightColor;
    float Padding2;
    uint NumPointLights;
    uint NumSpotLights;
    float2 Padding3;
};

StructuredBuffer<PointLight> PointLights : register(t0);
StructuredBuffer<SpotLight> SpotLights : register(t1);

Texture2D AlbedoTexture : register(t2);
Texture2D NormalTexture : register(t3);
Texture2D RoughnessTexture : register(t4);
Texture2D MetallicTexture : register(t5);
Texture2D AOTexture : register(t6);

Texture2D IrradianceMap : register(t7);
Texture2D PrefilterMap : register(t8);
Texture2D BRDFLookupTexture : register(t9);

SamplerState LinearSampler : register(s0);
SamplerState PointSampler : register(s1);

struct VSInput {
    float3 Position : POSITION;
    float2 UV : TEXCOORD0;
    float3 Normal : NORMAL;
    float3 Tangent : TANGENT;
    float3 BiTangent : BITANGENT;
};

struct PSInput {
    float4 Position : SV_POSITION;
    float2 UV : TEXCOORD0;
    float3 WorldPos : TEXCOORD1;
    float3 WorldNormal : TEXCOORD2;
    float3 Tangent : TEXCOORD3;
    float3 BiTangent : TEXCOORD4;
};

PSInput VSMain(VSInput input) {
    PSInput output;
    
    float4 worldPos = mul(float4(input.Position, 1.0), World);
    output.WorldPos = worldPos.xyz;
    output.Position = mul(worldPos, ViewProj);
    
    output.UV = input.UV;
    output.WorldNormal = mul(input.Normal, (float3x3)World);
    output.Tangent = mul(input.Tangent, (float3x3)World);
    output.BiTangent = mul(input.BiTangent, (float3x3)World);
    
    return output;
}

float4 PSMain(PSInput input) : SV_TARGET {
    // Sample all textures
    float3 albedo = AlbedoTexture.Sample(LinearSampler, input.UV).rgb;
    float3 normalMap = NormalTexture.Sample(LinearSampler, input.UV).rgb;
    float roughness = RoughnessTexture.Sample(LinearSampler, input.UV).r;
    float metallic = MetallicTexture.Sample(LinearSampler, input.UV).r;
    float ao = AOTexture.Sample(LinearSampler, input.UV).r;
    
    // Reconstruct normal from normal map
    normalMap = normalize(normalMap * 2.0 - 1.0);
    float3x3 TBN = float3x3(
        normalize(input.Tangent),
        normalize(input.BiTangent),
        normalize(input.WorldNormal)
    );
    float3 normal = normalize(mul(normalMap, TBN));
    
    // View direction
    float3 viewDir = normalize(CameraPosition - input.WorldPos);
    
    // Initialize color with ambient
    float3 color = AmbientLight * albedo * ao;
    
    // Main directional light
    {
        float3 lightDir = normalize(-MainLightDirection);
        float3 halfway = normalize(viewDir + lightDir);
        
        float ndl = max(dot(normal, lightDir), 0.0);
        float ndh = max(dot(normal, halfway), 0.0);
        float vdh = max(dot(viewDir, halfway), 0.0);
        
        float3 F = FresnelSchlick(max(dot(halfway, viewDir), 0.0), mix(float3(0.04), albedo, metallic));
        float D = DistributionGGX(ndh, roughness);
        float G = GeometrySmith(max(dot(normal, viewDir), 0.0), ndl, roughness);
        
        float3 kS = F;
        float3 kD = (1.0 - kS) * (1.0 - metallic);
        
        float3 specular = (D * F * G) / max(4.0 * max(dot(normal, viewDir), 0.0) * ndl, 0.001);
        float3 diffuse = kD * albedo / 3.14159;
        
        color += (diffuse + specular) * MainLightColor * MainLightIntensity * ndl;
    }
    
    // Point lights
    for (uint i = 0; i < NumPointLights && i < 16; ++i) {
        PointLight light = PointLights[i];
        float3 toLight = light.Position - input.WorldPos;
        float distance = length(toLight);
        float attenuation = 1.0 / (distance * distance + 0.1);
        
        float3 lightDir = normalize(toLight);
        float3 halfway = normalize(viewDir + lightDir);
        
        float ndl = max(dot(normal, lightDir), 0.0);
        float ndh = max(dot(normal, halfway), 0.0);
        
        if (ndl > 0.0) {
            float3 F = FresnelSchlick(max(dot(halfway, viewDir), 0.0), mix(float3(0.04), albedo, metallic));
            float D = DistributionGGX(ndh, roughness);
            float G = GeometrySmith(max(dot(normal, viewDir), 0.0), ndl, roughness);
            
            float3 kS = F;
            float3 kD = (1.0 - kS) * (1.0 - metallic);
            
            float3 specular = (D * F * G) / max(4.0 * max(dot(normal, viewDir), 0.0) * ndl, 0.001);
            float3 diffuse = kD * albedo / 3.14159;
            
            color += (diffuse + specular) * light.Color * light.Intensity * ndl * attenuation;
        }
    }
    
    // Apply gamma correction
    color = pow(color, 1.0 / 2.2);
    
    return float4(color, 1.0);
}
