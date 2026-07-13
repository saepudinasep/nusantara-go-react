import { useState } from 'react'
import Swal from 'sweetalert2'
import axiosClient from '../api/axiosClient'
import { useAuth } from '../context/AuthContext'
import Sidebar from '../components/Sidebar'
import Topbar from '../components/Topbar'

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

function firstDayOfMonthISO() {
    const d = new Date()
    return new Date(d.getFullYear(), d.getMonth(), 1).toISOString().slice(0, 10)
}

function todayISO() {
    return new Date().toISOString().slice(0, 10)
}

export default function Laporan() {
    const { user } = useAuth()
    const isAdmin = user.role === 'admin'
    const basePath = isAdmin ? '/admin' : '/petugas'

    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [tanggalDari, setTanggalDari] = useState(firstDayOfMonthISO())
    const [tanggalSampai, setTanggalSampai] = useState(todayISO())

    const [loading, setLoading] = useState(false)
    const [hasGenerated, setHasGenerated] = useState(false)
    const [summary, setSummary] = useState(null)
    const [breakdown, setBreakdown] = useState([])
    const [transactions, setTransactions] = useState([])

    const handleGenerate = async (e) => {
        e.preventDefault()
        setLoading(true)
        try {
            const res = await axiosClient.get(`${basePath}/laporan`, {
                params: { tanggal_dari: tanggalDari, tanggal_sampai: tanggalSampai },
            })
            const data = res.data.data
            setSummary(data.summary)
            setBreakdown(data.breakdown || [])
            setTransactions(data.transactions || [])
            setHasGenerated(true)
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal menyusun laporan' })
        } finally {
            setLoading(false)
        }
    }

    const handlePrint = () => {
        window.print()
    }

    return (
        <div className="app-shell">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

            <div className="main-shell">
                <Topbar title={isAdmin ? 'Laporan Global' : 'Laporan Harian'} onMenuClick={() => setSidebarOpen((prev) => !prev)} />

                <main className="main-content">
                    <div className="breadcrumb no-print">
                        SchoolApp · {isAdmin ? 'Admin' : 'Petugas'} · <span>{isAdmin ? 'Laporan Global' : 'Laporan Harian'}</span>
                    </div>
                    <div className="page-header no-print">
                        <h1>{isAdmin ? 'Laporan Global' : 'Laporan Harian'}</h1>
                        <p>
                            {isAdmin
                                ? 'Rekap keuangan seluruh sekolah pada rentang tanggal tertentu, termasuk rincian per petugas.'
                                : 'Rekap setoran yang kamu proses sendiri pada rentang tanggal tertentu.'}
                        </p>
                    </div>

                    <div className="card no-print">
                        <form onSubmit={handleGenerate} className="filter-bar" style={{ paddingBottom: 18 }}>
                            <label style={{ fontSize: 13, color: 'var(--text2)' }}>Dari</label>
                            <input type="date" value={tanggalDari} onChange={(e) => setTanggalDari(e.target.value)} required />
                            <label style={{ fontSize: 13, color: 'var(--text2)' }}>Sampai</label>
                            <input type="date" value={tanggalSampai} onChange={(e) => setTanggalSampai(e.target.value)} required />
                            <button type="submit" className="btn btn-primary" disabled={loading}>
                                {loading ? 'Menyusun...' : 'Tampilkan Laporan'}
                            </button>
                            {hasGenerated && (
                                <button type="button" className="btn btn-outline" onClick={handlePrint}>
                                    🖨️ Cetak Laporan
                                </button>
                            )}
                        </form>
                    </div>

                    {hasGenerated && summary && (
                        <div id="laporan-cetak">
                            <div className="print-only-header">
                                <h2>{isAdmin ? 'Laporan Keuangan Global' : 'Laporan Harian Petugas'}</h2>
                                <p>
                                    Periode: {summary.tanggal_dari} s/d {summary.tanggal_sampai}
                                    {!isAdmin && <> — Petugas: {user.username}</>}
                                </p>
                            </div>

                            <div className="stats-grid" style={{ gridTemplateColumns: 'repeat(2, 1fr)' }}>
                                <div className="stat-card blue">
                                    <div className="stat-card-top">
                                        <div className="stat-card-label">Jumlah Transaksi</div>
                                    </div>
                                    <div>
                                        <div className="stat-card-val">{summary.jumlah_transaksi}</div>
                                        <div className="stat-card-sub">
                                            {summary.tanggal_dari} s/d {summary.tanggal_sampai}
                                        </div>
                                    </div>
                                </div>
                                <div className="stat-card green">
                                    <div className="stat-card-top">
                                        <div className="stat-card-label">Total Nominal</div>
                                    </div>
                                    <div>
                                        <div className="stat-card-val">{formatRupiah(summary.total_nominal)}</div>
                                        <div className="stat-card-sub">Seluruh transaksi pada periode ini</div>
                                    </div>
                                </div>
                            </div>

                            {isAdmin && (
                                <div className="card">
                                    <div className="card-header">
                                        <div>
                                            <div className="card-title">Rekap per Petugas</div>
                                            <div className="card-subtitle">Jumlah transaksi &amp; total setoran masing-masing petugas</div>
                                        </div>
                                    </div>
                                    <div className="table-wrap">
                                        {breakdown.length === 0 ? (
                                            <div className="empty-state">Belum ada data petugas.</div>
                                        ) : (
                                            <table className="data-table">
                                                <thead>
                                                    <tr>
                                                        <th>Petugas</th>
                                                        <th>Jumlah Transaksi</th>
                                                        <th>Total Setoran</th>
                                                    </tr>
                                                </thead>
                                                <tbody>
                                                    {breakdown.map((b) => (
                                                        <tr key={b.staff_id}>
                                                            <td style={{ fontWeight: 600 }}>{b.staff_nama}</td>
                                                            <td>{b.jumlah_transaksi}</td>
                                                            <td>{formatRupiah(b.total_nominal)}</td>
                                                        </tr>
                                                    ))}
                                                </tbody>
                                            </table>
                                        )}
                                    </div>
                                </div>
                            )}

                            <div className="card">
                                <div className="card-header">
                                    <div>
                                        <div className="card-title">Detail Transaksi</div>
                                        <div className="card-subtitle">{transactions.length} transaksi pada periode ini</div>
                                    </div>
                                </div>
                                <div className="table-wrap">
                                    {transactions.length === 0 ? (
                                        <div className="empty-state">Tidak ada transaksi pada periode ini.</div>
                                    ) : (
                                        <table className="data-table">
                                            <thead>
                                                <tr>
                                                    <th>Tanggal</th>
                                                    <th>Siswa</th>
                                                    <th>Kelas</th>
                                                    <th>Bulan</th>
                                                    <th>Jumlah</th>
                                                    {isAdmin && <th>Petugas</th>}
                                                </tr>
                                            </thead>
                                            <tbody>
                                                {transactions.map((t) => (
                                                    <tr key={t.id}>
                                                        <td>{t.tanggal_bayar}</td>
                                                        <td style={{ fontWeight: 600 }}>
                                                            {t.student_nama}
                                                            <div style={{ fontSize: 11, color: 'var(--gray3)', fontWeight: 400 }}>{t.nisn}</div>
                                                        </td>
                                                        <td>
                                                            <span className="chip">{t.nama_kelas}</span>
                                                        </td>
                                                        <td>{t.bulan_dibayar}</td>
                                                        <td>{formatRupiah(t.jumlah_bayar)}</td>
                                                        {isAdmin && <td>{t.staff_nama}</td>}
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    )}
                                </div>
                            </div>
                        </div>
                    )}
                </main>
            </div>
        </div>
    )
}
