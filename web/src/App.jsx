import { useCallback, useEffect, useState } from "react";
import { api } from "./api.js";
import "./App.css";

function Status({ message }) {
  if (!message.text) return null;
  return <p className={`status ${message.kind}`}>{message.text}</p>;
}

export default function App() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [status, setStatus] = useState({ text: "", kind: "" });

  // Add-row form state
  const [newName, setNewName] = useState("");
  const [newEmail, setNewEmail] = useState("");
  const [creating, setCreating] = useState(false);

  // Per-row edit state: { [id]: { name, email, saving } }
  const [drafts, setDrafts] = useState({});

  // Lookup-by-id state
  const [lookupId, setLookupId] = useState("");

  const show = (text, kind = "") => setStatus({ text, kind });

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.listUsers();
      setUsers(Array.isArray(data) ? data : []);
      show("", "");
    } catch (err) {
      show(`Could not load users: ${err.message}`, "error");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const startEdit = (user) => {
    setDrafts((d) => ({
      ...d,
      [user.id]: { name: user.name, email: user.email, saving: false },
    }));
  };

  const cancelEdit = (id) => {
    setDrafts((d) => {
      const next = { ...d };
      delete next[id];
      return next;
    });
  };

  const saveEdit = async (id) => {
    const draft = drafts[id];
    if (!draft) return;
    if (!draft.name.trim() || !draft.email.trim()) {
      show("Name and email cannot be empty.", "error");
      return;
    }
    setDrafts((d) => ({ ...d, [id]: { ...draft, saving: true } }));
    try {
      const updated = await api.updateUser(id, draft.name.trim(), draft.email.trim());
      setUsers((rows) => rows.map((u) => (u.id === id ? updated : u)));
      cancelEdit(id);
      show(`User #${id} saved.`, "ok");
    } catch (err) {
      setDrafts((d) => ({ ...d, [id]: { ...draft, saving: false } }));
      show(`Save failed: ${err.message}`, "error");
    }
  };

  const handleCreate = async (e) => {
    e.preventDefault();
    if (!newName.trim() || !newEmail.trim()) {
      show("Enter a name and an email to add a row.", "error");
      return;
    }
    setCreating(true);
    try {
      const created = await api.createUser(newName.trim(), newEmail.trim());
      setUsers((rows) => [...rows, created]);
      setNewName("");
      setNewEmail("");
      show(`User #${created.id} added.`, "ok");
    } catch (err) {
      show(`Add failed: ${err.message}`, "error");
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm(`Delete user #${id}?`)) return;
    try {
      await api.deleteUser(id);
      setUsers((rows) => rows.filter((u) => u.id !== id));
      show(`User #${id} deleted.`, "ok");
    } catch (err) {
      show(`Delete failed: ${err.message}`, "error");
    }
  };

  const handleLookup = async (e) => {
    e.preventDefault();
    const id = lookupId.trim();
    if (!id) {
      refresh();
      return;
    }
    try {
      const user = await api.getUser(id);
      setUsers([user]);
      show(`Showing user #${id}. Clear to see all.`, "ok");
    } catch (err) {
      show(`Lookup failed: ${err.message}`, "error");
    }
  };

  return (
    <div className="page">
      <header className="hero">
        <div>
          <h1>Users sheet</h1>
          <p>
            Spreadsheet-style editor for your Go server — backed by{" "}
            <code>GET /users</code>, <code>POST /users</code>,{" "}
            <code>PUT /users/:id</code>, <code>DELETE /users/:id</code>.
          </p>
        </div>
        <button type="button" className="btn ghost" onClick={refresh}>
          Refresh
        </button>
      </header>

      <Status message={status} />

      <section className="card">
        <h2>Add row</h2>
        <form className="add-row" onSubmit={handleCreate}>
          <input
            placeholder="Name"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
          />
          <input
            placeholder="email@example.com"
            value={newEmail}
            onChange={(e) => setNewEmail(e.target.value)}
          />
          <button type="submit" className="btn primary" disabled={creating}>
            {creating ? "Adding…" : "Insert"}
          </button>
        </form>
      </section>

      <section className="card">
        <div className="row-head">
          <h2>Rows ({users.length})</h2>
          <form className="lookup" onSubmit={handleLookup}>
            <input
              placeholder="Lookup id…"
              value={lookupId}
              onChange={(e) => setLookupId(e.target.value)}
            />
            <button type="submit" className="btn ghost">
              Find
            </button>
            <button
              type="button"
              className="btn ghost"
              onClick={() => {
                setLookupId("");
                refresh();
              }}
            >
              Clear
            </button>
          </form>
        </div>

        {loading ? (
          <p className="muted">Loading…</p>
        ) : users.length === 0 ? (
          <p className="muted">No rows. Add one above.</p>
        ) : (
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>id</th>
                  <th>name</th>
                  <th>email</th>
                  <th>created_at</th>
                  <th>actions</th>
                </tr>
              </thead>
              <tbody>
                {users.map((u) => {
                  const draft = drafts[u.id];
                  const editing = Boolean(draft);
                  return (
                    <tr key={u.id}>
                      <td className="mono">{u.id}</td>
                      <td>
                        {editing ? (
                          <input
                            value={draft.name}
                            onChange={(e) =>
                              setDrafts((d) => ({
                                ...d,
                                [u.id]: { ...draft, name: e.target.value },
                              }))
                            }
                          />
                        ) : (
                          u.name
                        )}
                      </td>
                      <td>
                        {editing ? (
                          <input
                            value={draft.email}
                            onChange={(e) =>
                              setDrafts((d) => ({
                                ...d,
                                [u.id]: { ...draft, email: e.target.value },
                              }))
                            }
                          />
                        ) : (
                          u.email
                        )}
                      </td>
                      <td className="muted small">
                        {u.created_at
                          ? new Date(u.created_at).toLocaleString()
                          : "—"}
                      </td>
                      <td className="actions">
                        {editing ? (
                          <>
                            <button
                              type="button"
                              className="btn primary sm"
                              disabled={draft.saving}
                              onClick={() => saveEdit(u.id)}
                            >
                              {draft.saving ? "Saving…" : "Save"}
                            </button>
                            <button
                              type="button"
                              className="btn ghost sm"
                              onClick={() => cancelEdit(u.id)}
                            >
                              Cancel
                            </button>
                          </>
                        ) : (
                          <>
                            <button
                              type="button"
                              className="btn ghost sm"
                              onClick={() => startEdit(u)}
                            >
                              Edit
                            </button>
                            <button
                              type="button"
                              className="btn danger sm"
                              onClick={() => handleDelete(u.id)}
                            >
                              Delete
                            </button>
                          </>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </section>

      <footer className="muted small">
        Go server must be running on port 8080. Change the target in{" "}
        <code>src/api.js</code> (<code>API_BASE</code>) if yours differs.
      </footer>
    </div>
  );
}
