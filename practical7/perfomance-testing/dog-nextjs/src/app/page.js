// app/page.js
"use client";

import { useEffect, useState } from "react";
import { getRandomDog, getBreeds, getBreedImage } from "./services/dogApi";

export default function Home() {
  const [randomImage, setRandomImage] = useState("");
  const [breeds, setBreeds] = useState([]);
  const [breedImage, setBreedImage] = useState("");
  const [selectedBreed, setSelectedBreed] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    loadRandomDog();
    loadBreeds();
  }, []);

  async function loadRandomDog() {
    try {
      setLoading(true);
      setError("");
      const res = await getRandomDog();
      setRandomImage(res.message);
    } catch (e) {
      setError("Failed to load random dog.");
    } finally {
      setLoading(false);
    }
  }

  async function loadBreeds() {
    try {
      const res = await getBreeds();
      const list = Object.keys(res.message || {});
      setBreeds(list);
    } catch (e) {
      console.error(e);
    }
  }

  async function handleBreedSelect(e) {
    const breed = e.target.value;
    setSelectedBreed(breed);
    if (!breed) {
      setBreedImage("");
      return;
    }
    try {
      setLoading(true);
      const res = await getBreedImage(breed);
      setBreedImage(res.message);
    } catch (e) {
      setError("Failed to load breed image.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div style={{ padding: 20, fontFamily: "Arial, sans-serif" }}>
      <h1>Dog CEO — Demo App</h1>

      <section style={{ marginBottom: 20 }}>
        <h2>Random Dog</h2>
        <button onClick={loadRandomDog} disabled={loading}>
          {loading ? "Loading..." : "Get Random Dog"}
        </button>
        {randomImage && (
          <div style={{ marginTop: 10 }}>
            <img src={randomImage} width="350" alt="random dog" />
          </div>
        )}
      </section>

      <section style={{ marginBottom: 20 }}>
        <h2>Select Breed</h2>
        <select value={selectedBreed} onChange={handleBreedSelect}>
          <option value="">-- choose a breed --</option>
          {breeds.map((b) => (
            <option key={b} value={b}>
              {b}
            </option>
          ))}
        </select>

        {breedImage && (
          <div style={{ marginTop: 10 }}>
            <h3>{selectedBreed}</h3>
            <img src={breedImage} width="350" alt={`${selectedBreed} dog`} />
          </div>
        )}
      </section>

      {error && <div style={{ color: "red" }}>{error}</div>}
    </div>
  );
}
