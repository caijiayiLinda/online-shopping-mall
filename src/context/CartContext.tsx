'use client';

import { createContext, useContext } from 'react';
import { useCart } from '@/hooks/useCart';
import { Product } from '@/data/products';

interface CartItem extends Product {
  quantity: number;
}

const CartContext = createContext<{
  cartItems: CartItem[];
  addToCart: (product: Product, quantity?: number) => void;
  updateQuantity: (productId: string, quantity: number) => void;
} | null>(null);

export function CartProvider({ children }: { children: React.ReactNode }) {
  const { cartItems, addToCart, updateQuantity } = useCart();
  return (
    <CartContext.Provider value={{ cartItems, addToCart, updateQuantity }}>
      {children}
    </CartContext.Provider>
  );
}

export function useCartContext() {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCartContext must be used within a CartProvider');
  }
  return context;
}
