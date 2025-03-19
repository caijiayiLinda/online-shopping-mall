'use client'
import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import AdminForm from '@/components/AdminForm'

export default function AdminPage() {
  const { isAuthenticated, isAdmin, checkAuth } = useAuth()
  const router = useRouter()

  useEffect(() => {
    checkAuth().then(() => {
      if (!isAuthenticated) {
        router.push('/login')
      } else if (!isAdmin) {
        alert('当前用户无管理权限')
        router.push('/login')
      }
    })
  }, [isAuthenticated, isAdmin, checkAuth, router])

  if (!isAuthenticated || !isAdmin) {
    return null // 防止页面闪烁
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Manage Products</h1>
      <div>
        <h2 className="text-xl font-semibold mb-4">Add New Product</h2>
        <AdminForm />
      </div>
    </div>
  )
}
