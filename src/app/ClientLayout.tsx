'use client';

import { Suspense } from 'react';
import Cart from '@/components/Cart';
import { CartProvider } from '@/context/CartContext';

function Loading() {
  return <div>Loading...</div>;
}

export default function ClientLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <CartProvider>
      <div className="min-h-screen flex flex-col">
        <Suspense fallback={<Loading />}>
          <main className="flex-1">{children}</main>
        </Suspense>
      </div>
      <Cart isOpen={false} onClose={function (): void {
        throw new Error('Function not implemented.');
      } } />
    </CartProvider>
  );
}
