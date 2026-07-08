import { useCallback } from 'react'
import axiosClient from '../api/axiosClient'
import DashboardLayout from '../components/DashboardLayout'

export default function DashboardSiswa() {
  const fetchData = useCallback(async () => {
    const res = await axiosClient.get('/siswa/dashboard')
    return res.data.data
  }, [])

  return (
    <DashboardLayout
      title="Dashboard Siswa"
      subtitle="Lihat status pembayaran SPP dan riwayat transaksi kamu di sini."
      fetchData={fetchData}
    />
  )
}
