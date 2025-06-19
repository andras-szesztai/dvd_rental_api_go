const autocannon = require('autocannon');

// Test configuration
const BASE_URL = 'http://localhost:8080';
const TEST_ENDPOINT = '/v1/test';

// Rate limiter configuration (should match your Go app settings)
const RATE_LIMIT = 100; // requests per time frame
const TIME_FRAME = 60; // seconds

console.log('🚀 Starting Rate Limiter Test');
console.log(`📊 Rate Limit: ${RATE_LIMIT} requests per ${TIME_FRAME} seconds`);
console.log(`🎯 Target: ${BASE_URL}${TEST_ENDPOINT}`);
console.log('');

// Test 1: Normal load (should pass)
async function testNormalLoad() {
    console.log('📈 Test 1: Normal Load (50 requests)');
    console.log('Expected: All requests should succeed');
    console.log('');
    
    const result = await autocannon({
        url: `${BASE_URL}${TEST_ENDPOINT}`,
        connections: 10,
        duration: 5,
        pipelining: 1,
        headers: {
            'Content-Type': 'application/json'
        }
    });
    
    console.log(`✅ Requests: ${result.requests.total}`);
    console.log(`✅ 2xx responses: ${result['2xx']}`);
    console.log(`❌ 4xx responses: ${result['4xx']}`);
    console.log(`❌ 5xx responses: ${result['5xx']}`);
    console.log(`⏱️  Average latency: ${result.latency.average}ms`);
    console.log('');
}

// Test 2: High load (should hit rate limit)
async function testHighLoad() {
    console.log('🔥 Test 2: High Load (200 requests)');
    console.log('Expected: Some requests should be rate limited (429 responses)');
    console.log('');
    
    const result = await autocannon({
        url: `${BASE_URL}${TEST_ENDPOINT}`,
        connections: 20,
        duration: 10,
        pipelining: 1,
        headers: {
            'Content-Type': 'application/json'
        }
    });
    
    console.log(`📊 Requests: ${result.requests.total}`);
    console.log(`✅ 2xx responses: ${result['2xx']}`);
    console.log(`❌ 4xx responses: ${result['4xx']}`);
    console.log(`❌ 5xx responses: ${result['5xx']}`);
    console.log(`⏱️  Average latency: ${result.latency.average}ms`);
    console.log('');
}

// Test 3: Burst load (should definitely hit rate limit)
async function testBurstLoad() {
    console.log('💥 Test 3: Burst Load (500 requests)');
    console.log('Expected: Many requests should be rate limited (429 responses)');
    console.log('');
    
    const result = await autocannon({
        url: `${BASE_URL}${TEST_ENDPOINT}`,
        connections: 50,
        duration: 5,
        pipelining: 1,
        headers: {
            'Content-Type': 'application/json'
        }
    });
    
    console.log(`📊 Requests: ${result.requests.total}`);
    console.log(`✅ 2xx responses: ${result['2xx']}`);
    console.log(`❌ 4xx responses: ${result['4xx']}`);
    console.log(`❌ 5xx responses: ${result['5xx']}`);
    console.log(`⏱️  Average latency: ${result.latency.average}ms`);
    console.log('');
}

// Test 4: Sustained load over time frame
async function testSustainedLoad() {
    console.log('⏰ Test 4: Sustained Load (over rate limit window)');
    console.log('Expected: Should see rate limiting behavior over time');
    console.log('');
    
    const result = await autocannon({
        url: `${BASE_URL}${TEST_ENDPOINT}`,
        connections: 15,
        duration: TIME_FRAME + 5, // Test for slightly longer than the time frame
        pipelining: 1,
        headers: {
            'Content-Type': 'application/json'
        }
    });
    
    console.log(`📊 Requests: ${result.requests.total}`);
    console.log(`✅ 2xx responses: ${result['2xx']}`);
    console.log(`❌ 4xx responses: ${result['4xx']}`);
    console.log(`❌ 5xx responses: ${result['5xx']}`);
    console.log(`⏱️  Average latency: ${result.latency.average}ms`);
    console.log('');
}

// Run all tests
async function runAllTests() {
    try {
        await testNormalLoad();
        await new Promise(resolve => setTimeout(resolve, 2000)); // Wait between tests
        
        await testHighLoad();
        await new Promise(resolve => setTimeout(resolve, 2000)); // Wait between tests
        
        await testBurstLoad();
        await new Promise(resolve => setTimeout(resolve, 2000)); // Wait between tests
        
        await testSustainedLoad();
        
        console.log('🎉 All tests completed!');
        console.log('');
        console.log('📋 Summary:');
        console.log('- Test 1: Normal load should pass');
        console.log('- Test 2: High load should show some rate limiting');
        console.log('- Test 3: Burst load should show significant rate limiting');
        console.log('- Test 4: Sustained load should show rate limiting over time');
        
    } catch (error) {
        console.error('❌ Test failed:', error.message);
        process.exit(1);
    }
}

// Check if server is running
async function checkServer() {
    try {
        const result = await autocannon({
            url: `${BASE_URL}/v1/health`,
            connections: 1,
            duration: 1,
            pipelining: 1
        });
        
        if (result.errors > 0) {
            throw new Error('Server not responding');
        }
        
        console.log('✅ Server is running and responding');
        console.log('');
        return true;
    } catch (error) {
        console.error('❌ Server is not running or not responding');
        console.error('Please start your Go API server first:');
        console.error('  go run cmd/api/main.go');
        console.error('');
        process.exit(1);
    }
}

// Main execution
async function main() {
    await checkServer();
    await runAllTests();
}

main(); 