import UserImportPanel from "../../../src/components/users/UserImportPanel";

export default function TenantAdminUsersPage() {
  return (
    <main>
      <h1>Tenant Users</h1>
      <p>Manage users, invitations, and CSV imports.</p>
      <UserImportPanel />
    </main>
  );
}
