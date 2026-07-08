import { useCallback } from 'react'
import axiosClient from '../api/axiosClient'
import DashboardLayout from '../components/DashboardLayout'

export default function DashboardAdmin() {
  const fetchData = useCallback(async () => {
    const res = await axiosClient.get('/admin/dashboard')
    return res.data.data
  }, [])

  return (
    <DashboardLayout
      title="Dashboard Admin"
      subtitle="Ringkasan pengelolaan petugas, guru, siswa, kelas, dan pembayaran SPP."
      fetchData={fetchData}
    />
  )
}
