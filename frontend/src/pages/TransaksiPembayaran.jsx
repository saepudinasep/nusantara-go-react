import { useCallback, useEffect, useState } from 'react'
import Swal from 'sweetalert2'
import axiosClient from '../api/axiosClient'
import { useAuth } from '../context/AuthContext'
import Sidebar from '../components/Sidebar'
import Topbar from '../components/Topbar'

const PAGE_SIZE = 10
const DEFAULT_PAGINATION = { page: 1, limit: PAGE_SIZE, total: 0, total_pages: 1 }
const BULAN_OPTIONS = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember',
]

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

function todayISO() {
    return new Date().toISOString().slice(0, 10)
}

export default function TransaksiPembayaran() {
    const { user } = useAuth()
    const isAdmin = user.role === 'admin'
    const basePath = isAdmin ? '/admin' : '/petugas'

    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [list, setList] = useState([])
    const [loading, setLoading] = useState(true)
    const [page, setPage] = useState(1)
    const [pagination, setPagination] = useState(DEFAULT_PAGINATION)

    const [sppOptions, setSppOptions] = useState([])
    const [staffOptions, setStaffOptions] = useState([]) // hanya dipakai admin

    const [modalOpen, setModalOpen] = useState(false)
    const [formError, setFormError] = useState('')
    const [saving, setSaving] = useState(false)
    const [searchingNisn, setSearchingNisn] = useState(false)

    const [nisnInput, setNisnInput] = useState('')
    const [foundStudent, setFoundStudent] = useState(null)
    const [form, setForm] = useState({
        spp_id: '',
        staff_id: '',
        bulan_dibayar: '',
        tanggal_bayar: todayISO(),
        jumlah_bayar: '',
    })

    // Data pendukung form: daftar SPP (admin & petugas, read-only), daftar petugas (admin saja)
    useEffect(() => {
        axiosClient
            .get(`${basePath}/spp`, { params: { page: 1, limit: 100 } })
            .then((res) => setSppOptions(res.data.data?.items || []))
            .catch(() => setSppOptions([]))

        if (isAdmin) {
            axiosClient
                .get('/admin/petugas', { params: { page: 1, limit: 100 } })
                .then((res) => setStaffOptions(res.data.data?.items || []))
                .catch(() => setStaffOptions([]))
        }
    }, [basePath, isAdmin])

    const loadTransaksi = useCallback(
        async (targetPage) => {
            setLoading(true)
            try {
                const res = await axiosClient.get(`${basePath}/transaksi`, { params: { page: targetPage, limit: PAGE_SIZE } })
                const payload = res.data.data
                setList(payload?.items || [])
                setPagination(payload?.pagination || DEFAULT_PAGINATION)
            } catch (err) {
                Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal memuat data transaksi' })
            } finally {
                setLoading(false)
            }
        },
        [basePath]
    )

    useEffect(() => {
        loadTransaksi(page)
    }, [page, loadTransaksi])

    const openCreateModal = () => {
        setNisnInput('')
        setFoundStudent(null)
        setForm({ spp_id: '', staff_id: '', bulan_dibayar: '', tanggal_bayar: todayISO(), jumlah_bayar: '' })
        setFormError('')
        setModalOpen(true)
    }

    const closeModal = () => {
        if (saving) return
        setModalOpen(false)
    }

    const handleSearchNisn = async (e) => {
        e.preventDefault()
        setFormError('')
        setFoundStudent(null)
        setSearchingNisn(true)
        try {
            const res = await axiosClient.get(`${basePath}/siswa/cari`, { params: { nisn: nisnInput.trim() } })
            setFoundStudent(res.data.data)
        } catch (err) {
            setFormError(err.response?.data?.message || 'Siswa dengan NISN tersebut tidak ditemukan')
        } finally {
            setSearchingNisn(false)
        }
    }

    const handleSppChange = (sppId) => {
        const selected = sppOptions.find((s) => String(s.id) === String(sppId))
        setForm({ ...form, spp_id: sppId, jumlah_bayar: selected ? selected.nominal : '' })
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        if (!foundStudent) {
            setFormError('Cari siswa berdasarkan NISN terlebih dahulu')
            return
        }
        setFormError('')
        setSaving(true)

        try {
            const payload = {
                student_id: foundStudent.id,
                spp_id: Number(form.spp_id),
                bulan_dibayar: form.bulan_dibayar,
                tanggal_bayar: form.tanggal_bayar,
                jumlah_bayar: Number(form.jumlah_bayar),
            }
            if (isAdmin) {
                payload.staff_id = Number(form.staff_id)
            }

            const res = await axiosClient.post(`${basePath}/transaksi`, payload)
            setModalOpen(false)
            loadTransaksi(1)
            setPage(1)
            showReceipt(res.data.data)
        } catch (err) {
            setFormError(err.response?.data?.message || 'Gagal memproses pembayaran')
        } finally {
            setSaving(false)
        }
    }

    const showReceipt = (payment) => {
        Swal.fire({
            title: 'Pembayaran Berhasil',
            html: `
        <div id="struk-cetak" style="text-align:left; font-size:13px; line-height:1.8;">
          <p><strong>Nama Siswa:</strong> ${payment.student_nama}</p>
          <p><strong>NISN:</strong> ${payment.nisn}</p>
          <p><strong>Kelas:</strong> ${payment.nama_kelas}</p>
          <p><strong>Tahun Ajaran:</strong> ${payment.tahun_ajaran}</p>
          <p><strong>Bulan Dibayar:</strong> ${payment.bulan_dibayar}</p>
          <p><strong>Tanggal:</strong> ${payment.tanggal_bayar}</p>
          <p><strong>Jumlah:</strong> ${formatRupiah(payment.jumlah_bayar)}</p>
          <p><strong>Diproses oleh:</strong> ${payment.staff_nama}</p>
        </div>
      `,
            icon: 'success',
            showCancelButton: true,
            confirmButtonText: 'Cetak Struk',
            cancelButtonText: 'Tutup',
        }).then((result) => {
            if (result.isConfirmed) {
                window.print()
            }
        })
    }

    const handleCancel = async (payment) => {
        const result = await Swal.fire({
            icon: 'warning',
            title: `Batalkan transaksi ${payment.student_nama}?`,
            text: `Pembayaran ${formatRupiah(payment.jumlah_bayar)} untuk bulan ${payment.bulan_dibayar} akan dihapus permanen.`,
            showCancelButton: true,
            confirmButtonText: 'Ya, batalkan',
            cancelButtonText: 'Kembali',
            confirmButtonColor: '#ef4444',
            cancelButtonColor: '#9aa3b8',
            reverseButtons: true,
        })

        if (!result.isConfirmed) return

        try {
            await axiosClient.delete(`/admin/transaksi/${payment.id}`)
            Toast.fire({ icon: 'success', title: 'Transaksi berhasil dibatalkan' })
            loadTransaksi(page)
        } catch (err) {
            Toast.fire({ icon: 'error', title: err.response?.data?.message || 'Gagal membatalkan transaksi' })
        }
    }

    return (
        <div className="app-shell">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />

            <div className="main-shell">
                <Topbar title="Transaksi Pembayaran" onMenuClick={() => setSidebarOpen((prev) => !prev)} />

                <main className="main-content">
                    <div className="breadcrumb">
                        SchoolApp · {isAdmin ? 'Admin' : 'Petugas'} · <span>Transaksi Pembayaran</span>
                    </div>
                    <div className="page-header">
                        <h1>Transaksi Pembayaran</h1>
                        <p>
                            {isAdmin
                                ? 'Lihat seluruh transaksi pembayaran SPP dan proses pembayaran baru.'
                                : 'Proses pembayaran SPP siswa dan lihat riwayat transaksi yang kamu proses sendiri.'}
                        </p>
                    </div>

                    <div className="card">
                        <div className="toolbar">
                            <span style={{ fontSize: 13, color: 'var(--text2)' }}>
                                {pagination.total} transaksi {isAdmin ? 'tercatat' : 'kamu proses'}
                            </span>
                            <button type="button" className="btn btn-primary" onClick={openCreateModal}>
                                + Proses Pembayaran
                            </button>
                        </div>

                        <div className="table-wrap">
                            {loading ? (
                                <div className="empty-state">Memuat data...</div>
                            ) : list.length === 0 ? (
                                <div className="empty-state">Belum ada transaksi.</div>
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
                                            {isAdmin && <th style={{ textAlign: 'right' }}>Aksi</th>}
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {list.map((p) => (
                                            <tr key={p.id}>
                                                <td>{p.tanggal_bayar}</td>
                                                <td style={{ fontWeight: 600 }}>
                                                    {p.student_nama}
                                                    <div style={{ fontSize: 11, color: 'var(--gray3)', fontWeight: 400 }}>{p.nisn}</div>
                                                </td>
                                                <td>
                                                    <span className="chip">{p.nama_kelas}</span>
                                                </td>
                                                <td>{p.bulan_dibayar}</td>
                                                <td>{formatRupiah(p.jumlah_bayar)}</td>
                                                {isAdmin && <td>{p.staff_nama}</td>}
                                                {isAdmin && (
                                                    <td>
                                                        <div className="table-actions" style={{ justifyContent: 'flex-end' }}>
                                                            <button type="button" className="btn btn-danger-ghost" onClick={() => handleCancel(p)}>
                                                                Batalkan
                                                            </button>
                                                        </div>
                                                    </td>
                                                )}
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </div>

                        {!loading && pagination.total > 0 && (
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

            {modalOpen && (
                <div className="modal-overlay" onClick={closeModal}>
                    <div className="modal-box" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h3>Proses Pembayaran SPP</h3>
                            <button type="button" className="modal-close" onClick={closeModal}>
                                ✕
                            </button>
                        </div>

                        <div className="modal-body">
                            {formError && <div className="form-error">{formError}</div>}

                            {/* Langkah 1: cari siswa by NISN */}
                            <form onSubmit={handleSearchNisn} style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
                                <input
                                    type="text"
                                    placeholder="Masukkan NISN siswa"
                                    value={nisnInput}
                                    onChange={(e) => setNisnInput(e.target.value)}
                                    style={{
                                        flex: 1,
                                        padding: '10px 12px',
                                        border: '1px solid var(--gray2)',
                                        borderRadius: 8,
                                        fontSize: 13,
                                    }}
                                    required
                                />
                                <button type="submit" className="btn btn-outline" disabled={searchingNisn}>
                                    {searchingNisn ? 'Mencari...' : 'Cari'}
                                </button>
                            </form>

                            {foundStudent && (
                                <div
                                    style={{
                                        background: 'var(--gray0)',
                                        borderRadius: 8,
                                        padding: 12,
                                        marginBottom: 16,
                                        fontSize: 13,
                                    }}
                                >
                                    <strong>{foundStudent.nama}</strong> — {foundStudent.nama_kelas} (Tingkat {foundStudent.tingkat})
                                </div>
                            )}

                            {/* Langkah 2: detail pembayaran, hanya aktif setelah siswa ditemukan */}
                            {foundStudent && (
                                <form onSubmit={handleSubmit}>
                                    <div className="form-field">
                                        <label htmlFor="spp_id">Jenis SPP</label>
                                        <select
                                            id="spp_id"
                                            value={form.spp_id}
                                            onChange={(e) => handleSppChange(e.target.value)}
                                            required
                                        >
                                            <option value="" disabled>
                                                Pilih SPP
                                            </option>
                                            {sppOptions.map((s) => (
                                                <option key={s.id} value={s.id}>
                                                    {s.tahun_ajaran} — {formatRupiah(s.nominal)}
                                                </option>
                                            ))}
                                        </select>
                                    </div>

                                    {isAdmin && (
                                        <div className="form-field">
                                            <label htmlFor="staff_id">Diproses oleh Petugas</label>
                                            <select
                                                id="staff_id"
                                                value={form.staff_id}
                                                onChange={(e) => setForm({ ...form, staff_id: e.target.value })}
                                                required
                                            >
                                                <option value="" disabled>
                                                    Pilih petugas
                                                </option>
                                                {staffOptions.map((s) => (
                                                    <option key={s.id} value={s.id}>
                                                        {s.nama} ({s.posisi})
                                                    </option>
                                                ))}
                                            </select>
                                        </div>
                                    )}

                                    <div className="form-field">
                                        <label htmlFor="bulan_dibayar">Bulan Dibayar</label>
                                        <select
                                            id="bulan_dibayar"
                                            value={form.bulan_dibayar}
                                            onChange={(e) => setForm({ ...form, bulan_dibayar: e.target.value })}
                                            required
                                        >
                                            <option value="" disabled>
                                                Pilih bulan
                                            </option>
                                            {BULAN_OPTIONS.map((b) => (
                                                <option key={b} value={b}>
                                                    {b}
                                                </option>
                                            ))}
                                        </select>
                                    </div>

                                    <div className="form-field">
                                        <label htmlFor="tanggal_bayar">Tanggal Bayar</label>
                                        <input
                                            id="tanggal_bayar"
                                            type="date"
                                            value={form.tanggal_bayar}
                                            onChange={(e) => setForm({ ...form, tanggal_bayar: e.target.value })}
                                            required
                                        />
                                    </div>

                                    <div className="form-field">
                                        <label htmlFor="jumlah_bayar">Jumlah Bayar (Rp)</label>
                                        <input
                                            id="jumlah_bayar"
                                            type="text"
                                            // min="1"
                                            // step="1000"
                                            value={form.jumlah_bayar}
                                            onChange={(e) => setForm({ ...form, jumlah_bayar: e.target.value })}
                                            required
                                        />
                                    </div>

                                    <div className="modal-footer" style={{ padding: '16px 0 0', border: 'none' }}>
                                        <button type="button" className="btn btn-outline" onClick={closeModal} disabled={saving}>
                                            Batal
                                        </button>
                                        <button type="submit" className="btn btn-primary" disabled={saving}>
                                            {saving ? 'Memproses...' : 'Terima & Bayar'}
                                        </button>
                                    </div>
                                </form>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}
