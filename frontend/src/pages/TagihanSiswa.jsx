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

function formatRupiah(value) {
    const num = Number(value)
    if (Number.isNaN(num)) return value
    return 'Rp' + num.toLocaleString('id-ID')
}

export default function TagihanSiswa() {
    const [sidebarOpen, setSidebarOpen] = useState(false)

    const [loadingTagihan, setLoadingTagihan] = useState(true)
    const [tagihanList, setTagihanList] = useState([])

    const [loadingRiwayat, setLoadingRiwayat] = useState(true)
    const [riwayatList, setRiwayatList] = useState([])
    const [page, setPage] = useState(1)
    const [pagination, setPagination] = useState(DEFAULT_PAGINATION)

    useEffect(() => {
        setLoadingTagihan(true)
        axiosClient
            .get('/siswa/tagihan')
            .then((res) => setTagihanList(res.data.data?.tagihan || []))
            .catch((err) => Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal memuat tagihan' }))
            .finally(() => setLoadingTagihan(false))
    }, [])

    const loadRiwayat = useCallback(async (targetPage) => {
        setLoadingRiwayat(true)
        try {
            const res = await axiosClient.get('/siswa/riwayat', { params: { page: targetPage, limit: PAGE_SIZE } })
            const payload = res.data.data
            setRiwayatList(payload?.items || [])
            setPagination(payload?.pagination || DEFAULT_PAGINATION)
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal memuat riwayat pembayaran' })
        } finally {
            setLoadingRiwayat(false)
        }
    }, [])

    useEffect(() => {
        loadRiwayat(page)
    }, [page, loadRiwayat])

    return (
        <div className="app-shell">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

            <div className="main-shell">
                <Topbar title="Tagihan &amp; Riwayat" onMenuClick={() => setSidebarOpen((prev) => !prev)} />

                <main className="main-content">
                    <div className="breadcrumb">
                        SchoolApp · Siswa · <span>Tagihan &amp; Riwayat</span>
                    </div>
                    <div className="page-header">
                        <h1>Tagihan &amp; Riwayat Pembayaran</h1>
                        <p>Lihat status pembayaran SPP kamu per bulan, dan riwayat transaksi yang sudah dilakukan.</p>
                    </div>

                    {/* ---- Status Tagihan per Bulan ---- */}
                    {loadingTagihan ? (
                        <div className="card">
                            <div className="empty-state">Memuat data tagihan...</div>
                        </div>
                    ) : tagihanList.length === 0 ? (
                        <div className="card">
                            <div className="empty-state">Belum ada jenis SPP yang terdaftar.</div>
                        </div>
                    ) : (
                        tagihanList.map((spp) => (
                            <div className="card" key={spp.spp_id}>
                                <div className="card-header">
                                    <div>
                                        <div className="card-title">SPP Tahun Ajaran {spp.tahun_ajaran}</div>
                                        <div className="card-subtitle">Nominal per bulan: {formatRupiah(spp.nominal)}</div>
                                    </div>
                                </div>
                                <div className="bulan-grid">
                                    {spp.bulanan.map((b) => (
                                        <div key={b.bulan} className={`bulan-item ${b.status === 'Lunas' ? 'lunas' : 'belum'}`}>
                                            <div className="bulan-name">{b.bulan}</div>
                                            <div className="bulan-status">{b.status}</div>
                                            {b.status === 'Lunas' && <div className="bulan-tanggal">{b.tanggal_bayar}</div>}
                                        </div>
                                    ))}
                                </div>
                            </div>
                        ))
                    )}

                    {/* ---- Riwayat Pembayaran Lengkap ---- */}
                    <div className="card">
                        <div className="card-header">
                            <div>
                                <div className="card-title">Riwayat Pembayaran</div>
                                <div className="card-subtitle">Seluruh transaksi pembayaran SPP yang sudah kamu lakukan</div>
                            </div>
                        </div>

                        <div className="table-wrap">
                            {loadingRiwayat ? (
                                <div className="empty-state">Memuat riwayat...</div>
                            ) : riwayatList.length === 0 ? (
                                <div className="empty-state">Belum ada riwayat pembayaran.</div>
                            ) : (
                                <table className="data-table">
                                    <thead>
                                        <tr>
                                            <th>Tanggal</th>
                                            <th>Bulan</th>
                                            <th>Tahun Ajaran</th>
                                            <th>Jumlah</th>
                                            <th>Diproses oleh</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {riwayatList.map((r) => (
                                            <tr key={r.id}>
                                                <td>{r.tanggal_bayar}</td>
                                                <td>{r.bulan_dibayar}</td>
                                                <td>
                                                    <span className="chip">{r.tahun_ajaran}</span>
                                                </td>
                                                <td>{formatRupiah(r.jumlah_bayar)}</td>
                                                <td>{r.staff_nama}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </div>

                        {!loadingRiwayat && pagination.total > 0 && (
                            <div className="pagination-bar">
                                <span className="pagination-info">
                                    Menampilkan halaman {pagination.page} dari {pagination.total_pages} ({pagination.total} transaksi total)
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
