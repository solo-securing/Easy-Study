"use client";

import { FormEvent, useState } from "react";

type FormState = {
  name: string;
  subdomain: string;
  ownerEmail: string;
  planCode: "free" | "trial" | "pro" | "enterprise";
};

const initialState: FormState = {
  name: "",
  subdomain: "",
  ownerEmail: "",
  planCode: "trial",
};

export default function TenantCreateForm() {
  const [state, setState] = useState<FormState>(initialState);
  const [message, setMessage] = useState<string>("");
  const [submitting, setSubmitting] = useState(false);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setSubmitting(true);
    setMessage("");

    try {
      const response = await fetch("http://localhost:8080/api/v1/super-admin/tenants", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Tenant-ID": "platform",
        },
        body: JSON.stringify(state),
      });

      if (!response.ok) {
        const error = (await response.json()) as { error?: string };
        setMessage(error.error ?? "Failed to create tenant");
        return;
      }

      setState(initialState);
      setMessage("Tenant created successfully");
    } catch {
      setMessage("Failed to create tenant");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={submit}>
      <h2>Create Tenant</h2>
      <input
        placeholder="Tenant name"
        value={state.name}
        onChange={(e) => setState((prev) => ({ ...prev, name: e.target.value }))}
      />
      <input
        placeholder="Subdomain"
        value={state.subdomain}
        onChange={(e) => setState((prev) => ({ ...prev, subdomain: e.target.value }))}
      />
      <input
        placeholder="Owner email"
        value={state.ownerEmail}
        onChange={(e) => setState((prev) => ({ ...prev, ownerEmail: e.target.value }))}
      />
      <select
        value={state.planCode}
        onChange={(e) =>
          setState((prev) => ({ ...prev, planCode: e.target.value as FormState["planCode"] }))
        }
      >
        <option value="free">free</option>
        <option value="trial">trial</option>
        <option value="pro">pro</option>
        <option value="enterprise">enterprise</option>
      </select>
      <button disabled={submitting} type="submit">
        {submitting ? "Creating..." : "Create tenant"}
      </button>
      {message ? <p>{message}</p> : null}
    </form>
  );
}
