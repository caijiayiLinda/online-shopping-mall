'use client';

import { useState } from 'react';
import axios from 'axios';

interface Order {
  id: string;
  date: string;
  items: {
    name: string;
    quantity: number;
    price: number;
  }[];
  total: number;
}

export default function OrderHistory() {
  const [isHovering, setIsHovering] = useState(false);
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const handleMouseEnter = async () => {
    setIsHovering(true);
    if (orders.length === 0 && !isLoading) {
      try {
        setIsLoading(true);
        const response = await axios.get('/api/orders/recent');
        setOrders(response.data);
      } catch (error) {
        console.error('Failed to fetch orders:', error);
      } finally {
        setIsLoading(false);
      }
    }
  };

  const handleMouseLeave = () => {
    setIsHovering(false);
  };

  return (
    <div className="relative">
      <button
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
        className="p-2 rounded-full hover:bg-gray-200"
        aria-label="Order history"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
          />
        </svg>
      </button>

      {isHovering && (
        <div className="absolute right-0 mt-2 w-64 bg-white rounded-md shadow-lg z-50 border border-gray-200">
          <div className="p-4">
            <h3 className="text-lg font-medium mb-2">Recent Orders</h3>
            {isLoading ? (
              <div className="text-center py-4">Loading...</div>
            ) : orders.length === 0 ? (
              <div className="text-center py-4">No recent orders</div>
            ) : (
              <div className="space-y-4">
                {orders.map((order) => (
                  <div key={order.id} className="border-b pb-2 last:border-b-0">
                    <div className="flex justify-between text-sm">
                      <span className="font-medium">Order #{order.id}</span>
                      <span className="text-gray-500">{new Date(order.date).toLocaleDateString()}</span>
                    </div>
                    <ul className="mt-1 space-y-1">
                      {order.items.map((item, index) => (
                        <li key={index} className="flex justify-between text-xs">
                          <span>{item.name} × {item.quantity}</span>
                          <span>${item.price.toFixed(2)}</span>
                        </li>
                      ))}
                    </ul>
                    <div className="mt-1 text-right text-sm font-medium">
                      Total: ${order.total.toFixed(2)}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
