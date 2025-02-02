'use client';

import { Product } from '@/data/products';
import { useCartContext } from '@/context/CartContext';

interface AddToCartButtonProps {
  product: Product;
}

export default function AddToCartButton({ product }: AddToCartButtonProps) {
  const { addToCart } = useCartContext();

  return (
    <button
      onClick={() => addToCart(product)}
      className="w-full bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition-colors"
    >
      加入购物车
    </button>
  );
}
