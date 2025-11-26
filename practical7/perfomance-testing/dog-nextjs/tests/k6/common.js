// tests/k6/common.js
export const BASE_URL = __ENV.BASE_URL || "https://ec7168d5843e.ngrok-free.app";

export function safeJson(body) {
  try {
    return JSON.parse(body);
  } catch (e) {
    return null;
  }
}
