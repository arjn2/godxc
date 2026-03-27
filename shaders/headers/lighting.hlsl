// Lighting calculations and utilities

#ifndef LIGHTING_HLSL
#define LIGHTING_HLSL

// ============================================================
// LIGHT STRUCTURES
// ============================================================

struct DirectionalLight {
    float3 Direction;
    float _padding1;
    float3 Color;
    float Intensity;
};

struct PointLight {
    float3 Position;
    float Range;
    float3 Color;
    float Intensity;
};

struct SpotLight {
    float3 Position;
    float Range;
    float3 Direction;
    float InnerCone;
    float3 Color;
    float OuterCone;
    float Intensity;
};

// ============================================================
// LIGHTING CALCULATIONS
// ============================================================

float3 CalculateDirectionalLight(float3 normal, DirectionalLight light) {
    float diff = max(dot(normal, -light.Direction), 0.0);
    return light.Color * diff * light.Intensity;
}

float CalculatePointLightAttenuation(float3 position, PointLight light) {
    float distance = length(light.Position - position);
    float attenuation = 1.0 / (1.0 + distance * distance / (light.Range * light.Range));
    return attenuation;
}

float3 CalculatePointLight(float3 position, float3 normal, PointLight light) {
    float3 lightDir = normalize(light.Position - position);
    float diff = max(dot(normal, lightDir), 0.0);
    float attenuation = CalculatePointLightAttenuation(position, light);
    return light.Color * diff * light.Intensity * attenuation;
}

#endif
