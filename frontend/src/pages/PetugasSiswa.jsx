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

// Halaman ini SENGAJA read-only (mengikuti pola yang sama seperti Data Kelas untuk petugas):
// tidak ada tombol tambah/ubah/hapus sama sekali, dan backend-nya cuma mendaftarkan route GET
// untuk role petugas (lihat router.go) — jadi bukan cuma disembunyikan di UI.
export default function PetugasSiswa() {
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [siswaList, setSiswaList] = useState([])
    const [loading, setLoading] = useState(true)
    const [search, setSearch] = useState('')

    const [page, setPage] = useState(1)
    const [pagination, setPagination] = useState(DEFAULT_PAGINATION)

    const loadSiswa = useCallback(async (targetPage) => {
        setLoading(true)
        try {
            const res = await axiosClient.get('/petugas/siswa', { params: { page: targetPage, limit: PAGE_SIZE } })
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

    const filteredList = siswaList.filter((s) => {
        const q = search.trim().toLowerCase()
        if (!q) return true
        return (
            s.nama.toLowerCase().includes(q) ||
            s.nisn.toLowerCase().includes(q) ||
            s.nama_kelas.toLowerCase().includes(q)
        )
    })

    return (
        <div className="app-shell">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

            <div className="main-shell">
                <Topbar title="Data Siswa" onMenuClick={() => setSidebarOpen((prev) => !prev)} />

                <main className="main-content">
                    <div className="breadcrumb">
                        SchoolApp · Petugas · <span>Data Siswa</span>
                    </div>
                    <div className="page-header">
                        <h1>Data Siswa</h1>
                        <p>Telusuri data siswa untuk keperluan pencarian tagihan SPP.</p>
                    </div>

                    <div className="card">
                        <div className="toolbar">
                            <input
                                type="text"
                                className="toolbar-search"
                                placeholder="Cari di halaman ini (nama / NISN / kelas)..."
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                            />
                        </div>

                        <div className="table-wrap">
                            {loading ? (
                                <div className="empty-state">Memuat data...</div>
                            ) : filteredList.length === 0 ? (
                                <div className="empty-state">
                                    {search ? 'Tidak ada siswa yang cocok dengan pencarian di halaman ini.' : 'Belum ada data siswa.'}
                                </div>
                            ) : (
                                <table className="data-table">
                                    <thead>
                                        <tr>
                                            <th>Nama</th>
                                            <th>NISN</th>
                                            <th>Kelas</th>
                                            <th>No. Telp</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {filteredList.map((s) => (
                                            <tr key={s.id}>
                                                <td style={{ fontWeight: 600 }}>{s.nama}</td>
                                                <td>{s.nisn}</td>
                                                <td>
                                                    <span className="chip">{s.nama_kelas}</span>
                                                </td>
                                                <td>{s.no_telp || '-'}</td>
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
        </div>
    )
}
