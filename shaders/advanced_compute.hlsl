// Advanced Compute Shader with shared memory and synchronization
// Demonstrates parallel reduction pattern

cbuffer Constants : register(b0) {
    uint ElementCount;
    float Scale;
    uint Padding1;
    uint Padding2;
};

StructuredBuffer<float> InputA : register(t0);
StructuredBuffer<float> InputB : register(t1);
RWStructuredBuffer<float> Output : register(u0);
RWStructuredBuffer<float> SumOutput : register(u1);

// Shared memory for parallel reduction
groupshared float SharedData[256];

[numthreads(256, 1, 1)]
void main(uint3 DTid : SV_DispatchThreadID, uint3 GTid : SV_GroupThreadID) {
    uint index = DTid.x;
    
    // Bounds check
    if (index >= ElementCount) {
        SharedData[GTid.x] = 0.0f;
    } else {
        // Load data and compute
        float a = InputA[index];
        float b = InputB[index];
        float result = (a + b) * Scale;
        
        Output[index] = result;
        SharedData[GTid.x] = result;
    }
    
    // Synchronize all threads in the group
    GroupMemoryBarrierWithGroupSync();
    
    // Parallel reduction in shared memory
    for (uint stride = 128; stride > 0; stride >>= 1) {
        if (GTid.x < stride) {
            SharedData[GTid.x] += SharedData[GTid.x + stride];
        }
        GroupMemoryBarrierWithGroupSync();
    }
    
    // Write sum result
    if (GTid.x == 0) {
        SumOutput[DTid.z] = SharedData[0];
    }
}
