import http from "k6/http";
import { check, sleep } from "k6";
import { Counter } from "k6/metrics";

const totalRequests = new Counter("total_requests");

export const options = {
  scenarios: {
    light_load: {
      executor: "constant-vus",
      vus: 10,
      duration: "1m",
    },
    spike_test: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "10s", target: 50 }, // Spike to 50 users
        { duration: "30s", target: 50 }, // Stay at 50
        { duration: "10s", target: 0 }, // Drop back
      ],
      startTime: "1m",
    },
  },
  thresholds: {
    http_req_duration: ["p(90)<1000"], // 90% under 1s
    http_req_failed: ["rate<0.05"], // Less than 5% errors
  },
};