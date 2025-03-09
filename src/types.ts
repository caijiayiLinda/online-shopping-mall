export interface Product {
  id: string;
  name: string;
  price: number;
  description: string;
  image_url: string;
  thumbnail_url: string;
  category_id: string;
  createdAt?: string;
  updatedAt?: string;
}
