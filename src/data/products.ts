export const categories = [
  'Clothing',
  'Tools', 
  'Toys',
  'Beauty',
  'Pets'
];

export interface Product {
  id: string;
  name: string;
  price: number;
  description: string;
  image_url: string;
  category_id: string;
  quantity?: number;
}

import axios from 'axios';

export const getProducts = async (): Promise<Product[]> => {
  try {
    const response = await axios.get('/api/products');
    return response.data.map((product: any) => ({
      ...product,
      image_url: product.image_url,
      category_id: product.category_id
    }));
  } catch (error) {
    console.error('Failed to fetch products:', error);
    return [];
  }
};
