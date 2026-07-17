import ProtectedRoute from './ProtectedRoute';

function AdminRoute({ children }: { children: React.ReactNode }) {
  return <ProtectedRoute requireAdmin>{children}</ProtectedRoute>;
}

export default AdminRoute;
