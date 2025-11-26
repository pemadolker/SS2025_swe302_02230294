// tests/k6/spike-load.js
import http from "k6/http";
import { check, sleep } from "k6";
import { BASE_URL } from "./common.js";

export const options = {
  scenarios: {
    baseline: {
      executor: "constant-vus",
      vus: 5,
      duration: "30s",
    },
    spike: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "10s", target: 0 },
        { duration: "10s", target: 95 }, // spike (adjust if laptop can't handle)
        { duration: "30s", target: 95 },
        { duration: "10s", target: 0 },
      ],
      startTime: "30s",
    },
  },
  thresholds: {
    http_req_duration: ["p(95)<1500"],
    http_req_failed: ["rate<0.1"],
  },
};

export default function () {
  const r1 = http.get(`${BASE_URL}/`);
  check(r1, { "home 200": (r) => r.status === 200 });
  sleep(1);

  const r2 = http.get("https://dog.ceo/api/breeds/list/all");
  check(r2, { "breeds 200": (r) => r.status === 200 });
  sleep(0.5);

  const r3 = http.get("https://dog.ceo/api/breeds/image/random");
  check(r3, { "random 200": (r) => r.status === 200 });
  sleep(1);
}
