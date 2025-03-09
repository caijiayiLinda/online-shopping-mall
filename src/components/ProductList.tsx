'use client';

import Link from 'next/link';
import Image from 'next/image';
import AddToCartButton from './AddToCartButton';
import { useState, useEffect } from 'react';
import { Product } from '@/types';

interface ProductListProps {
  categoryId?: string;
}

export default function ProductList({ categoryId }: ProductListProps) {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const url = categoryId 
          ? `/api/products/category?category_id=${categoryId}`
          : '/api/products';
        const response = await fetch(url);
        if (!response.ok) {
          throw new Error('Failed to fetch products');
        }
        const data = await response.json();
        setProducts(data);
      } catch (error) {
        console.error('Error fetching products:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchProducts();
  }, [categoryId]);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      {products.map((product) => (
        <div
          key={product.id}
          className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow flex flex-col"
        >
          <Link href={`/products/${product.id}`}>
            <div className="cursor-pointer flex flex-col">
              <div className="w-[300px] h-[300px] relative">
                <Image
                  src={product.thumbnail_url}
                  alt={product.name}
                  width={300}
                  height={300}
                  className="object-cover"
                />
              </div>
              <div className="p-4 flex-1 flex flex-col">
                <h2 className="text-lg font-semibold mb-2 overflow-hidden text-ellipsis whitespace-nowrap">{product.name}</h2>
                <p className="text-gray-600 mb-4 overflow-hidden text-ellipsis whitespace-nowrap">${product.price}</p>
              </div>
            </div>
          </Link>
          <div className="p-4 border-t">
            <AddToCartButton product={product} />
          </div>
        </div>
      ))}
    </div>
  );
}
