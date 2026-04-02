"use client";

import { useState } from "react";

export default function UserImportPanel() {
  const [fileAssetId, setFileAssetId] = useState("");
  const [email, setEmail] = useState("");
  const [fullName, setFullName] = useState("");
  const [message, setMessage] = useState("");

  const importUsers = async () => {
    setMessage("");
    const response = await fetch("http://localhost:8080/api/v1/tenants/tenant-001/users/import-jobs", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": "tenant-001",
      },
      body: JSON.stringify({
        fileAssetId,
        defaultRole: "student",
        sendInvitation: true,
      }),
    });

    if (!response.ok) {
      const json = (await response.json()) as { error?: string };
      setMessage(json.error ?? "Import request failed");
      return;
    }
    setMessage("CSV import job queued");
  };

  const inviteUser = async () => {
    setMessage("");
    const response = await fetch("http://localhost:8080/api/v1/tenants/tenant-001/users", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": "tenant-001",
      },
      body: JSON.stringify({
        email,
        fullName,
        role: "student",
        sendInvitation: true,
      }),
    });

    if (!response.ok) {
      const json = (await response.json()) as { error?: string };
      setMessage(json.error ?? "Invitation failed");
      return;
    }
    setMessage("User created and invitation queued");
  };

  return (
    <section>
      <h2>CSV Import</h2>
      <input
        value={fileAssetId}
        onChange={(event) => setFileAssetId(event.target.value)}
        placeholder="File asset ID"
      />
      <button onClick={importUsers}>Queue CSV Import</button>

      <h2>Invite User</h2>
      <input value={email} onChange={(event) => setEmail(event.target.value)} placeholder="Email" />
      <input value={fullName} onChange={(event) => setFullName(event.target.value)} placeholder="Full name" />
      <button onClick={inviteUser}>Create + Invite</button>

      {message ? <p>{message}</p> : null}
    </section>
  );
}
