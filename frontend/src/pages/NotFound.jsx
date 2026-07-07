import { Link } from 'react-router-dom'

export function NotFound() {
  return (
    <div className="notfound-screen">
      <div style={{ textAlign: 'center' }}>
        <h2>404 - Halaman tidak ditemukan</h2>
        <p><Link to="/login">Kembali ke Login</Link></p>
      </div>
    </div>
  )
}

export function Unauthorized() {
  return (
    <div className="notfound-screen">
      <div style={{ textAlign: 'center' }}>
        <h2>403 - Anda tidak memiliki akses ke halaman ini</h2>
        <p><Link to="/login">Kembali ke Login</Link></p>
      </div>
    </div>
  )
}
