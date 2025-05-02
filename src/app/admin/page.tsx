'use client';

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import AdminForm from '@/components/AdminForm'

interface Order {
  id: number;
  invoice: string;
  username: string;
  total_price: number;
  created_at: string;
  status: string;
  products: {
    product_id: number;
    quantity: number;
    price: number;
    name?: string;
  }[];
}

function OrderHistoryTable() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const response = await fetch('/api/admin/orders');
        const data = await response.json();
        setOrders(data);
      } catch (error) {
        console.error('Failed to fetch orders:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchOrders();
  }, []);

  if (loading) {
    return <div>Loading orders...</div>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full bg-white">
        <thead>
          <tr>
            <th className="py-2 px-4 border-b border-gray-200 bg-gray-50 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
              Order ID
            </th>
            <th className="py-2 px-4 border-b border-gray-200 bg-gray-50 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
              Invoice
            </th>
            <th className="py-2 px-4 border-b border-gray-200 bg-gray-50 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
              Customer
            </th>
            <th className="py-2 px-4 border-b border-gray-200 bg-gray-50 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
              Date
            </th>
            <th className="py-2 px-4 border-b border-gray-200 bg-gray-50 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
              Total
            </th>
            <th className="py-2 px-4 border-b border-gray-200 bg-gray-50 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
              Status
            </th>
          </tr>
        </thead>
        <tbody>
          {orders.map((order) => (
            <tr key={order.id}>
              <td className="py-2 px-4 border-b border-gray-200">
                {order.id}
              </td>
              <td className="py-2 px-4 border-b border-gray-200">
                {order.invoice}
              </td>
              <td className="py-2 px-4 border-b border-gray-200">
                {order.username}
              </td>
              <td className="py-2 px-4 border-b border-gray-200">
                {new Date(order.created_at).toLocaleDateString()}
              </td>
              <td className="py-2 px-4 border-b border-gray-200 font-medium">
                ${order.total_price.toFixed(2)}
              </td>
              <td className="py-2 px-4 border-b border-gray-200">
                {order.status}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default function AdminPage() {
  const { isAuthenticated, isAdmin, checkAuth } = useAuth()
  const router = useRouter()
  const [authChecking, setAuthChecking] = useState(true)

  useEffect(() => {
    let mounted = true
    const verifyAdmin = async () => {
      await checkAuth()
      if (mounted) {
        setAuthChecking(false)
        if (!isAuthenticated || !isAdmin) {
          router.replace('/login?from=/admin')
        }
      }
    }
    verifyAdmin()
    return () => { mounted = false }
  }, [checkAuth, isAuthenticated, isAdmin, router])

  if (authChecking || !isAuthenticated || !isAdmin) {
    return null
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Admin Panel</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div>
          <h2 className="text-xl font-semibold mb-4">Add New Product</h2>
          <AdminForm />
        </div>
        <div>
          <h2 className="text-xl font-semibold mb-4">Order Management</h2>
          <OrderHistoryTable />
        </div>
      </div>
    </div>
  )
}
