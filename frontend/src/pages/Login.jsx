import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const ROLE_REDIRECT = {
  admin: '/admin/dashboard',
  guru: '/guru/dashboard',
  siswa: '/siswa/dashboard',
}

export default function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const { login, loading } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    const result = await login(email, password)
    if (!result.success) {
      setError(result.message)
      return
    }

    navigate(ROLE_REDIRECT[result.user.role] || '/login')
  }

  return (
    <div className="login-wrap">
      <div className="login-deco" />

      {/* Kiri — branding & hero, disembunyikan di layar sempit */}
      <div className="login-left">
        <div className="login-logo">
          <div className="login-logo-icon">
            <svg viewBox="0 0 24 24">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 14H9V8h2v8zm4 0h-2V8h2v8z" />
            </svg>
          </div>
          <div className="login-logo-name">
            School<span>App</span>
          </div>
        </div>

        <div className="login-hero">
          <h1>
            Sistem Informasi
            <br />
            Sekolah Terpadu
          </h1>
          <p>
            Platform digital untuk mengelola aktivitas belajar-mengajar,
            presensi, dan nilai secara terpusat untuk admin, guru, dan siswa.
          </p>
        </div>

        <div className="login-stats">
          <div className="stat-item">
            <div className="stat-num">3</div>
            <div className="stat-label">Level Akses</div>
          </div>
          <div className="stat-item">
            <div className="stat-num">100%</div>
            <div className="stat-label">Tercatat Digital</div>
          </div>
          <div className="stat-item">
            <div className="stat-num">Real-time</div>
            <div className="stat-label">Data Terbaru</div>
          </div>
        </div>
      </div>

      {/* Kanan — form login */}
      <div className="login-right">
        <div className="login-card">
          <h2>Selamat Datang</h2>
          <p className="login-subtitle">Masuk ke Sistem Informasi Sekolah</p>

          {error && <div className="login-error">{error}</div>}

          <form onSubmit={handleSubmit}>
            <div className="login-field">
              <label htmlFor="email">Email</label>
              <input
                id="email"
                type="email"
                placeholder="nama@sekolah.id"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                autoComplete="username"
                autoFocus
                required
              />
            </div>

            <div className="login-field">
              <label htmlFor="password">Password</label>
              <input
                id="password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                autoComplete="current-password"
                required
              />
            </div>

            <button type="submit" className="login-btn" disabled={loading}>
              {loading && <span className="login-spinner" />}
              {loading ? 'Memproses...' : 'Masuk ke Sistem →'}
            </button>
          </form>

          <div className="login-hint">
            <strong>Akun contoh (password: password123)</strong>
            Admin: admin@sekolah.com
            <br />
            Guru: guru@sekolah.com
            <br />
            Siswa: siswa@sekolah.com
          </div>
        </div>
      </div>
    </div>
  )
}
