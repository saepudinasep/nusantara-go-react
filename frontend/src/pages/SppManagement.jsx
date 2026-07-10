import { useCallback, useEffect, useState } from 'react'
import Swal from 'sweetalert2'
import axiosClient from '../api/axiosClient'
import Sidebar from '../components/Sidebar'
import Topbar from '../components/Topbar'

const EMPTY_FORM = { tahun_ajaran: '', nominal: '' }
const PAGE_SIZE = 10
const DEFAULT_PAGINATION = { page: 1, limit: PAGE_SIZE, total: 0, total_pages: 1 }

const Toast = Swal.mixin({
    toast: true,
    position: 'top-end',
    showConfirmButton: false,
    timer: 2500,
    timerProgressBar: true,
})

function formatRupiah(value) {
    const num = Number(value)
    if (Number.isNaN(num)) return value
    return 'Rp' + num.toLocaleString('id-ID')
}

export default function SppManagement() {
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [sppList, setSppList] = useState([])
    const [loading, setLoading] = useState(true)

    const [page, setPage] = useState(1)
    const [pagination, setPagination] = useState(DEFAULT_PAGINATION)

    const [modalOpen, setModalOpen] = useState(false)
    const [editingId, setEditingId] = useState(null)
    const [form, setForm] = useState(EMPTY_FORM)
    const [formError, setFormError] = useState('')
    const [saving, setSaving] = useState(false)

    const loadSpp = useCallback(async (targetPage) => {
        setLoading(true)
        try {
            const res = await axiosClient.get('/admin/spp', { params: { page: targetPage, limit: PAGE_SIZE } })
            const payload = res.data.data

            if (Array.isArray(payload)) {
                console.warn(
                    'Backend mengembalikan array SPP tanpa info pagination. ' +
                    'Kemungkinan server backend belum di-restart setelah update — jalankan ulang "go run cmd/api/main.go".'
                )
                setSppList(payload)
                setPagination({ ...DEFAULT_PAGINATION, total: payload.length })
                return
            }

            setSppList(payload?.items || [])
            setPagination(payload?.pagination || DEFAULT_PAGINATION)
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal memuat data SPP' })
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        loadSpp(page)
    }, [page, loadSpp])

    const openCreateModal = () => {
        setEditingId(null)
        setForm(EMPTY_FORM)
        setFormError('')
        setModalOpen(true)
    }

    const openEditModal = (spp) => {
        setEditingId(spp.id)
        setForm({ tahun_ajaran: spp.tahun_ajaran, nominal: spp.nominal })
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
            const payload = { tahun_ajaran: form.tahun_ajaran, nominal: Number(form.nominal) }
            if (editingId) {
                await axiosClient.put(`/admin/spp/${editingId}`, payload)
                Toast.fire({ icon: 'success', title: 'Data SPP berhasil diperbarui' })
            } else {
                await axiosClient.post('/admin/spp', payload)
                Toast.fire({ icon: 'success', title: 'Data SPP berhasil ditambahkan' })
            }
            setModalOpen(false)
            loadSpp(page)
        } catch (err) {
            // Pesan di sini sudah pasti aman ditampilkan — backend menerjemahkan semua error SQL
            // (duplikat, FK constraint, lock timeout) jadi pesan domain sebelum sampai ke response.
            setFormError(err.response?.data?.message || 'Gagal menyimpan data SPP')
        } finally {
            setSaving(false)
        }
    }

    const handleDelete = async (spp) => {
        const result = await Swal.fire({
            icon: 'warning',
            title: `Hapus data SPP "${spp.tahun_ajaran}"?`,
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
            await axiosClient.delete(`/admin/spp/${spp.id}`)
            Toast.fire({ icon: 'success', title: 'Data SPP berhasil dihapus' })

            if (sppList.length === 1 && page > 1) {
                setPage(page - 1)
            } else {
                loadSpp(page)
            }
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal menghapus data SPP' })
        }
    }

    return (
        <div className="app-shell">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

            <div className="main-shell">
                <Topbar title="Data SPP" onMenuClick={() => setSidebarOpen((prev) => !prev)} />

                <main className="main-content">
                    <div className="breadcrumb">
                        SchoolApp · Admin · <span>Data SPP</span>
                    </div>
                    <div className="page-header">
                        <h1>Data SPP</h1>
                        <p>Kelola master nominal SPP per tahun ajaran.</p>
                    </div>

                    <div className="card">
                        <div className="toolbar">
                            <span style={{ fontSize: 13, color: 'var(--text2)' }}>
                                {pagination.total} data SPP terdaftar
                            </span>
                            <button type="button" className="btn btn-primary" onClick={openCreateModal}>
                                + Tambah SPP
                            </button>
                        </div>

                        <div className="table-wrap">
                            {loading ? (
                                <div className="empty-state">Memuat data...</div>
                            ) : sppList.length === 0 ? (
                                <div className="empty-state">Belum ada data SPP.</div>
                            ) : (
                                <table className="data-table">
                                    <thead>
                                        <tr>
                                            <th>Tahun Ajaran</th>
                                            <th>Nominal</th>
                                            <th style={{ textAlign: 'right' }}>Aksi</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {sppList.map((s) => (
                                            <tr key={s.id}>
                                                <td style={{ fontWeight: 600 }}>{s.tahun_ajaran}</td>
                                                <td>{formatRupiah(s.nominal)}</td>
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
                                    Menampilkan halaman {pagination.page} dari {pagination.total_pages} ({pagination.total} data total)
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
                            <h3>{editingId ? 'Ubah Data SPP' : 'Tambah Data SPP'}</h3>
                            <button type="button" className="modal-close" onClick={closeModal}>
                                ✕
                            </button>
                        </div>

                        <form onSubmit={handleSubmit}>
                            <div className="modal-body">
                                {formError && <div className="form-error">{formError}</div>}

                                <div className="form-field">
                                    <label htmlFor="tahun_ajaran">Tahun Ajaran</label>
                                    <input
                                        id="tahun_ajaran"
                                        type="text"
                                        placeholder="Contoh: 2025/2026"
                                        value={form.tahun_ajaran}
                                        onChange={(e) => setForm({ ...form, tahun_ajaran: e.target.value })}
                                        required
                                    />
                                </div>

                                <div className="form-field">
                                    <label htmlFor="nominal">Nominal (Rp)</label>
                                    <input
                                        id="nominal"
                                        type="number"
                                        min="1"
                                        // step="1000"
                                        placeholder="Contoh: 150000"
                                        value={form.nominal}
                                        onChange={(e) => setForm({ ...form, nominal: e.target.value })}
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
