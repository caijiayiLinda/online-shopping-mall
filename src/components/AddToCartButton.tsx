'use client';

import { useCartContext } from '@/context/CartContext';
import { Product } from '@/types';

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
      Add to cart
    </button>
  );
}
