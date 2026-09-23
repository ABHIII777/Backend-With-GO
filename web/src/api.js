// Frontend-only helpers for the custom Go REST API.
// No backend logic here — just fetch() wrappers.
//
// Configure this one constant to point at your Go server.
// Dev default: Go listens on :8080 (see RESTapi/cmd/server/main.go).
export const API_BASE = "http://localhost:8080";

async function parseBody(res) {
  const text = await res.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

async function request(path, options = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  const data = await parseBody(res);
  if (!res.ok) {
    const msg =
      data && typeof data === "object" && data.error
        ? data.error
        : `Request failed (${res.status})`;
    throw new Error(msg);
  }
  return data;
}

export const api = {
  // GET /users -> [{id, name, email, created_at}]
  listUsers: () => request("/users"),
  // GET /users/:id -> {id, name, email, created_at}
  getUser: (id) => request(`/users/${id}`),
  // POST /users {name, email} -> 201 user
  createUser: (name, email) =>
    request("/users", {
      method: "POST",
      body: JSON.stringify({ name, email }),
    }),
  // PUT /users/:id {name, email} -> 200 user (full row save)
  updateUser: (id, name, email) =>
    request(`/users/${id}`, {
      method: "PUT",
      body: JSON.stringify({ name, email }),
    }),
  // DELETE /users/:id -> 204
  deleteUser: (id) => request(`/users/${id}`, { method: "DELETE" }),
};
