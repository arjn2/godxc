// Basic Compute Shader
// Demonstrates a simple parallel addition

StructuredBuffer<float> InputA : register(t0);
StructuredBuffer<float> InputB : register(t1);
RWStructuredBuffer<float> Output : register(u0);

cbuffer Constants : register(b0) {
    uint ElementCount;
    float Scale;
};

[numthreads(256, 1, 1)]
void main(uint3 DTid : SV_DispatchThreadID) {
    uint index = DTid.x;
    
    // Bounds check
    if (index >= ElementCount) {
        return;
    }
    
    // Simple computation: C = (A + B) * Scale
    float a = InputA[index];
    float b = InputB[index];
    Output[index] = (a + b) * Scale;
}
