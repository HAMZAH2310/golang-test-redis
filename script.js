import http from 'k6/http';
import { check, sleep, group } from 'k6';

export const options = {
  scenarios: {
    api_load_test: {
      executor: 'constant-arrival-rate',
      rate: 100,            // 100 request per detik
      timeUnit: '1s',
      duration: '1m',
      preAllocatedVUs: 50,
      maxVUs: 200,
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],        // < 1% error
    http_req_duration: ['p(95)<400'],      // p95 < 400ms
  },
};

const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzAyMDg2MTgsImlkIjo5LCJ1c2VyX2lkIjo5fQ.Fuq7ZNqyFmf3yI4GcgLbV5R5kPMtzy0rx1hX6CeAm0E';

export default function () {
  group('GET /api/v1/me', () => {
    const res = http.get(
      'http://localhost:3000/api/v1/me',
      {
        headers: {
          Authorization: `Bearer ${token}`,
          Accept: 'application/json',
        },
      }
    );

    check(res, {
      'status is 200': (r) => r.status === 200,
    });
  });

  sleep(1); // think time (simulasi user)
}
