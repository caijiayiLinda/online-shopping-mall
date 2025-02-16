'use client';

import { useState } from 'react';
import Image from 'next/image';
import { useCartContext } from '@/context/CartContext';

export default function Cart() {
  const [isOpen, setIsOpen] = useState(false);
  const { cartItems, updateQuantity } = useCartContext();

  const totalPrice = cartItems.reduce(
    (total, item) => total + item.price * item.quantity,
    0
  );

  return (
    <div
      className={`fixed inset-0 transition-opacity z-50 ${
        isOpen ? 'opacity-100 visible' : 'opacity-0 invisible'
      } bg-black bg-opacity-50`}
      onClick={() => setIsOpen(false)}
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
            onClick={() => setIsOpen(false)}
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

        <div className="p-4">
          {cartItems.length === 0 ? (
            <p className="text-gray-600">您的购物车是空的</p>
          ) : (
            <div className="space-y-4">
              {cartItems.map((item) => (
                <div key={item.id} className="flex items-center space-x-4">
                  <Image
                    src={item.image}
                    alt={item.name}
                    width={64}
                    height={64}
                    className="object-cover rounded"
                  />
                  <div className="flex-1">
                    <h3 className="font-medium">{item.name}</h3>
                    <p className="text-gray-600">${item.price}</p>
                  </div>
                  <div className="flex flex-col items-center">
                    <input
                      type="number"
                      min="1"
                      value={item.quantity}
                      onChange={(e) => {
                        const quantity = parseInt(e.target.value);
                        if (quantity > 0) {
                          updateQuantity(item.id, quantity);
                        }
                      }}
                      className="w-16 px-2 py-1 border rounded mb-2"
                      id={`quantity-input-${item.id}`}
                    />
                    <button
                      className="bg-blue-500 text-white py-1 px-2 rounded text-sm hover:bg-blue-600 transition-colors"
                      onClick={() => {
                        const quantityInput = document.querySelector(`#quantity-input-${item.id}`) as HTMLInputElement;
                        const quantity = parseInt(quantityInput.value);
                        if (quantity > 0) {
                          updateQuantity(item.id, quantity);
                        }
                      }}
                    >
                      更新
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="absolute bottom-0 left-0 right-0 p-4 border-t bg-white">
          <div className="flex justify-between items-center mb-4">
            <span className="font-medium">总计：</span>
            <span className="text-xl font-semibold">${totalPrice.toFixed(2)}</span>
          </div>
          <button
            className="w-full bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition-colors mb-2"
            onClick={async () => {
              try {
                const response = await fetch('/api/checkout', {
                  method: 'POST',
                  headers: {
                    'Content-Type': 'application/json',
                  },
                  body: JSON.stringify({
                    items: cartItems,
                    total: totalPrice.toFixed(2),
                  }),
                });

                const { url } = await response.json();
                window.location.href = url;
              } catch (error) {
                console.error('Checkout error:', error);
                alert('支付处理失败，请稍后重试');
              }
            }}
          >
            结算
          </button>
          <button
            className="w-full text-blue-500 py-2 px-4 rounded hover:text-blue-600 transition-colors border border-blue-500"
          >
            继续购物
          </button>
        </div>
      </div>
    </div>
  );
}
