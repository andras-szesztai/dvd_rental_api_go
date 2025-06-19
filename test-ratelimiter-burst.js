const autocannon = require('autocannon');

console.log('🚀 Quick Rate Limiter Burst Test');
console.log('🎯 Target: http://localhost:8080/v1/test');
console.log('');

// Quick burst test
async function runBurstTest() {
    console.log('💥 Running burst test (100 requests in 2 seconds)...');
    console.log('Expected: Some requests should be rate limited (429 responses)');
    console.log('');
    
    const result = await autocannon({
        url: 'http://localhost:8080/v1/test',
        connections: 20,
        duration: 2,
        pipelining: 1,
        headers: {
            'Content-Type': 'application/json'
        }
    });
    
    console.log('📊 Results:');
    console.log(`   Total Requests: ${result.requests.total}`);
    console.log(`   ✅ 2xx responses: ${result['2xx']}`);
    console.log(`   ❌ 4xx responses: ${result['4xx']}`);
    console.log(`   ❌ 5xx responses: ${result['5xx']}`);
    console.log(`   ⏱️  Average latency: ${result.latency.average}ms`);
    console.log(`   🚀 Requests/sec: ${result.requests.average}`);
    console.log('');
    
    if (result['4xx'] > 0) {
        console.log('🎉 Rate limiter is working! Some requests were blocked.');
    } else {
        console.log('⚠️  No rate limiting detected. Check your configuration.');
    }
}

runBurstTest().catch(console.error); 