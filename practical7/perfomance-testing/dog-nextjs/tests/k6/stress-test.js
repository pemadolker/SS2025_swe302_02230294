// tests/k6/stress-test.js
import http from "k6/http";
import { check, sleep } from "k6";
import { BASE_URL } from "./common.js";

export const options = {
  stages: [
    { duration: "1m", target: 10 },
    { duration: "2m", target: 50 },
    { duration: "3m", target: 99 }, // push hard (adjust local)
    { duration: "2m", target: 0 },
  ],
  thresholds: {
    http_req_failed: ["rate<0.2"],
    http_req_duration: ["p(95)<2000"],
  },
};

export default function () {
  const r1 = http.get("https://dog.ceo/api/breeds/image/random");
  check(r1, { "random 200": (r) => r.status === 200 });
  sleep(1);

  const r2 = http.get("https://dog.ceo/api/breeds/list/all");
  check(r2, { "breeds 200": (r) => r.status === 200 });
  sleep(0.5);

  const breed = "husky";
  const r3 = http.get(`https://dog.ceo/api/breed/${breed}/images/random`);
  check(r3, { "breed 200": (r) => r.status === 200 });
  sleep(1.5);
}
