'use client';

import { useState } from 'react';
import { Product } from '@/data/products';

interface CartItem extends Product {
  quantity: number;
}

export function useCart() {
  const [cartItems, setCartItems] = useState<CartItem[]>([]);

  const addToCart = (product: Product, quantity: number = 1) => {
    setCartItems((prev) => {
      const existingItem = prev.find((item) => item.id === product.id);
      if (existingItem) {
        return prev.map((item) =>
          item.id === product.id
            ? { ...item, quantity: item.quantity + quantity }
            : item
        );
      }
      return [...prev, { ...product, quantity }];
    });
  };

  const updateQuantity = (productId: string, quantity: number) => {
    setCartItems((prev) => {
      return prev.map((item) =>
        item.id === productId ? { ...item, quantity: quantity } : item
      );
    });
  };

  return {
    cartItems,
    addToCart,
    updateQuantity,
  };
}
