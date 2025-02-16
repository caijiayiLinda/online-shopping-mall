'use client';

import Cart from '@/components/Cart';
import { CartProvider } from '@/context/CartContext';

export default function ClientLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <CartProvider>
      <div className="min-h-screen flex flex-col">
        <main className="flex-1">{children}</main>
      </div>
      <Cart />
    </CartProvider>
  );
}
