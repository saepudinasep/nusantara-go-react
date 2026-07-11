import { useCallback, useEffect, useState } from 'react'
import Swal from 'sweetalert2'
import axiosClient from '../api/axiosClient'
import Sidebar from '../components/Sidebar'
import Topbar from '../components/Topbar'

const EMPTY_FORM = { username: '', password: '', nama: '', posisi: '' }
const PAGE_SIZE = 10
const DEFAULT_PAGINATION = { page: 1, limit: PAGE_SIZE, total: 0, total_pages: 1 }

const Toast = Swal.mixin({
    toast: true,
    position: 'top-end',
    showConfirmButton: false,
    timer: 2500,
    timerProgressBar: true,
})

export default function PetugasManagement() {
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [petugasList, setPetugasList] = useState([])
    const [loading, setLoading] = useState(true)

    const [page, setPage] = useState(1)
    const [pagination, setPagination] = useState(DEFAULT_PAGINATION)

    const [modalOpen, setModalOpen] = useState(false)
    const [editingId, setEditingId] = useState(null)
    const [form, setForm] = useState(EMPTY_FORM)
    const [formError, setFormError] = useState('')
    const [saving, setSaving] = useState(false)

    const loadPetugas = useCallback(async (targetPage) => {
        setLoading(true)
        try {
            const res = await axiosClient.get('/admin/petugas', { params: { page: targetPage, limit: PAGE_SIZE } })
            const payload = res.data.data
            setPetugasList(payload?.items || [])
            setPagination(payload?.pagination || DEFAULT_PAGINATION)
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal memuat data petugas' })
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        loadPetugas(page)
    }, [page, loadPetugas])

    const openCreateModal = () => {
        setEditingId(null)
        setForm(EMPTY_FORM)
        setFormError('')
        setModalOpen(true)
    }

    const openEditModal = (petugas) => {
        setEditingId(petugas.id)
        setForm({ username: petugas.username, password: '', nama: petugas.nama, posisi: petugas.posisi })
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
                await axiosClient.put(`/admin/petugas/${editingId}`, { nama: form.nama, posisi: form.posisi })
                Toast.fire({ icon: 'success', title: 'Data petugas berhasil diperbarui' })
            } else {
                await axiosClient.post('/admin/petugas', {
                    username: form.username,
                    password: form.password,
                    nama: form.nama,
                    posisi: form.posisi,
                })
                Toast.fire({ icon: 'success', title: 'Petugas dan akun login berhasil dibuat' })
            }
            setModalOpen(false)
            loadPetugas(page)
        } catch (err) {
            setFormError(err.response?.data?.message || 'Gagal menyimpan data petugas')
        } finally {
            setSaving(false)
        }
    }

    const handleDelete = async (petugas) => {
        const result = await Swal.fire({
            icon: 'warning',
            title: `Hapus petugas "${petugas.nama}"?`,
            text: 'Akun login petugas ini akan ikut terhapus. Tindakan ini tidak bisa dibatalkan.',
            showCancelButton: true,
            confirmButtonText: 'Ya, hapus',
            cancelButtonText: 'Batal',
            confirmButtonColor: '#ef4444',
            cancelButtonColor: '#9aa3b8',
            reverseButtons: true,
        })

        if (!result.isConfirmed) return

        try {
            await axiosClient.delete(`/admin/petugas/${petugas.id}`)
            Toast.fire({ icon: 'success', title: 'Petugas dan akun login berhasil dihapus' })

            if (petugasList.length === 1 && page > 1) {
                setPage(page - 1)
            } else {
                loadPetugas(page)
            }
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal menghapus petugas' })
        }
    }

    return (
        <div className="app-shell">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

            <div className="main-shell">
                <Topbar title="Data Petugas" onMenuClick={() => setSidebarOpen((prev) => !prev)} />

                <main className="main-content">
                    <div className="breadcrumb">
                        SchoolApp · Admin · <span>Data Petugas</span>
                    </div>
                    <div className="page-header">
                        <h1>Data Petugas</h1>
                        <p>Kelola data petugas SPP. Menambah petugas baru otomatis membuat akun login untuknya.</p>
                    </div>

                    <div className="card">
                        <div className="toolbar">
                            <span style={{ fontSize: 13, color: 'var(--text2)' }}>
                                {pagination.total} petugas terdaftar
                            </span>
                            <button type="button" className="btn btn-primary" onClick={openCreateModal}>
                                + Tambah Petugas
                            </button>
                        </div>

                        <div className="table-wrap">
                            {loading ? (
                                <div className="empty-state">Memuat data...</div>
                            ) : petugasList.length === 0 ? (
                                <div className="empty-state">Belum ada data petugas.</div>
                            ) : (
                                <table className="data-table">
                                    <thead>
                                        <tr>
                                            <th>Nama</th>
                                            <th>Username</th>
                                            <th>Posisi / Jabatan</th>
                                            <th style={{ textAlign: 'right' }}>Aksi</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {petugasList.map((p) => (
                                            <tr key={p.id}>
                                                <td style={{ fontWeight: 600 }}>{p.nama}</td>
                                                <td>{p.username}</td>
                                                <td>
                                                    <span className="chip">{p.posisi}</span>
                                                </td>
                                                <td>
                                                    <div className="table-actions" style={{ justifyContent: 'flex-end' }}>
                                                        <button type="button" className="btn btn-icon-ghost" onClick={() => openEditModal(p)}>
                                                            Ubah
                                                        </button>
                                                        <button type="button" className="btn btn-danger-ghost" onClick={() => handleDelete(p)}>
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
                                    Menampilkan halaman {pagination.page} dari {pagination.total_pages} ({pagination.total} petugas total)
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
                            <h3>{editingId ? 'Ubah Data Petugas' : 'Tambah Petugas Baru'}</h3>
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
                                            <label htmlFor="username">Username (untuk login petugas)</label>
                                            <input
                                                id="username"
                                                type="text"
                                                placeholder="Contoh: petugas3"
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
                                    <label htmlFor="nama">Nama Lengkap</label>
                                    <input
                                        id="nama"
                                        type="text"
                                        placeholder="Contoh: Siti Rahma"
                                        value={form.nama}
                                        onChange={(e) => setForm({ ...form, nama: e.target.value })}
                                        required
                                    />
                                </div>

                                <div className="form-field">
                                    <label htmlFor="posisi">Posisi / Jabatan</label>
                                    <input
                                        id="posisi"
                                        type="text"
                                        placeholder="Contoh: Kasir / Tata Usaha"
                                        value={form.posisi}
                                        onChange={(e) => setForm({ ...form, posisi: e.target.value })}
                                        required
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
