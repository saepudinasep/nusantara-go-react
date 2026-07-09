import { useCallback, useEffect, useMemo, useState } from 'react'
import Swal from 'sweetalert2'
import axiosClient from '../api/axiosClient'
import Sidebar from '../components/Sidebar'
import Topbar from '../components/Topbar'

const TINGKAT_OPTIONS = [10, 11, 12]
const EMPTY_FORM = { nama_kelas: '', tingkat: 10 }
const PAGE_SIZE = 10

// Toast kecil di pojok kanan atas untuk notifikasi sukses/gagal (menggantikan toast custom sebelumnya)
const Toast = Swal.mixin({
    toast: true,
    position: 'top-end',
    showConfirmButton: false,
    timer: 2500,
    timerProgressBar: true,
})

export default function KelasManagement() {
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [kelasList, setKelasList] = useState([])
    const [loading, setLoading] = useState(true)
    const [search, setSearch] = useState('')

    const [page, setPage] = useState(1)
    const [pagination, setPagination] = useState({ page: 1, limit: PAGE_SIZE, total: 0, total_pages: 1 })
    console.log('pagination', pagination)

    const [modalOpen, setModalOpen] = useState(false)
    const [editingId, setEditingId] = useState(null)
    const [form, setForm] = useState(EMPTY_FORM)
    const [formError, setFormError] = useState('')
    const [saving, setSaving] = useState(false)

    const loadKelas = useCallback(async (targetPage) => {
        setLoading(true)
        try {
            const res = await axiosClient.get('/admin/kelas', { params: { page: targetPage, limit: PAGE_SIZE } })
            const { items, pagination: meta } = res.data.data
            setKelasList(items || [])
            setPagination(meta)
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal memuat data kelas' })
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        loadKelas(page)
    }, [page, loadKelas])

    const filteredList = useMemo(() => {
        const q = search.trim().toLowerCase()
        if (!q) return kelasList
        return kelasList.filter(
            (k) => k.nama_kelas.toLowerCase().includes(q) || String(k.tingkat).includes(q)
        )
    }, [kelasList, search])

    const openCreateModal = () => {
        setEditingId(null)
        setForm(EMPTY_FORM)
        setFormError('')
        setModalOpen(true)
    }

    const openEditModal = (kelas) => {
        setEditingId(kelas.id)
        setForm({ nama_kelas: kelas.nama_kelas, tingkat: kelas.tingkat })
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
            const payload = { nama_kelas: form.nama_kelas, tingkat: Number(form.tingkat) }
            if (editingId) {
                await axiosClient.put(`/admin/kelas/${editingId}`, payload)
                Toast.fire({ icon: 'success', title: 'Kelas berhasil diperbarui' })
            } else {
                await axiosClient.post('/admin/kelas', payload)
                Toast.fire({ icon: 'success', title: 'Kelas berhasil ditambahkan' })
            }
            setModalOpen(false)
            loadKelas(page)
        } catch (err) {
            setFormError(err.response?.data?.message || 'Gagal menyimpan data kelas')
        } finally {
            setSaving(false)
        }
    }

    const handleDelete = async (kelas) => {
        const result = await Swal.fire({
            icon: 'warning',
            title: `Hapus kelas "${kelas.nama_kelas}"?`,
            text: 'Tindakan ini tidak bisa dibatalkan.',
            showCancelButton: true,
            confirmButtonText: 'Ya, hapus',
            cancelButtonText: 'Batal',
            confirmButtonColor: '#ef4444',
            cancelButtonColor: '#9aa3b8',
            reverseButtons: true,
        })

        if (!result.isConfirmed) return

        try {
            await axiosClient.delete(`/admin/kelas/${kelas.id}`)
            Toast.fire({ icon: 'success', title: 'Kelas berhasil dihapus' })

            // Kalau item terakhir di halaman ini dihapus dan bukan halaman pertama, mundur satu halaman
            if (kelasList.length === 1 && page > 1) {
                setPage(page - 1)
            } else {
                loadKelas(page)
            }
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal menghapus kelas' })
        }
    }

    return (
        <div className="app-shell">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

            <div className="main-shell">
                <Topbar title="Data Kelas" onMenuClick={() => setSidebarOpen((prev) => !prev)} />

                <main className="main-content">
                    <div className="breadcrumb">
                        SchoolApp · Admin · <span>Data Kelas</span>
                    </div>
                    <div className="page-header">
                        <h1>Data Kelas</h1>
                        <p>Kelola daftar kelas: tambah, ubah, atau hapus data kelas.</p>
                    </div>

                    <div className="card">
                        <div className="toolbar">
                            <input
                                type="text"
                                className="toolbar-search"
                                placeholder="Cari di halaman ini (nama kelas / tingkat)..."
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                            />
                            <button type="button" className="btn btn-primary" onClick={openCreateModal}>
                                + Tambah Kelas
                            </button>
                        </div>

                        <div className="table-wrap">
                            {loading ? (
                                <div className="empty-state">Memuat data...</div>
                            ) : filteredList.length === 0 ? (
                                <div className="empty-state">
                                    {search ? 'Tidak ada kelas yang cocok dengan pencarian di halaman ini.' : 'Belum ada data kelas.'}
                                </div>
                            ) : (
                                <table className="data-table">
                                    <thead>
                                        <tr>
                                            <th>Nama Kelas</th>
                                            <th>Tingkat</th>
                                            <th style={{ textAlign: 'right' }}>Aksi</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {filteredList.map((k) => (
                                            <tr key={k.id}>
                                                <td style={{ fontWeight: 600 }}>{k.nama_kelas}</td>
                                                <td>
                                                    <span className="chip">{k.tingkat}</span>
                                                </td>
                                                <td>
                                                    <div className="table-actions" style={{ justifyContent: 'flex-end' }}>
                                                        <button type="button" className="btn btn-icon-ghost" onClick={() => openEditModal(k)}>
                                                            Ubah
                                                        </button>
                                                        <button
                                                            type="button"
                                                            className="btn btn-danger-ghost"
                                                            onClick={() => handleDelete(k)}
                                                        >
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
                                    Menampilkan halaman {pagination.page} dari {pagination.total_pages} ({pagination.total} kelas total)
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
                            <h3>{editingId ? 'Ubah Kelas' : 'Tambah Kelas'}</h3>
                            <button type="button" className="modal-close" onClick={closeModal}>
                                ✕
                            </button>
                        </div>

                        <form onSubmit={handleSubmit}>
                            <div className="modal-body">
                                {formError && <div className="form-error">{formError}</div>}

                                <div className="form-field">
                                    <label htmlFor="nama_kelas">Nama Kelas</label>
                                    <input
                                        id="nama_kelas"
                                        type="text"
                                        placeholder="Contoh: XA, XI RPL"
                                        value={form.nama_kelas}
                                        onChange={(e) => setForm({ ...form, nama_kelas: e.target.value })}
                                        required
                                    />
                                </div>

                                <div className="form-field">
                                    <label htmlFor="tingkat">Tingkat</label>
                                    <select
                                        id="tingkat"
                                        value={form.tingkat}
                                        onChange={(e) => setForm({ ...form, tingkat: e.target.value })}
                                    >
                                        {TINGKAT_OPTIONS.map((t) => (
                                            <option key={t} value={t}>
                                                {t}
                                            </option>
                                        ))}
                                    </select>
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
