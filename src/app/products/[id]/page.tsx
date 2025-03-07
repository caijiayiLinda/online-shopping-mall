'use client';

import { notFound } from 'next/navigation';
import { products } from '@/data/products';
import AddToCartButton from '@/components/AddToCartButton';
import Nav from '@/components/Nav';
import Image from 'next/image';
import Cart from '@/components/Cart';

import { useParams } from 'next/navigation';
import { useState } from 'react';
import EditProductForm from '@/components/EditProductForm';

export default function ProductPage() {
  const params = useParams();
  const [isEditing, setIsEditing] = useState(false);

  const handleUpdateProduct = async (updatedProduct: typeof product) => {
    if (!product) return;
    const formData = new FormData();
    
    // 添加文本字段
    formData.append('name', updatedProduct.name);
    formData.append('price', updatedProduct.price.toString());
    formData.append('description', updatedProduct.description);
    
    // 如果有新图片则添加
    if (updatedProduct.image.startsWith('data:')) {
      const blob = await fetch(updatedProduct.image).then(r => r.blob());
      formData.append('image', blob, 'product-image.jpg');
    }

    try {
      const response = await fetch(`/api/products?id=${updatedProduct.id}`, {
        method: 'PUT',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('更新产品失败');
      }

      const data = await response.json();
      
      // 更新本地产品数据
      product.name = data.name;
      product.price = data.price;
      product.description = data.description;
      product.image = data.image_url;

      setIsEditing(false);
    } catch (error) {
      console.error('更新产品出错:', error);
      alert('更新产品失败，请稍后重试');
    }
  };
  if (!params || !params.id) {
    notFound();
  }

      const product = products.find((p) => p.id === params.id) as typeof products[0];

  if (!product) {
    notFound();
  }

  return (
    <Nav category={product.category} product={product.name}>
      <div className="text-right mb-4">
        <button 
          onClick={() => setIsEditing(!isEditing)}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          {isEditing ? '取消Edit' : 'Edit产品'}
        </button>
      </div>

      {isEditing ? (
        <EditProductForm product={product} onSubmit={handleUpdateProduct} />
      ) : (
      <main className="max-w-6xl mx-auto py-8 px-4">
        <div className="grid md:grid-cols-2 gap-8">
          <div>
            <Image
              src={product.image}
              alt={product.name}
              width={500}
              height={500}
              className="w-full rounded-lg shadow-md"
            />
          </div>
          <div>
            <h1 className="text-3xl font-bold mb-4">{product.name}</h1>
            <p className="text-gray-600 text-xl mb-6">${product.price}</p>
            <p className="text-gray-700 mb-8">{product.description}</p>
            <AddToCartButton product={product} />
          </div>
        </div>
      </main>
      )}
      <Cart />
    </Nav>
  );
}
