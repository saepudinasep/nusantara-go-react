import { useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const ICONS = {
  dashboard:
    'M4 13h6a1 1 0 001-1V4a1 1 0 00-1-1H4a1 1 0 00-1 1v8a1 1 0 001 1zm0 8h6a1 1 0 001-1v-4a1 1 0 00-1-1H4a1 1 0 00-1 1v4a1 1 0 001 1zm10 0h6a1 1 0 001-1v-8a1 1 0 00-1-1h-6a1 1 0 00-1 1v8a1 1 0 001 1zm0-18h6a1 1 0 001-1V4a1 1 0 00-1-1h-6a1 1 0 00-1 1v4a1 1 0 001 1z',
  users:
    'M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z',
  book: 'M18 2H6c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h13v-2H6.5a.5.5 0 010-1H19a1 1 0 001-1V4a2 2 0 00-2-2zm-2 12H8V8h8v6z',
  kelas: 'M5 13.18v4L12 21l7-3.82v-4L12 17l-7-3.82zM12 3L1 9l11 6 9-4.91V17h2V9L12 3z',
  laporan:
    'M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 3c1.93 0 3.5 1.57 3.5 3.5S13.93 13 12 13s-3.5-1.57-3.5-3.5S10.07 6 12 6zm7 13H5v-.23c0-.62.28-1.2.76-1.58C7.47 15.82 9.64 15 12 15s4.53.82 6.24 2.19c.48.38.76.97.76 1.58V19z',
  calendar:
    'M19 4h-1V2h-2v2H8V2H6v2H5c-1.11 0-1.99.9-1.99 2L3 20a2 2 0 002 2h14c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 16H5V10h14v10zm0-12H5V6h14v2z',
  check: 'M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z',
  profil:
    'M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z',
}

// Menu per role. `route` diisi kalau sudah tersedia halamannya; kalau tidak,
// item ditampilkan non-aktif dengan label "Segera" (placeholder untuk pengembangan lanjutan).
const MENUS = {
  admin: [
    { name: 'Dashboard', icon: 'dashboard', route: '/admin/dashboard' },
    { name: 'Data Kelas', icon: 'kelas', route: '/admin/kelas' },
    { name: 'Data Petugas', icon: 'users' },
    { name: 'Data Guru', icon: 'users' },
    { name: 'Data Siswa', icon: 'users', route: '/admin/siswa' },
    { name: 'Data SPP', icon: 'laporan', route: '/admin/spp' },
    { name: 'Profil', icon: 'profil', route: '/admin/profile' },
  ],
  petugas: [
    { name: 'Dashboard', icon: 'dashboard', route: '/petugas/dashboard' },
    { name: 'Transaksi Pembayaran', icon: 'check' },
    { name: 'Data Siswa', icon: 'users' },
    { name: 'Profil', icon: 'profil', route: '/petugas/profile' },
  ],
  guru: [
    { name: 'Dashboard', icon: 'dashboard', route: '/guru/dashboard' },
    { name: 'Jadwal Mengajar', icon: 'calendar' },
    { name: 'Materi & Kuis', icon: 'book' },
    { name: 'Profil', icon: 'profil', route: '/guru/profile' },
  ],
  siswa: [
    { name: 'Dashboard', icon: 'dashboard', route: '/siswa/dashboard' },
    { name: 'Tagihan SPP', icon: 'book' },
    { name: 'Riwayat Pembayaran', icon: 'calendar' },
    { name: 'Profil', icon: 'profil', route: '/siswa/profile' },
  ],
}

export default function Sidebar({ isOpen, onClose }) {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  const menus = MENUS[user?.role] || []

  const initials = (user?.username || '?').trim().slice(0, 2).toUpperCase()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <>
      {isOpen && <div className="sidebar-overlay" onClick={onClose} />}

      <aside className={'sidebar' + (isOpen ? ' open' : '')}>
        <div className="sidebar-logo">
          <div className="sidebar-logo-icon">
            <svg viewBox="0 0 24 24">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 14H9V8h2v8zm4 0h-2V8h2v8z" />
            </svg>
          </div>
          <div className="sidebar-logo-text">
            SchoolApp
            <small>Sistem Informasi Sekolah</small>
          </div>
        </div>

        <nav className="sidebar-section">
          {menus.map((menu) => {
            const isRoute = Boolean(menu.route)
            const isActive = isRoute && location.pathname === menu.route

            if (!isRoute) {
              return (
                <button key={menu.name} type="button" className="nav-item disabled" disabled>
                  <svg className="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                    <path d={ICONS[menu.icon]} />
                  </svg>
                  {menu.name}
                  <span className="nav-soon">Segera</span>
                </button>
              )
            }

            return (
              <button
                key={menu.name}
                type="button"
                className={'nav-item' + (isActive ? ' active' : '')}
                onClick={() => {
                  onClose?.()
                  navigate(menu.route)
                }}
              >
                <svg className="nav-icon" viewBox="0 0 24 24" fill="currentColor">
                  <path d={ICONS[menu.icon]} />
                </svg>
                {menu.name}
              </button>
            )
          })}
        </nav>

        {user && (
          <div className="sidebar-footer">
            <div className="user-card">
              <div className="user-avatar">{initials}</div>
              <div>
                <div className="user-name">{user.username}</div>
                <div className="user-role">{user.role}</div>
              </div>
            </div>
            <button type="button" className="logout-btn" onClick={handleLogout}>
              ← Keluar
            </button>
          </div>
        )}
      </aside>
    </>
  )
}
