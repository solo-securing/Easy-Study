type Tenant = {
  id: string;
  name: string;
  subdomain: string;
  status: "active" | "suspended" | "deactivated";
  planCode: string;
};

async function fetchTenants(): Promise<Tenant[]> {
  const response = await fetch("http://localhost:8080/api/v1/super-admin/tenants", {
    headers: {
      "X-Tenant-ID": "platform",
    },
    cache: "no-store",
  });

  if (!response.ok) {
    return [];
  }

  const json = (await response.json()) as { items?: Tenant[] };
  return json.items ?? [];
}

export default async function SuperAdminTenantsPage() {
  const tenants = await fetchTenants();

  return (
    <main>
      <h1>Tenant Management</h1>
      {tenants.length === 0 ? (
        <p>No tenants found.</p>
      ) : (
        <ul>
          {tenants.map((tenant) => (
            <li key={tenant.id}>
              <strong>{tenant.name}</strong> ({tenant.subdomain}) - {tenant.status} [{tenant.planCode}]
            </li>
          ))}
        </ul>
      )}
    </main>
  );
}
