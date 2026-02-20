import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 100 },
    { duration: '15s', target: 100 },
    { duration: '5s', target: 500 },
  ],
};

export default function () {
  const url = 'http://localhost:9000/api/url';
  const payload = JSON.stringify({ name: 'John' + new Date().getTime(), url: "https://google.com", idempotencyKey: new Date().getTime() });
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
    // 'response has expected data': (r) => r.json().name === 'John',
  });
}
