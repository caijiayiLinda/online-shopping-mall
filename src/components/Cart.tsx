'use client';

import { useState, useEffect } from 'react';
import Image from 'next/image';
import { useCartContext } from '@/context/CartContext';
import { Product } from '@/types';

interface CartProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function Cart({ isOpen, onClose }: CartProps) {
  const { cartItems, updateQuantity, removeFromCart /*, user*/ } = useCartContext();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchProducts = async () => {
      if (cartItems.length === 0) return;
      
      setLoading(true);
      try {
        const productPromises = cartItems.map(async (item) => {
          const response = await fetch(`/api/products/${item.id}`);
          if (!response.ok) {
            throw new Error('Failed to fetch product');
          }
          return response.json();
        });

        const productsData = await Promise.all(productPromises);
        setProducts(productsData);
      } catch (error) {
        console.error('Error fetching products:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, [cartItems]);

  const totalPrice = cartItems.reduce(
    (total, item) => total + item.price * item.quantity,
    0
  );

  return (
    <div
      className={`fixed inset-0 transition-opacity z-50 ${
        isOpen ? 'opacity-100 visible' : 'opacity-0 invisible'
      } bg-black bg-opacity-50`}
      onClick={onClose}
    >
      <div
        className={`absolute right-0 top-0 h-full w-96 bg-white shadow-lg transform transition-transform ${
          isOpen ? 'translate-x-0' : 'translate-x-full'
        } hover:shadow-xl transition-shadow duration-200 z-50`}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex justify-between items-center p-4 border-b z-50 relative z-50">
          <h2 className="text-xl font-semibold">Shopping Cart</h2>
          <button
            onClick={onClose}
            className="text-gray-600 hover:text-gray-900 z-50 "
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
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <div className="p-4 h-[calc(100vh-200px)] overflow-y-auto">
          {cartItems.length === 0 ? (
            <p className="text-gray-600">Empty Shopping Cart</p>
          ) : (
            <div className="space-y-4 pb-4">
              {cartItems.map((item) => (
                <div key={item.id} className="flex items-center space-x-4">
                  {loading ? (
                    <div className="w-16 h-16 bg-gray-200 animate-pulse rounded" />
                  ) : (
                    <Image
                      src={products.find(p => p.id === item.id)?.thumbnail_url || '/placeholder.jpg'}
                      alt={item.name}
                      width={64}
                      height={64}
                      className="object-cover rounded"
                    />
                  )}
                  <div className="flex-1">
                    <h3 className="font-medium">{item.name}</h3>
                    <p className="text-gray-600">${item.price}</p>
                  </div>
                  <div className="flex flex-col items-center space-y-2 z-50 relative">
                    <div className="flex items-center space-x-2 z-50">
                      <button
                        className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center hover:bg-gray-300 transition-colors"
                        onClick={() => updateQuantity(item.id, item.quantity - 1)}
                      >
                        -
                      </button>
                      <input
                        type="number"
                        min="0"
                        value={item.quantity}
                        onChange={(e) => {
                          const quantity = parseInt(e.target.value);
                          if (quantity > 0) {
                            updateQuantity(item.id, quantity);
                          }
                        }}
                        className="w-16 px-2 py-1 border rounded text-center"
                        id={`quantity-input-${item.id}`}
                      />
                      <button
                        className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center hover:bg-gray-300 transition-colors"
                        onClick={() => updateQuantity(item.id, item.quantity + 1)}
                      >
                        +
                      </button>
                    </div>
                    <button
                      className="text-red-500 hover:text-red-700 text-sm"
                      onClick={() => removeFromCart(item.id)}
                    >
                      delete
                    </button>
                    {/* <p className="text-xs text-gray-500">ID: {item.id}</p> */}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="fixed bottom-0 left-0 right-0 p-4 border-t bg-white" style={{ width: '384px' }}>
          <div className="flex justify-between items-center mb-4">
            <span className="font-medium">总计：</span>
            <span className="text-xl font-semibold">${totalPrice.toFixed(2)}</span>
          </div>
          <form onSubmit={async (e) => {
            e.preventDefault();
            
            // if (!user?.email) {
            //   alert('请先登录');
            //   return;
            // }

            try {
              const response = await fetch('/api/checkout/paypal', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                  cartItems: cartItems.map(item => ({
                    id: item.id,
                    name: item.name,
                    price: item.price,
                    quantity: item.quantity
                  })),
                  invoice: `INV-${Date.now()}`,
                  // email: user.email
                })
              });

              const data = await response.json();
              window.location.href = data.approvalUrl;
            } catch (error) {
              console.error('结账失败:', error);
              alert('结账失败，请检查控制台');
            }
          }}>
            <button
              type="submit"
              className="w-full bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition-colors mb-2"
            >
              Checkout with PayPal
            </button>
          </form>
          <button
            className="w-full text-blue-500 py-2 px-4 rounded hover:text-blue-600 transition-colors border border-blue-500"
          >
            back
          </button>
        </div>
      </div>
    </div>
  );
}
