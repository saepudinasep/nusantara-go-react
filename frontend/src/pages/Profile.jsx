import { useEffect, useState } from 'react'
import axiosClient from '../api/axiosClient'
import { useAuth } from '../context/AuthContext'
import Sidebar from '../components/Sidebar'
import Topbar from '../components/Topbar'

export default function Profile() {
    const { user } = useAuth()
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [profile, setProfile] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')

    useEffect(() => {
        let active = true
        setLoading(true)
        axiosClient
            .get(`/${user.role}/profile`)
            .then((res) => {
                if (active) setProfile(res.data.data)
            })
            .catch((err) => {
                if (active) setError(err.response?.data?.message || 'Gagal memuat profil')
            })
            .finally(() => active && setLoading(false))
        return () => {
            active = false
        }
    }, [user.role])

    const initials = (user?.username || '?').trim().slice(0, 2).toUpperCase()

    // Field yang ditampilkan berbeda per role, sesuai sumber data di backend:
    // admin & guru murni dari tabel users, petugas dari staffs, siswa dari students+classes.
    const renderFields = () => {
        if (!profile) return null

        switch (user.role) {
            case 'admin':
            case 'guru':
                return (
                    <>
                        <ProfileField label="Username" value={profile.username} />
                        <ProfileField label="Role" value={profile.role} capitalize />
                    </>
                )
            case 'petugas':
                return (
                    <>
                        <ProfileField label="Username" value={profile.username} />
                        <ProfileField label="Nama" value={profile.nama || '-'} />
                        <ProfileField label="Posisi / Jabatan" value={profile.posisi || '-'} />
                        <ProfileField label="Role" value={profile.role} capitalize />
                    </>
                )
            case 'siswa':
                return (
                    <>
                        <ProfileField label="Username" value={profile.username} />
                        <ProfileField label="Nama" value={profile.nama || '-'} />
                        <ProfileField label="NISN" value={profile.nisn || '-'} />
                        <ProfileField label="Kelas" value={profile.nama_kelas ? `${profile.nama_kelas} (Tingkat ${profile.tingkat})` : '-'} />
                        <ProfileField label="Alamat" value={profile.alamat || '-'} />
                        <ProfileField label="No. Telepon" value={profile.no_telp || '-'} />
                        <ProfileField label="Role" value={profile.role} capitalize />
                    </>
                )
            default:
                return null
        }
    }

    const showIncompleteNotice =
        profile && (user.role === 'petugas' || user.role === 'siswa') && profile.has_profile === false

    return (
        <div className="app-shell">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

            <div className="main-shell">
                <Topbar title="Profil" onMenuClick={() => setSidebarOpen((prev) => !prev)} />

                <main className="main-content">
                    <div className="breadcrumb">
                        SchoolApp · <span>Profil</span>
                    </div>
                    <div className="page-header">
                        <h1>Profil Saya</h1>
                        <p>Informasi akun kamu di Sistem Informasi Sekolah.</p>
                    </div>

                    <div className="card">
                        {loading ? (
                            <div className="empty-state">Memuat profil...</div>
                        ) : error ? (
                            <div className="empty-state">{error}</div>
                        ) : (
                            <div className="card-body">
                                <div className="profile-header">
                                    <div className="profile-avatar">{initials}</div>
                                    <div>
                                        <div className="profile-name">{profile.nama || profile.username}</div>
                                        <span className="badge">{profile.role}</span>
                                    </div>
                                </div>

                                {showIncompleteNotice && (
                                    <div className="form-error" style={{ marginTop: 20 }}>
                                        Profil detail kamu belum dibuat oleh admin. Beberapa data di bawah masih kosong.
                                    </div>
                                )}

                                <div className="profile-fields">{renderFields()}</div>
                            </div>
                        )}
                    </div>
                </main>
            </div>
        </div>
    )
}

function ProfileField({ label, value, capitalize }) {
    return (
        <div className="profile-field">
            <div className="profile-field-label">{label}</div>
            <div className={'profile-field-value' + (capitalize ? ' capitalize' : '')}>{value}</div>
        </div>
    )
}
