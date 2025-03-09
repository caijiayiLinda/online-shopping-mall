'use client';

import { Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import ProductList from '@/components/ProductList';
import Nav from '@/components/Nav';
import Cart from '@/components/Cart';

function Loading() {
  return <div>Loading...</div>;
}

export default function Home() {
  const searchParams = useSearchParams()!;
  const categoryId = searchParams.get('categoryId') || undefined;

  return (
    <Suspense fallback={<Loading />}>
      <Nav categoryId={categoryId}>
        <ProductList categoryId={categoryId} />
        <Cart isOpen={false} onClose={function (): void {
          throw new Error('Function not implemented.');
        } } />
      </Nav>
    </Suspense>
  );
}
