import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 5,
    duration: '1s', // Try to hit server simultaneously
    iterations: 5,  // Total 5 requests
};

const BASE_URL = 'http://localhost:8080/api';
const API_KEY = 'biytf7rciyubyt6r7g89py';

export function setup() {
    const loginPayload = JSON.stringify({
        email: 'budi@mail.com',
        password: 'pAsswj@1873',
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(`${BASE_URL}/auth/login`, loginPayload, params);
    check(res, { 'login successful': (r) => r.status === 200 });

    return { token: res.json('access_token') };
}

export default function (data) {
    const payload = JSON.stringify({
        contract_number: `CTR-K6-${__VU}-${__ITER}`,
        otr: 600000,
        admin_fee: 10000,
        installment_amount: 105000,
        interest_amount: 10000,
        asset_name: 'K6 Test Asset',
        tenor: 6,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'X-API-KEY': API_KEY,
            'Authorization': `Bearer ${data.token}`,
        },
    };

    const res = http.post(`${BASE_URL}/transaction/`, payload, params);

    // We expect mostly 400 (Insufficient Limit) and exactly one 201 (Created)
    // But due to network, maybe more than one 201 if race condition exists!
    check(res, {
        'status is 201 or 400': (r) => r.status === 201 || r.status === 400,
    });

    if (res.status === 201) {
        console.log(`VU ${__VU} created transaction successfully!`);
    } else if (res.status === 400) {
        console.log(`VU ${__VU} failed: ${res.json('error')}`);
    } else {
        console.log(`VU ${__VU} unexpected status: ${res.status} body: ${res.body}`);
    }
}
