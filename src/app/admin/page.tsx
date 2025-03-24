'use client';

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import AdminForm from '@/components/AdminForm'

export default function AdminPage() {
  const { isAuthenticated, isAdmin, checkAuth } = useAuth()
  const router = useRouter()

  console.log('AdminPage auth state:', { isAuthenticated, isAdmin })

  const [initialCheckDone, setInitialCheckDone] = useState(false)

  useEffect(() => {
    console.log('AdminPage useEffect triggered')
    checkAuth().then(() => {
      console.log('AdminPage checkAuth completed', { isAuthenticated, isAdmin })
      setInitialCheckDone(true)
    })
  }, [checkAuth])

  useEffect(() => {
    if (!initialCheckDone) return
    
    if (!isAuthenticated) {
      console.log('AdminPage redirecting to login')
      router.push('/login?from=/admin')
    }
  }, [isAuthenticated, initialCheckDone, router])

  if (!isAuthenticated || !isAdmin) {
    console.log('AdminPage rendering null (not authenticated or not admin)')
    return null // 防止页面闪烁
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Admin Panel</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div>
          <h2 className="text-xl font-semibold mb-4">Add New Product</h2>
          <AdminForm />
        </div>
        
      </div>
    </div>
  )
}
