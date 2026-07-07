import { useCallback } from 'react'
import axiosClient from '../api/axiosClient'
import DashboardLayout from '../components/DashboardLayout'

export default function DashboardGuru() {
  const fetchData = useCallback(async () => {
    const res = await axiosClient.get('/guru/dashboard')
    return res.data.data
  }, [])

  return (
    <DashboardLayout
      title="Dashboard Guru"
      subtitle="Ringkasan kelas, tugas, dan presensi yang kamu ampu."
      fetchData={fetchData}
    />
  )
}
