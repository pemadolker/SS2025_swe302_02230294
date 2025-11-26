import http from "k6/http";
import { check } from "k6";
import { Trend, Counter } from "k6/metrics";

// Custom metrics
const dogFetchTime = new Trend("dog_fetch_duration");
const breedsFetchTime = new Trend("breeds_fetch_duration");
const totalDogsFetched = new Counter("total_dogs_fetched");

export const options = {
  vus: 5,
  duration: "30s",
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:3000";

export default function () {
  // Fetch and measure breeds endpoint
  let start = new Date().getTime();
  let response = http.get(`${BASE_URL}/api/dogs/breeds`);
  let duration = new Date().getTime() - start;
  breedsFetchTime.add(duration);

  // Fetch and measure dog image endpoint
  start = new Date().getTime();
  response = http.get(`${BASE_URL}/api/dogs`);
  duration = new Date().getTime() - start;
  dogFetchTime.add(duration);

  if (response.status === 200) {
    totalDogsFetched.add(1);
  }
}