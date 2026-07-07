import { createContext, useContext, useState } from 'react'
import axiosClient from '../api/axiosClient'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => {
    const stored = localStorage.getItem('user')
    return stored ? JSON.parse(stored) : null
  })
  const [loading, setLoading] = useState(false)

  // login memanggil POST /auth/login, menyimpan token & user ke localStorage + state
  const login = async (email, password) => {
    setLoading(true)
    try {
      const res = await axiosClient.post('/auth/login', { email, password })
      const { token, user: loggedInUser } = res.data.data

      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(loggedInUser))
      setUser(loggedInUser)

      return { success: true, user: loggedInUser }
    } catch (err) {
      const message = err.response?.data?.message || 'Gagal login, silakan coba lagi'
      return { success: false, message }
    } finally {
      setLoading(false)
    }
  }

  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, login, logout, loading }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth harus dipakai di dalam AuthProvider')
  }
  return context
}
