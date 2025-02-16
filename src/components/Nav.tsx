'use client';

import Link from 'next/link';
import CartButton from './CartButton';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { categories } from '@/data/products';

interface NavProps {
  category?: string;
  product?: string;
  children: React.ReactNode;
}

export default function Nav({ category, product, children }: NavProps) {
  const router = useRouter();
  const [selectedCategory, setSelectedCategory] = useState(category || '');

  const handleCategoryClick = (category: string) => {
    setSelectedCategory(category);
    router.push(`/?category=${category}`);
  };

  const uniqueCategories = [...new Set(categories)];

  return (
    <div className="bg-white">
      <div className="max-w-6xl mx-auto px-4">
        {/* 顶部导航栏 */}
        <div className="flex justify-between items-center h-16">
          <Link href="/" className="text-xl font-bold text-gray-800">
            Online shopping mall
          </Link>
          <div className="hidden md:flex space-x-4"> {/* Horizontal category navigation for larger screens */}
                <button
                  onClick={() => handleCategoryClick('')}
                  className={`px-3 py-2 text-sm rounded-md ${
                    selectedCategory === ''
                      ? 'bg-blue-500 text-white'
                      : 'text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  All
                </button>
            {uniqueCategories.map((category) => (
              <button
                key={category}
                onClick={() => handleCategoryClick(category)}
                className={`px-3 py-2 text-sm rounded-md ${
                  selectedCategory === category
                    ? 'bg-blue-500 text-white'
                    : 'text-gray-700 hover:bg-gray-200'
                }`}
              >
                {category}
              </button>
            ))}
          </div>
          <CartButton />
        </div>
        {/* Breadcrumbs */}
        <div className="py-2">
          <Link href="/" className="text-blue-500 hover:underline">Home</Link>
          {category && category !== '' && (
            <span className="inline-flex items-center">
              <span className="mx-2">{'>'}</span>
              <Link href={`/?category=${category}`} className="text-blue-500 hover:underline">{category}</Link>
              {product && product !== '' && (
                <span className="inline-flex items-center">
                  <span className="mx-2">{'>'}</span>
                  <span>{product}</span>
                </span>
              )}
            </span>
          )}
        </div>
        {/* Content area */}
        <div>
          {children}
        </div>
      </div>
    </div>
  );
}
