// tests/k6/soak-test.js
import http from "k6/http";
import { check, sleep } from "k6";
import { Trend } from "k6/metrics";
import { BASE_URL } from "./common.js";

const dogFetch = new Trend("dog_fetch_duration");

export const options = {
  stages: [
    { duration: "2m", target: 20 },
    { duration: "5m", target: 20 },
    { duration: "2m", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<1500"],
    http_req_failed: ["rate<0.05"],
    dog_fetch_duration: ["p(95)<1500"],
  },
};

export default function () {
  const start = Date.now();
  const r = http.get("https://dog.ceo/api/breeds/image/random");
  dogFetch.add(Date.now() - start);
  check(r, { "random 200": (res) => res.status === 200 });
  sleep(2);

  const r2 = http.get("https://dog.ceo/api/breeds/list/all");
  check(r2, { "breeds 200": (res) => res.status === 200 });
  sleep(3);
}
