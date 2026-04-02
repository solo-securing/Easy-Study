"use client";

import { useState } from "react";

type TenantStatus = "active" | "suspended" | "deactivated";

type Props = {
  tenantId: string;
  currentStatus: TenantStatus;
};

export default function TenantStatusActions({ tenantId, currentStatus }: Props) {
  const [status, setStatus] = useState<TenantStatus>(currentStatus);
  const [message, setMessage] = useState("");

  const updateStatus = async (nextStatus: TenantStatus) => {
    setMessage("");
    const response = await fetch(`http://localhost:8080/api/v1/super-admin/tenants/${tenantId}/status`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": "platform",
      },
      body: JSON.stringify({ status: nextStatus }),
    });

    if (!response.ok) {
      const error = (await response.json()) as { error?: string };
      setMessage(error.error ?? "Update failed");
      return;
    }

    setStatus(nextStatus);
    setMessage("Status updated");
  };

  return (
    <div>
      <p>Current status: {status}</p>
      <button onClick={() => updateStatus("active")} disabled={status === "active"}>
        Activate
      </button>
      <button onClick={() => updateStatus("suspended")} disabled={status === "suspended"}>
        Suspend
      </button>
      <button onClick={() => updateStatus("deactivated")} disabled={status === "deactivated"}>
        Deactivate
      </button>
      {message ? <p>{message}</p> : null}
    </div>
  );
}
