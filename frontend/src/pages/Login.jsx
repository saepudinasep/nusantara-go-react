import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const ROLE_REDIRECT = {
  admin: '/admin/dashboard',
  petugas: '/petugas/dashboard',
  guru: '/guru/dashboard',
  siswa: '/siswa/dashboard',
}

export default function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const { login, loading } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    const result = await login(username, password)
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
            Pembayaran SPP
          </h1>
          <p>
            Platform digital untuk mengelola pembayaran SPP, data kelas, dan
            data siswa secara terpusat untuk admin, petugas, guru, dan siswa.
          </p>
        </div>

        <div className="login-stats">
          <div className="stat-item">
            <div className="stat-num">4</div>
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
              <label htmlFor="username">Username</label>
              <input
                id="username"
                type="text"
                placeholder="Masukkan username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
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
            Admin: admin
            <br />
            Petugas: petugas1
            <br />
            Guru: guru1
            <br />
            Siswa: siswa1
          </div>
        </div>
      </div>
    </div>
  )
}
