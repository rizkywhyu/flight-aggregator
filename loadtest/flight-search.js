import http from 'k6/http';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '1m', target: 50 },
    { duration: '30s', target: 0 },
  ],
};

const BASE_URL = 'http://localhost:8080';

export default function () {
  const payload = JSON.stringify({
    origin: 'CGK',
    destination: 'DPS',
    departureDate: '2025-12-15',
    passengers: 1,
    cabinClass: 'economy'
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'X-Tracer-ID': `test-${__VU}-${__ITER}`
    },
  };

  const response = http.post(`${BASE_URL}/api/flights/search`, payload, params);
  
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 2000ms': (r) => r.timings.duration < 2000,
    'has flights': (r) => JSON.parse(r.body).flights.length > 0,
  });
}