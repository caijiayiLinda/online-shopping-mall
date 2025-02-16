'use client';

import Link from 'next/link';
import Image from 'next/image';
import AddToCartButton from './AddToCartButton';
import { products } from '@/data/products';

interface ProductListProps {
  category?: string;
}

export default function ProductList({ category }: ProductListProps) {
  const filteredProducts = category
    ? products.filter((p) => p.category === category)
    : products;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      {filteredProducts.map((product) => (
        <div
          key={product.id}
          className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow flex flex-col"
        >
          <Link href={`/products/${product.id}`}>
            <div className="cursor-pointer flex flex-col">
              <div className="w-full h-48 relative">
                <Image
                  src={product.image}
                  alt={product.name}
                  fill
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
