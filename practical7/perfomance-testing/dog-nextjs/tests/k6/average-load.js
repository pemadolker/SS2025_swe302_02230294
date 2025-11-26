// tests/k6/average-load.js
import http from "k6/http";
import { check, sleep } from "k6";
import { BASE_URL, safeJson } from "./common.js";

export const options = {
  stages: [
    { duration: "30s", target: 10 },
    { duration: "2m", target: 10 },
    { duration: "30s", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(95)<800"],
    http_req_failed: ["rate<0.05"],
  },
};

export default function () {
  const resHome = http.get(`${BASE_URL}/`);
  check(resHome, {
    "home 200": (r) => r.status === 200,
  });
  sleep(2);

  // Simulate client fetching breed list from external Dog API via UI: call Dog CEO directly
  const resBreeds = http.get("https://dog.ceo/api/breeds/list/all");
  check(resBreeds, { "breeds 200": (r) => r.status === 200 });
  sleep(1);

  const resRandom = http.get("https://dog.ceo/api/breeds/image/random");
  check(resRandom, { "random 200": (r) => r.status === 200 });
  sleep(1);
}
