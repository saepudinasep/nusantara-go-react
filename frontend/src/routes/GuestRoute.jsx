import { Navigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const ROLE_REDIRECT = {
  admin: '/admin/dashboard',
  guru: '/guru/dashboard',
  siswa: '/siswa/dashboard',
}

// GuestRoute membungkus halaman yang HANYA boleh diakses jika belum login (misal: /login).
// Jika user sudah login, otomatis di-redirect ke dashboard sesuai role-nya.
export default function GuestRoute({ children }) {
  const { user } = useAuth()

  if (user) {
    return <Navigate to={ROLE_REDIRECT[user.role] || '/login'} replace />
  }

  return children
}
