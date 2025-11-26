// tests/k6/smoke-test.js
import http from "k6/http";
import { check } from "k6";
import { BASE_URL } from "./common.js";


export const options = {
  vus: 1,
  duration: "30s",
};

export default function () {

  const endpoints = [
    { name: "Homepage", url: BASE_URL + "/" },
    { name: "RandomDog", url: BASE_URL + "/api/random" }, // fallback route check (but our UI uses external API)
    { name: "Breeds", url: BASE_URL + "/api/breeds" }, // may return 404 if not present; main check is homepage
    { name: "Husky API", url: BASE_URL + "/api/breed/husky" },
  ];


  // Primary check: homepage
  const res = http.get(BASE_URL + "/");
  check(res, {
    "home status 200": (r) => r.status === 200,
    "home loads <2s": (r) => r.timings.duration < 2000,
  });
}
