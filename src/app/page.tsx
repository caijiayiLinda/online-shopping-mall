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
  const selectedCategory = searchParams.get('category') || undefined;

  return (
    <Suspense fallback={<Loading />}>
      <Nav category={selectedCategory}>
        <ProductList category={selectedCategory} />
        <Cart />
      </Nav>
    </Suspense>
  );
}
