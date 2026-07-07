import { useEffect, useState } from 'react'
import { useAuth } from '../context/AuthContext'
import Sidebar from './Sidebar'
import Topbar from './Topbar'

const STAT_ICONS = {
  users: 'M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z',
  book: 'M18 2H6c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h13v-2H6.5a.5.5 0 010-1H19a1 1 0 001-1V4a2 2 0 00-2-2zm-2 12H8V8h8v6z',
  kelas: 'M5 13.18v4L12 21l7-3.82v-4L12 17l-7-3.82zM12 3L1 9l11 6 9-4.91V17h2V9L12 3z',
  calendar:
    'M19 4h-1V2h-2v2H8V2H6v2H5c-1.11 0-1.99.9-1.99 2L3 20a2 2 0 002 2h14c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 16H5V10h14v10zm0-12H5V6h14v2z',
  check: 'M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z',
  dashboard:
    'M4 13h6a1 1 0 001-1V4a1 1 0 00-1-1H4a1 1 0 00-1 1v8a1 1 0 001 1zm0 8h6a1 1 0 001-1v-4a1 1 0 00-1-1H4a1 1 0 00-1 1v4a1 1 0 001 1zm10 0h6a1 1 0 001-1v-8a1 1 0 00-1-1h-6a1 1 0 00-1 1v8a1 1 0 001 1zm0-18h6a1 1 0 001-1V4a1 1 0 00-1-1h-6a1 1 0 00-1 1v4a1 1 0 001 1z',
}

export default function DashboardLayout({ title, subtitle, fetchData }) {
  const { user } = useAuth()
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [stats, setStats] = useState([])
  const [activities, setActivities] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Tutup sidebar otomatis kalau layar dibesarkan kembali ke ukuran desktop
    const handleResize = () => {
      if (window.innerWidth > 1024) setSidebarOpen(false)
    }
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  useEffect(() => {
    let active = true
    setLoading(true)
    fetchData()
      .then((data) => {
        if (!active) return
        setStats(data.stats || [])
        setActivities(data.activities || [])
      })
      .finally(() => active && setLoading(false))
    return () => {
      active = false
    }
  }, [fetchData])

  return (
    <div className="app-shell">
      <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

      <div className="main-shell">
        <Topbar title={title} onMenuClick={() => setSidebarOpen((prev) => !prev)} />

        <main className="main-content">
          <div className="breadcrumb">
            SchoolApp · <span>{title}</span>
          </div>
          <div className="page-header">
            <h1>{title}</h1>
            <p>{subtitle || `Selamat datang kembali, ${user?.name} 👋`}</p>
          </div>

          {loading ? (
            <div className="empty-state">Memuat data...</div>
          ) : (
            <>
              <div className="stats-grid">
                {stats.map((s) => (
                  <div className={`stat-card ${s.color}`} key={s.label}>
                    <div className="stat-card-top">
                      <div className="stat-card-label">{s.label}</div>
                      <div className="stat-card-icon">
                        <svg viewBox="0 0 24 24">
                          <path d={STAT_ICONS[s.icon] || STAT_ICONS.dashboard} />
                        </svg>
                      </div>
                    </div>
                    <div>
                      <div className="stat-card-val">{s.value}</div>
                      <div className="stat-card-sub">{s.sub}</div>
                    </div>
                  </div>
                ))}
              </div>

              <div className="two-col">
                <div className="card">
                  <div className="card-header">
                    <div>
                      <div className="card-title">Ringkasan</div>
                      <div className="card-subtitle">Ikhtisar aktivitas terbaru di akun kamu</div>
                    </div>
                  </div>
                  <div style={{ padding: '22px' }}>
                    <p style={{ fontSize: 13, color: 'var(--text2)', lineHeight: 1.7 }}>
                      Ini adalah dashboard contoh untuk role <strong style={{ textTransform: 'capitalize' }}>{user?.role}</strong>.
                      Kartu di atas dan panel aktivitas di samping bisa dihubungkan ke data asli
                      dari database sesuai kebutuhan (misalnya jumlah siswa, rata-rata nilai, atau presensi).
                    </p>
                  </div>
                </div>

                <div className="card">
                  <div className="card-header">
                    <div className="card-title">Aktivitas Terkini</div>
                  </div>
                  {activities.length === 0 ? (
                    <div className="empty-state">Belum ada aktivitas.</div>
                  ) : (
                    <div className="timeline">
                      {activities.map((a, i) => (
                        <div className="tl-item" key={i}>
                          <div className="tl-dot" />
                          <div className="tl-content">
                            <div className="tl-label">{a.label}</div>
                            <div className="tl-sub">{a.sub}</div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </>
          )}
        </main>
      </div>
    </div>
  )
}
