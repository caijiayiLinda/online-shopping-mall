'use client'
import AdminForm from '@/components/AdminForm'

export default function AdminPage() {
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
