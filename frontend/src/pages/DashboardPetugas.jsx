import { useCallback } from 'react'
import axiosClient from '../api/axiosClient'
import DashboardLayout from '../components/DashboardLayout'

export default function DashboardPetugas() {
    const fetchData = useCallback(async () => {
        const res = await axiosClient.get('/petugas/dashboard')
        return res.data.data
    }, [])

    return (
        <DashboardLayout
            title="Dashboard Petugas"
            subtitle="Ringkasan transaksi pembayaran SPP dan data siswa."
            fetchData={fetchData}
        />
    )
}
