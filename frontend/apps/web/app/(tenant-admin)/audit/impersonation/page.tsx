type ImpersonationAuditRow = {
  sessionId: string;
  superAdminId: string;
  startedAt: string;
  endedAt: string;
  actions: string[];
};

async function fetchAuditRows(): Promise<ImpersonationAuditRow[]> {
  const response = await fetch("http://localhost:8080/api/v1/tenants/tenant-001/impersonation-audit", {
    headers: {
      "X-Tenant-ID": "tenant-001",
    },
    cache: "no-store",
  });
  if (!response.ok) {
    return [];
  }
  const json = (await response.json()) as { items?: ImpersonationAuditRow[] };
  return json.items ?? [];
}

export default async function TenantImpersonationAuditPage() {
  const rows = await fetchAuditRows();

  return (
    <main>
      <h1>Impersonation Audit</h1>
      {rows.length === 0 ? (
        <p>No impersonation sessions found.</p>
      ) : (
        <ul>
          {rows.map((row) => (
            <li key={row.sessionId}>
              <strong>{row.sessionId}</strong> by {row.superAdminId} | {row.startedAt} - {row.endedAt}
              <div>Actions: {row.actions.join(", ") || "none"}</div>
            </li>
          ))}
        </ul>
      )}
    </main>
  );
}
