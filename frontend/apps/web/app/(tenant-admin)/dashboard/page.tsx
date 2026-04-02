type DashboardResponse = {
  tenantId: string;
  activeUserCount: number;
  activeCourseCount: number;
  completionRate: number;
};

async function fetchDashboard(): Promise<DashboardResponse | null> {
  const response = await fetch("http://localhost:8080/api/v1/tenants/tenant-001/reports/dashboard", {
    headers: {
      "X-Tenant-ID": "tenant-001",
    },
    cache: "no-store",
  });
  if (!response.ok) {
    return null;
  }
  return (await response.json()) as DashboardResponse;
}

export default async function TenantAdminDashboardPage() {
  const dashboard = await fetchDashboard();

  return (
    <main>
      <h1>Tenant Dashboard</h1>
      {!dashboard ? (
        <p>No dashboard data.</p>
      ) : (
        <ul>
          <li>Tenant: {dashboard.tenantId}</li>
          <li>Active users: {dashboard.activeUserCount}</li>
          <li>Active courses: {dashboard.activeCourseCount}</li>
          <li>Completion rate: {dashboard.completionRate}%</li>
        </ul>
      )}
    </main>
  );
}
