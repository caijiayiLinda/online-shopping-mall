'use client';

import { notFound } from 'next/navigation';
import AddToCartButton from '@/components/AddToCartButton';
import Nav from '@/components/Nav';
import Image from 'next/image';
import Cart from '@/components/Cart';
import { Product } from '@/types';

import { useParams } from 'next/navigation';
import { useState, useEffect } from 'react';

export default function ProductPage() {
  const params = useParams();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        if (!params?.id) return;
        
        const response = await fetch(`/api/products/${params.id}`);
        if (!response.ok) {
          throw new Error('Failed to fetch product');
        }
        const data = await response.json();
        setProduct(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchProduct();
  }, [params?.id]);

  if (!params || !params.id) {
    notFound();
  }

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!product) {
    notFound();
  }

  return (
    <Nav categoryId={product.category_id} product={product}>
      <main className="max-w-6xl mx-auto py-8 px-4">
        <div className="grid md:grid-cols-2 gap-8">
          <div>
            <Image
              src={product.image_url}
              alt={product.name}
              width={500}
              height={500}
              className="w-full rounded-lg shadow-md"
            />
          </div>
          <div>
            <h1 className="text-3xl font-bold mb-4">{product.name}</h1>
            <p className="text-gray-600 text-xl mb-6">${product.price}</p>
            <p className="text-gray-700 mb-8">{product.description}</p>
            <AddToCartButton product={product} />
          </div>
        </div>
      </main>
      <Cart isOpen={false} onClose={function (): void {
        throw new Error('Function not implemented.');
      } } />
    </Nav>
  );
}
