import { useCallback, useEffect, useState } from 'react'
import Swal from 'sweetalert2'
import axiosClient from '../api/axiosClient'
import Sidebar from '../components/Sidebar'
import Topbar from '../components/Topbar'

const EMPTY_FORM = { username: '', password: '', nisn: '', nama: '', class_id: '', alamat: '', no_telp: '' }
const PAGE_SIZE = 10
const DEFAULT_PAGINATION = { page: 1, limit: PAGE_SIZE, total: 0, total_pages: 1 }

const Toast = Swal.mixin({
    toast: true,
    position: 'top-end',
    showConfirmButton: false,
    timer: 2500,
    timerProgressBar: true,
})

export default function SiswaManagement() {
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [siswaList, setSiswaList] = useState([])
    const [kelasOptions, setKelasOptions] = useState([])
    const [loading, setLoading] = useState(true)

    const [page, setPage] = useState(1)
    const [pagination, setPagination] = useState(DEFAULT_PAGINATION)

    const [modalOpen, setModalOpen] = useState(false)
    const [editingId, setEditingId] = useState(null)
    const [form, setForm] = useState(EMPTY_FORM)
    const [formError, setFormError] = useState('')
    const [saving, setSaving] = useState(false)

    // Ambil daftar kelas sekali di awal untuk mengisi dropdown pada form tambah/ubah siswa
    useEffect(() => {
        axiosClient
            .get('/admin/kelas', { params: { page: 1, limit: 100 } })
            .then((res) => setKelasOptions(res.data.data?.items || []))
            .catch(() => setKelasOptions([]))
    }, [])

    const loadSiswa = useCallback(async (targetPage) => {
        setLoading(true)
        try {
            const res = await axiosClient.get('/admin/siswa', { params: { page: targetPage, limit: PAGE_SIZE } })
            const payload = res.data.data
            setSiswaList(payload?.items || [])
            setPagination(payload?.pagination || DEFAULT_PAGINATION)
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal memuat data siswa' })
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        loadSiswa(page)
    }, [page, loadSiswa])

    const openCreateModal = () => {
        setEditingId(null)
        setForm(EMPTY_FORM)
        setFormError('')
        setModalOpen(true)
    }

    const openEditModal = (siswa) => {
        setEditingId(siswa.id)
        setForm({
            username: siswa.username,
            password: '',
            nisn: siswa.nisn,
            nama: siswa.nama,
            class_id: siswa.class_id,
            alamat: siswa.alamat || '',
            no_telp: siswa.no_telp || '',
        })
        setFormError('')
        setModalOpen(true)
    }

    const closeModal = () => {
        if (saving) return
        setModalOpen(false)
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setFormError('')
        setSaving(true)

        try {
            if (editingId) {
                // Update TIDAK menyertakan username/password — kredensial login tidak diubah lewat form ini
                await axiosClient.put(`/admin/siswa/${editingId}`, {
                    nisn: form.nisn,
                    nama: form.nama,
                    class_id: Number(form.class_id),
                    alamat: form.alamat,
                    no_telp: form.no_telp,
                })
                Toast.fire({ icon: 'success', title: 'Data siswa berhasil diperbarui' })
            } else {
                await axiosClient.post('/admin/siswa', {
                    username: form.username,
                    password: form.password,
                    nisn: form.nisn,
                    nama: form.nama,
                    class_id: Number(form.class_id),
                    alamat: form.alamat,
                    no_telp: form.no_telp,
                })
                Toast.fire({ icon: 'success', title: 'Siswa dan akun login berhasil dibuat' })
            }
            setModalOpen(false)
            loadSiswa(page)
        } catch (err) {
            setFormError(err.response?.data?.message || 'Gagal menyimpan data siswa')
        } finally {
            setSaving(false)
        }
    }

    const handleDelete = async (siswa) => {
        const result = await Swal.fire({
            icon: 'warning',
            title: `Hapus siswa "${siswa.nama}"?`,
            text: 'Akun login siswa ini akan ikut terhapus. Tindakan ini tidak bisa dibatalkan.',
            showCancelButton: true,
            confirmButtonText: 'Ya, hapus',
            cancelButtonText: 'Batal',
            confirmButtonColor: '#ef4444',
            cancelButtonColor: '#9aa3b8',
            reverseButtons: true,
        })

        if (!result.isConfirmed) return

        try {
            await axiosClient.delete(`/admin/siswa/${siswa.id}`)
            Toast.fire({ icon: 'success', title: 'Siswa dan akun login berhasil dihapus' })

            if (siswaList.length === 1 && page > 1) {
                setPage(page - 1)
            } else {
                loadSiswa(page)
            }
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal menghapus siswa' })
        }
    }

    return (
        <div className="app-shell">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

            <div className="main-shell">
                <Topbar title="Data Siswa" onMenuClick={() => setSidebarOpen((prev) => !prev)} />

                <main className="main-content">
                    <div className="breadcrumb">
                        SchoolApp · Admin · <span>Data Siswa</span>
                    </div>
                    <div className="page-header">
                        <h1>Data Siswa</h1>
                        <p>Kelola data siswa. Menambah siswa baru otomatis membuat akun login untuknya.</p>
                    </div>

                    <div className="card">
                        <div className="toolbar">
                            <span style={{ fontSize: 13, color: 'var(--text2)' }}>
                                {pagination.total} siswa terdaftar
                            </span>
                            <button type="button" className="btn btn-primary" onClick={openCreateModal}>
                                + Tambah Siswa
                            </button>
                        </div>

                        <div className="table-wrap">
                            {loading ? (
                                <div className="empty-state">Memuat data...</div>
                            ) : siswaList.length === 0 ? (
                                <div className="empty-state">Belum ada data siswa.</div>
                            ) : (
                                <table className="data-table">
                                    <thead>
                                        <tr>
                                            <th>Nama</th>
                                            <th>NISN</th>
                                            <th>Username</th>
                                            <th>Kelas</th>
                                            <th>No. Telp</th>
                                            <th style={{ textAlign: 'right' }}>Aksi</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {siswaList.map((s) => (
                                            <tr key={s.id}>
                                                <td style={{ fontWeight: 600 }}>{s.nama}</td>
                                                <td>{s.nisn}</td>
                                                <td>{s.username}</td>
                                                <td>
                                                    <span className="chip">{s.nama_kelas}</span>
                                                </td>
                                                <td>{s.no_telp || '-'}</td>
                                                <td>
                                                    <div className="table-actions" style={{ justifyContent: 'flex-end' }}>
                                                        <button type="button" className="btn btn-icon-ghost" onClick={() => openEditModal(s)}>
                                                            Ubah
                                                        </button>
                                                        <button type="button" className="btn btn-danger-ghost" onClick={() => handleDelete(s)}>
                                                            Hapus
                                                        </button>
                                                    </div>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </div>

                        {!loading && pagination.total > 0 && (
                            <div className="pagination-bar">
                                <span className="pagination-info">
                                    Menampilkan halaman {pagination.page} dari {pagination.total_pages} ({pagination.total} siswa total)
                                </span>
                                <div className="pagination-controls">
                                    <button
                                        type="button"
                                        className="btn btn-outline"
                                        disabled={pagination.page <= 1}
                                        onClick={() => setPage((p) => Math.max(1, p - 1))}
                                    >
                                        ← Sebelumnya
                                    </button>
                                    <button
                                        type="button"
                                        className="btn btn-outline"
                                        disabled={pagination.page >= pagination.total_pages}
                                        onClick={() => setPage((p) => Math.min(pagination.total_pages, p + 1))}
                                    >
                                        Selanjutnya →
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>
                </main>
            </div>

            {modalOpen && (
                <div className="modal-overlay" onClick={closeModal}>
                    <div className="modal-box" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h3>{editingId ? 'Ubah Data Siswa' : 'Tambah Siswa Baru'}</h3>
                            <button type="button" className="modal-close" onClick={closeModal}>
                                ✕
                            </button>
                        </div>

                        <form onSubmit={handleSubmit}>
                            <div className="modal-body">
                                {formError && <div className="form-error">{formError}</div>}

                                {!editingId && (
                                    <>
                                        <div className="form-field">
                                            <label htmlFor="username">Username (untuk login siswa)</label>
                                            <input
                                                id="username"
                                                type="text"
                                                placeholder="Contoh: siswa3"
                                                value={form.username}
                                                onChange={(e) => setForm({ ...form, username: e.target.value })}
                                                required
                                            />
                                        </div>

                                        <div className="form-field">
                                            <label htmlFor="password">Password Awal</label>
                                            <input
                                                id="password"
                                                type="text"
                                                placeholder="Minimal 6 karakter"
                                                value={form.password}
                                                onChange={(e) => setForm({ ...form, password: e.target.value })}
                                                required
                                            />
                                        </div>
                                    </>
                                )}

                                {editingId && (
                                    <div className="form-field">
                                        <label>Username</label>
                                        <input type="text" value={form.username} disabled style={{ opacity: 0.6 }} />
                                    </div>
                                )}

                                <div className="form-field">
                                    <label htmlFor="nisn">NISN</label>
                                    <input
                                        id="nisn"
                                        type="text"
                                        placeholder="Contoh: 0051234567"
                                        value={form.nisn}
                                        onChange={(e) => setForm({ ...form, nisn: e.target.value })}
                                        required
                                    />
                                </div>

                                <div className="form-field">
                                    <label htmlFor="nama">Nama Lengkap</label>
                                    <input
                                        id="nama"
                                        type="text"
                                        placeholder="Contoh: Budi Santoso"
                                        value={form.nama}
                                        onChange={(e) => setForm({ ...form, nama: e.target.value })}
                                        required
                                    />
                                </div>

                                <div className="form-field">
                                    <label htmlFor="class_id">Kelas</label>
                                    <select
                                        id="class_id"
                                        value={form.class_id}
                                        onChange={(e) => setForm({ ...form, class_id: e.target.value })}
                                        required
                                    >
                                        <option value="" disabled>
                                            Pilih kelas
                                        </option>
                                        {kelasOptions.map((k) => (
                                            <option key={k.id} value={k.id}>
                                                {k.nama_kelas} (Tingkat {k.tingkat})
                                            </option>
                                        ))}
                                    </select>
                                </div>

                                <div className="form-field">
                                    <label htmlFor="alamat">Alamat</label>
                                    <input
                                        id="alamat"
                                        type="text"
                                        placeholder="Opsional"
                                        value={form.alamat}
                                        onChange={(e) => setForm({ ...form, alamat: e.target.value })}
                                    />
                                </div>

                                <div className="form-field">
                                    <label htmlFor="no_telp">No. Telepon</label>
                                    <input
                                        id="no_telp"
                                        type="text"
                                        placeholder="Opsional"
                                        value={form.no_telp}
                                        onChange={(e) => setForm({ ...form, no_telp: e.target.value })}
                                    />
                                </div>
                            </div>

                            <div className="modal-footer">
                                <button type="button" className="btn btn-outline" onClick={closeModal} disabled={saving}>
                                    Batal
                                </button>
                                <button type="submit" className="btn btn-primary" disabled={saving}>
                                    {saving ? 'Menyimpan...' : 'Simpan'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    )
}
