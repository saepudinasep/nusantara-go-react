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
//
// Kolom "Status SPP Bulan Ini" + filter "hanya yang nunggak" dipakai supaya petugas bisa langsung
// tahu SIAPA SAJA siswa yang belum bayar (bukan cuma angka total seperti di dashboard).
export default function PetugasSiswa() {
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [siswaList, setSiswaList] = useState([])
    const [loading, setLoading] = useState(true)
    const [search, setSearch] = useState('')
    const [onlyUnpaid, setOnlyUnpaid] = useState(false)

    const [page, setPage] = useState(1)
    const [pagination, setPagination] = useState(DEFAULT_PAGINATION)

    const loadSiswa = useCallback(async (targetPage, unpaidOnly) => {
        setLoading(true)
        try {
            const res = await axiosClient.get('/petugas/siswa', {
                params: { page: targetPage, limit: PAGE_SIZE, with_status: true, only_unpaid: unpaidOnly },
            })
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
        loadSiswa(page, onlyUnpaid)
    }, [page, onlyUnpaid, loadSiswa])

    const handleToggleUnpaid = (checked) => {
        setOnlyUnpaid(checked)
        setPage(1)
    }

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
                        <p>Telusuri data siswa beserta status pembayaran SPP bulan ini untuk keperluan penagihan.</p>
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
                            <label style={{ display: 'flex', alignItems: 'center', gap: 8, fontSize: 13, color: 'var(--text2)' }}>
                                <input
                                    type="checkbox"
                                    checked={onlyUnpaid}
                                    onChange={(e) => handleToggleUnpaid(e.target.checked)}
                                />
                                Tampilkan hanya yang nunggak
                            </label>
                        </div>

                        <div className="table-wrap">
                            {loading ? (
                                <div className="empty-state">Memuat data...</div>
                            ) : filteredList.length === 0 ? (
                                <div className="empty-state">
                                    {onlyUnpaid
                                        ? 'Tidak ada siswa yang menunggak. 🎉'
                                        : search
                                            ? 'Tidak ada siswa yang cocok dengan pencarian di halaman ini.'
                                            : 'Belum ada data siswa.'}
                                </div>
                            ) : (
                                <table className="data-table">
                                    <thead>
                                        <tr>
                                            <th>Nama</th>
                                            <th>NISN</th>
                                            <th>Kelas</th>
                                            <th>No. Telp</th>
                                            <th>Status SPP Bulan Ini</th>
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
                                                <td>
                                                    <span className={`status-badge ${s.status_spp_bulan_ini === 'Lunas' ? 'lunas' : 'nunggak'}`}>
                                                        {s.status_spp_bulan_ini}
                                                    </span>
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
                                    Menampilkan halaman {pagination.page} dari {pagination.total_pages} ({pagination.total} siswa
                                    {onlyUnpaid ? ' menunggak' : ' total'})
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
