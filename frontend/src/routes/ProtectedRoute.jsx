import { Navigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

// ProtectedRoute membungkus halaman yang butuh login.
// Jika allowedRoles diisi, hanya role yang tercantum yang boleh mengakses.
export default function ProtectedRoute({ children, allowedRoles }) {
  const { user } = useAuth()

  if (!user) {
    return <Navigate to="/login" replace />
  }

  if (allowedRoles && !allowedRoles.includes(user.role)) {
    return <Navigate to="/unauthorized" replace />
  }

  return children
}
