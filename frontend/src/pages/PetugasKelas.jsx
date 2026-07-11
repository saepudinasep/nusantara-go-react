import { useCallback, useEffect, useState } from 'react'
import Swal from 'sweetalert2'
import axiosClient from '../api/axiosClient'
import Sidebar from '../components/Sidebar'
import Topbar from '../components/Topbar'

const PAGE_SIZE = 10
const DEFAULT_PAGINATION = { page: 1, limit: PAGE_SIZE, total: 0, total_pages: 1 }

const Toast = Swal.mixin({
    toast: true,
    position: 'top-end',
    showConfirmButton: false,
    timer: 2500,
    timerProgressBar: true,
})

// Halaman ini SENGAJA read-only (sesuai checklist: petugas hanya boleh melihat daftar kelas
// untuk keperluan mencari data tagihan siswa) — tidak ada tombol tambah/ubah/hapus sama sekali,
// dan backend-nya pun cuma mendaftarkan route GET untuk role petugas (lihat router.go).
export default function PetugasKelas() {
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [kelasList, setKelasList] = useState([])
    const [loading, setLoading] = useState(true)
    const [search, setSearch] = useState('')

    const [page, setPage] = useState(1)
    const [pagination, setPagination] = useState(DEFAULT_PAGINATION)

    const loadKelas = useCallback(async (targetPage) => {
        setLoading(true)
        try {
            const res = await axiosClient.get('/petugas/kelas', { params: { page: targetPage, limit: PAGE_SIZE } })
            const payload = res.data.data
            setKelasList(payload?.items || [])
            setPagination(payload?.pagination || DEFAULT_PAGINATION)
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal memuat data kelas' })
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        loadKelas(page)
    }, [page, loadKelas])

    const filteredList = kelasList.filter((k) => {
        const q = search.trim().toLowerCase()
        if (!q) return true
        return k.nama_kelas.toLowerCase().includes(q) || String(k.tingkat).includes(q)
    })

    return (
        <div className="app-shell">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

            <div className="main-shell">
                <Topbar title="Data Kelas" onMenuClick={() => setSidebarOpen((prev) => !prev)} />

                <main className="main-content">
                    <div className="breadcrumb">
                        SchoolApp · Petugas · <span>Data Kelas</span>
                    </div>
                    <div className="page-header">
                        <h1>Data Kelas</h1>
                        <p>Lihat daftar kelas untuk membantu mencari data tagihan SPP siswa.</p>
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
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {filteredList.map((k) => (
                                            <tr key={k.id}>
                                                <td style={{ fontWeight: 600 }}>{k.nama_kelas}</td>
                                                <td>
                                                    <span className="chip">{k.tingkat}</span>
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
        </div>
    )
}
