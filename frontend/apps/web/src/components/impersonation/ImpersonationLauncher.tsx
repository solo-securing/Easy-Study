"use client";

import { FormEvent, useState } from "react";

export default function ImpersonationLauncher() {
  const [tenantId, setTenantId] = useState("tenant-001");
  const [reason, setReason] = useState("Support request");
  const [sessionId, setSessionId] = useState("");
  const [message, setMessage] = useState("");

  const startImpersonation = async (event: FormEvent) => {
    event.preventDefault();
    setMessage("");

    const response = await fetch("http://localhost:8080/api/v1/admin/impersonations", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": "platform",
      },
      body: JSON.stringify({ tenantId, reason }),
    });
    if (!response.ok) {
      setMessage("Failed to start impersonation session");
      return;
    }
    const json = (await response.json()) as { sessionId?: string };
    setSessionId(json.sessionId ?? "");
    setMessage("Impersonation session started");
  };

  const endImpersonation = async () => {
    if (!sessionId) {
      setMessage("No active session");
      return;
    }
    const response = await fetch(`http://localhost:8080/api/v1/admin/impersonations/${sessionId}`, {
      method: "DELETE",
      headers: {
        "X-Tenant-ID": "platform",
      },
    });
    if (!response.ok) {
      setMessage("Failed to end impersonation session");
      return;
    }
    setMessage("Impersonation session ended");
    setSessionId("");
  };

  return (
    <section>
      <h2>Impersonation Launcher</h2>
      <form onSubmit={startImpersonation}>
        <label>
          Tenant ID
          <input value={tenantId} onChange={(event) => setTenantId(event.target.value)} />
        </label>
        <label>
          Reason
          <input value={reason} onChange={(event) => setReason(event.target.value)} />
        </label>
        <button type="submit">Start Impersonation</button>
      </form>
      <button type="button" onClick={endImpersonation} disabled={!sessionId}>
        End Impersonation
      </button>
      {sessionId ? <p>Session ID: {sessionId}</p> : null}
      {message ? <p>{message}</p> : null}
    </section>
  );
}
