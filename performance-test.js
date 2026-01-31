import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const loginDuration = new Trend('login_duration');
const profileDuration = new Trend('get_profile_duration');
const limitsDuration = new Trend('get_limits_duration');
const transactionDuration = new Trend('create_transaction_duration');
const getTransactionsDuration = new Trend('get_transactions_duration');

// Configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const API_KEY = __ENV.API_KEY || 'biytf7rciyubyt6r7g89py';

// Test options with multiple scenarios
export const options = {
    scenarios: {
        // Smoke test - minimal load to verify system works
        smoke: {
            executor: 'constant-vus',
            vus: 1,
            duration: '10s',
            startTime: '0s',
            tags: { test_type: 'smoke' },
        },
        // Load test - normal expected load
        load: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '30s', target: 10 },  // Ramp up
                { duration: '1m', target: 10 },   // Stay at 10 VUs
                { duration: '30s', target: 0 },   // Ramp down
            ],
            startTime: '10s',
            tags: { test_type: 'load' },
        },
        // Stress test - beyond normal load
        stress: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '30s', target: 15 },  // Ramp up to 15
                { duration: '1m', target: 15 },   // Stay at 15 VUs
                { duration: '30s', target: 20 },  // Push to 20
                { duration: '1m', target: 20 },   // Stay at 20 VUs
                { duration: '30s', target: 0 },   // Ramp down
            ],
            startTime: '2m30s',
            tags: { test_type: 'stress' },
        },
    },

    // Performance thresholds (SLA)
    thresholds: {
        'http_req_duration': ['p(95)<500', 'p(99)<1000'],  // 95% under 500ms, 99% under 1s
        'http_req_failed': ['rate<0.01'],                   // Less than 1% errors
        'errors': ['rate<0.05'],                            // Custom error rate under 5%
        'login_duration': ['p(95)<300'],
        'get_profile_duration': ['p(95)<200'],
        'get_limits_duration': ['p(95)<200'],
        'create_transaction_duration': ['p(95)<500'],
        'get_transactions_duration': ['p(95)<300'],
    },
};

// Test users
const users = [
    { email: 'budi@mail.com', password: 'pAsswj@1873' },
    { email: 'annisa@mail.com', password: 'pAsswj@1763' },
];

// Setup - runs once before all VUs
export function setup() {
    console.log('Setting up test...');

    // Login with first user to verify system is up
    const loginRes = http.post(`${BASE_URL}/api/auth/login`, JSON.stringify({
        email: users[0].email,
        password: users[0].password,
    }), {
        headers: { 'Content-Type': 'application/json' },
    });

    const setupCheck = check(loginRes, {
        'setup: login successful': (r) => r.status === 200,
    });

    if (!setupCheck) {
        console.error('Setup failed! Make sure the server is running.');
        return { ready: false };
    }

    return { ready: true };
}

// Main test function
export default function (data) {
    if (!data.ready) {
        console.error('Test aborted - setup failed');
        return;
    }

    // Randomly select a user
    const user = users[Math.floor(Math.random() * users.length)];
    let token = null;

    // Group 1: Authentication
    group('Authentication', function () {
        const loginPayload = JSON.stringify({
            email: user.email,
            password: user.password,
        });

        const loginRes = http.post(`${BASE_URL}/api/auth/login`, loginPayload, {
            headers: { 'Content-Type': 'application/json' },
        });

        loginDuration.add(loginRes.timings.duration);

        const loginSuccess = check(loginRes, {
            'login: status is 200': (r) => r.status === 200,
            'login: has access_token': (r) => r.json('access_token') !== undefined,
        });

        errorRate.add(!loginSuccess);

        if (loginSuccess) {
            token = loginRes.json('access_token');
        }
    });

    if (!token) {
        console.error('Login failed, skipping remaining tests');
        return;
    }

    const authHeaders = {
        'Content-Type': 'application/json',
        'X-API-KEY': API_KEY,
        'Authorization': `Bearer ${token}`,
    };

    sleep(0.5);

    // Group 2: Read Operations
    group('Read Operations', function () {
        // Get Profile
        const profileRes = http.get(`${BASE_URL}/api/user/profile`, {
            headers: authHeaders,
        });

        profileDuration.add(profileRes.timings.duration);

        const profileSuccess = check(profileRes, {
            'get profile: status is 200': (r) => r.status === 200,
            'get profile: has data': (r) => r.json('data') !== undefined,
        });

        errorRate.add(!profileSuccess);

        sleep(0.3);

        // Get Limits
        const limitsRes = http.get(`${BASE_URL}/api/limit/`, {
            headers: authHeaders,
        });

        limitsDuration.add(limitsRes.timings.duration);

        const limitsSuccess = check(limitsRes, {
            'get limits: status is 200': (r) => r.status === 200,
        });

        errorRate.add(!limitsSuccess);

        sleep(0.3);

        // Get Transactions
        const transactionsRes = http.get(`${BASE_URL}/api/transaction/`, {
            headers: authHeaders,
        });

        getTransactionsDuration.add(transactionsRes.timings.duration);

        const transactionsSuccess = check(transactionsRes, {
            'get transactions: status is 200': (r) => r.status === 200,
        });

        errorRate.add(!transactionsSuccess);
    });

    sleep(0.5);

    // Group 3: Write Operations (Create Transaction)
    group('Write Operations', function () {
        const tenors = [1, 2, 3, 6];
        const randomTenor = tenors[Math.floor(Math.random() * tenors.length)];
        const randomOTR = Math.floor(Math.random() * 5000) + 1000;

        const payload = JSON.stringify({
            contract_number: `CTR-K6-${__VU}-${__ITER}-${Date.now()}`,
            otr: randomOTR,
            admin_fee: 10000,
            installment_amount: Math.floor(randomOTR / randomTenor) + 5000,
            interest_amount: 10000,
            asset_name: 'K6 Performance Test Asset',
            tenor: randomTenor,
        });

        const createRes = http.post(`${BASE_URL}/api/transaction/`, payload, {
            headers: authHeaders,
        });

        transactionDuration.add(createRes.timings.duration);

        // 201 = Created, 400 = Insufficient limit (expected when limit exhausted)
        const createSuccess = check(createRes, {
            'create transaction: valid response': (r) => r.status === 201 || r.status === 400,
        });

        if (createRes.status !== 201 && createRes.status !== 400) {
            console.error(`Unexpected status: ${createRes.status} - ${createRes.body}`);
            errorRate.add(true);
        } else {
            errorRate.add(false);
        }
    });

    sleep(1);
}

// Teardown - runs once after all VUs finish
export function teardown(data) {
    console.log('Test completed!');
    console.log('Check the results above for performance metrics.');
}

// Health check scenario (can run separately)
export function healthCheck() {
    const res = http.get(`${BASE_URL}/health`);
    check(res, {
        'health: status is 200': (r) => r.status === 200,
        'health: status is UP': (r) => r.json('status') === 'UP',
    });
}
