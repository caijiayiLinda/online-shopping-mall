'use client';

'use client';

import { useSearchParams } from 'next/navigation';
import ProductList from '@/components/ProductList';
import Nav from '@/components/Nav';
import Cart from '@/components/Cart';

export default function Home() {
  const searchParams = useSearchParams();
  const selectedCategory = searchParams.get('category') || undefined;

  return (
    <Nav category={selectedCategory}>
      <ProductList category={selectedCategory} />
      <Cart />
    </Nav>
  );
}
