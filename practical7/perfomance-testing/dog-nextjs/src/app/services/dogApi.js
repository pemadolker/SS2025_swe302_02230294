// app/services/dogApi.js
import axios from "axios";

const BASE_URL = "https://dog.ceo/api";

export async function getRandomDog() {
  const res = await axios.get(`${BASE_URL}/breeds/image/random`);
  return res.data;
}

export async function getBreeds() {
  const res = await axios.get(`${BASE_URL}/breeds/list/all`);
  return res.data;
}

export async function getBreedImage(breed) {
  const res = await axios.get(`${BASE_URL}/breed/${breed}/images/random`);
  return res.data;
}
