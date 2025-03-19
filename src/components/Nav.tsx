'use client';

import Link from 'next/link';
import { Product } from '@/types';
import CartButton from './CartButton';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import axios from 'axios';
import { useAuth } from '@/hooks/useAuth';

const useCategories = () => {
  const [categories, setCategories] = useState<string[]>([]);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        interface Category {
          id: number;
          name: string;
        }
        const response = await axios.get<Category[]>('/api/categories');
        const categoryNames = response.data.map((c: Category) => c.name);
        setCategories(categoryNames);
      } catch (error) {
        console.error('Failed to fetch categories:', error);
      }
    };
    fetchCategories();
  }, []);

  return categories;
};

interface NavProps {
  categoryId?: string;
  product?: Product;
  children: React.ReactNode;
}

function AuthStatus() {
  const { isAuthenticated, isAdmin, logout } = useAuth();

  return (
    <div className="flex items-center gap-2">
      {isAuthenticated ? (
        <>
          <span className="text-sm">Welcome, {isAdmin ? 'Admin' : 'User'}</span>
          <button
            onClick={logout}
            className="px-3 py-1 text-sm text-white bg-red-500 rounded hover:bg-red-600"
          >
            Logout
          </button>
        </>
      ) : (
        <Link
          href="/login"
          className="px-3 py-1 text-sm text-white bg-blue-500 rounded hover:bg-blue-600"
        >
          Login
        </Link>
      )}
    </div>
  );
}

export default function Nav({ categoryId, product, children }: NavProps) {
  const router = useRouter();
  const [selectedCategory, setSelectedCategory] = useState(categoryId || '');

  const categoryMap = {
    'Clothing': 1,
    'Tools': 2,
    'Toys': 3,
    'Beauty': 4,
    'Pets': 5
  };

  const handleCategoryClick = (category: string) => {
    setSelectedCategory(category);
    if (category) {
      const categoryId = categoryMap[category as keyof typeof categoryMap];
      router.push(`/?categoryId=${categoryId}`);
    } else {
      router.push('/');
    }
  };

  const categories = useCategories();
  const uniqueCategories = [...new Set(categories)];

  return (
    <div className="bg-white">
      <div className="max-w-6xl mx-auto px-4">
        {/* 顶部导航栏 */}
        <div className="flex justify-between items-center h-16">
          <Link href="/" className="text-xl font-bold text-gray-800">
            Online shopping mall
          </Link>
          <div className="hidden md:flex space-x-4">
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
              <Link
                key={category}
                href={`/?categoryId=${categoryMap[category as keyof typeof categoryMap]}`}
                className={`px-3 py-2 text-sm rounded-md ${
                  selectedCategory === category
                    ? 'bg-blue-500 text-white'
                    : 'text-gray-700 hover:bg-gray-200'
                }`}
              >
                {category}
              </Link>
            ))}
          </div>
          <div className="flex items-center gap-4">
            <CartButton />
            <AuthStatus />
          </div>
        </div>
        {/* Breadcrumbs */}
        <div className="py-2">
          <Link href="/" className="text-blue-500 hover:underline">Home</Link>
          {(categoryId || (typeof product === 'object' && product?.category_id)) && (
            <span className="inline-flex items-center">
              <span className="mx-2">{'>'}</span>
              <Link 
                href={`?categoryId=${(typeof product === 'object' ? product?.category_id : null) || categoryId}`} 
                className="text-blue-500 hover:underline"
              >
                {Object.entries(categoryMap).find(([, id]) => id.toString() === ((typeof product === 'object' ? product?.category_id : null) || categoryId)?.toString())?.[0]}
              </Link>
              {typeof product === 'object' && product?.name && (
                <span className="inline-flex items-center">
                  <span className="mx-2">{'>'}</span>
                  <span>{product.name}</span>
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
