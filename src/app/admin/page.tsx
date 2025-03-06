'use client'
import AdminForm from '@/components/AdminForm'

export default function AdminPage() {
  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">产品管理</h1>
      <div>
        <h2 className="text-xl font-semibold mb-4">添加新产品</h2>
        <AdminForm />
      </div>
    </div>
  )
}
