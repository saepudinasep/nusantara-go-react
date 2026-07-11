import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import ProtectedRoute from './routes/ProtectedRoute'
import GuestRoute from './routes/GuestRoute'

import Login from './pages/Login'
import DashboardAdmin from './pages/DashboardAdmin'
import DashboardPetugas from './pages/DashboardPetugas'
import DashboardGuru from './pages/DashboardGuru'
import DashboardSiswa from './pages/DashboardSiswa'
import KelasManagement from './pages/KelasManagement'
import SppManagement from './pages/SppManagement'
import SiswaManagement from './pages/SiswaManagement'
import PetugasManagement from './pages/PetugasManagement'
import PetugasKelas from './pages/PetugasKelas'
import Profile from './pages/Profile'
import { NotFound, Unauthorized } from './pages/NotFound'

function RootRedirect() {
  const { user } = useAuth()
  if (!user) return <Navigate to="/login" replace />

  const map = {
    admin: '/admin/dashboard',
    petugas: '/petugas/dashboard',
    guru: '/guru/dashboard',
    siswa: '/siswa/dashboard',
  }
  return <Navigate to={map[user.role] || '/login'} replace />
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<RootRedirect />} />
          <Route
            path="/login"
            element={
              <GuestRoute>
                <Login />
              </GuestRoute>
            }
          />
          <Route path="/unauthorized" element={<Unauthorized />} />

          <Route
            path="/admin/dashboard"
            element={
              <ProtectedRoute allowedRoles={['admin']}>
                <DashboardAdmin />
              </ProtectedRoute>
            }
          />

          <Route
            path="/admin/kelas"
            element={
              <ProtectedRoute allowedRoles={['admin']}>
                <KelasManagement />
              </ProtectedRoute>
            }
          />

          <Route
            path="/admin/spp"
            element={
              <ProtectedRoute allowedRoles={['admin']}>
                <SppManagement />
              </ProtectedRoute>
            }
          />

          <Route
            path="/admin/siswa"
            element={
              <ProtectedRoute allowedRoles={['admin']}>
                <SiswaManagement />
              </ProtectedRoute>
            }
          />

          <Route
            path="/admin/petugas"
            element={
              <ProtectedRoute allowedRoles={['admin']}>
                <PetugasManagement />
              </ProtectedRoute>
            }
          />

          <Route
            path="/admin/profile"
            element={
              <ProtectedRoute allowedRoles={['admin']}>
                <Profile />
              </ProtectedRoute>
            }
          />

          <Route
            path="/petugas/profile"
            element={
              <ProtectedRoute allowedRoles={['petugas']}>
                <Profile />
              </ProtectedRoute>
            }
          />

          <Route
            path="/guru/profile"
            element={
              <ProtectedRoute allowedRoles={['guru']}>
                <Profile />
              </ProtectedRoute>
            }
          />

          <Route
            path="/siswa/profile"
            element={
              <ProtectedRoute allowedRoles={['siswa']}>
                <Profile />
              </ProtectedRoute>
            }
          />

          <Route
            path="/petugas/dashboard"
            element={
              <ProtectedRoute allowedRoles={['petugas']}>
                <DashboardPetugas />
              </ProtectedRoute>
            }
          />

          <Route
            path="/petugas/kelas"
            element={
              <ProtectedRoute allowedRoles={['petugas']}>
                <PetugasKelas />
              </ProtectedRoute>
            }
          />

          <Route
            path="/guru/dashboard"
            element={
              <ProtectedRoute allowedRoles={['guru']}>
                <DashboardGuru />
              </ProtectedRoute>
            }
          />

          <Route
            path="/siswa/dashboard"
            element={
              <ProtectedRoute allowedRoles={['siswa']}>
                <DashboardSiswa />
              </ProtectedRoute>
            }
          />

          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}
