'use client';

import { useState, useEffect } from 'react'
import axios from 'axios'
import { toast } from 'react-hot-toast'
import EditProductForm from './EditProductForm'
import { Product } from '../types'

interface Category {
  id: number
  name: string
  created_at: string
  updated_at: string
}

export default function AdminForm() {
  const [categories, setCategories] = useState<Category[]>([])
  const [selectedCategory, setSelectedCategory] = useState<number>(0)
  const [name, setName] = useState('')
  const [price, setPrice] = useState('')
  const [description, setDescription] = useState('')
  const [image, setImage] = useState<File | null>(null)
  const [loading, setLoading] = useState(false)
  const [products, setProducts] = useState<Product[]>([])

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const response = await axios.get('/api/categories')
        const formattedCategories = response.data.map((category: Category) => ({
          id: category.id,
          name: category.name,
          created_at: category.created_at,
          updated_at: category.updated_at
        }))
        console.log('Formatted categories:', formattedCategories)
        setCategories(formattedCategories)
      } catch (error) {
        console.error('Failed to fetch categories:', error)
        toast.error('获取Category失败')
      }
    }
    fetchCategories()
    fetchProducts()
  }, [])

  const fetchProducts = async () => {
    try {
      const response = await axios.get('/api/products')
      const formattedProducts = response.data.map((product: Product & { catid: string }) => ({
        ...product,
        image_url: product.image_url,
        category_id: product.catid
      }))
      setProducts(formattedProducts)
    } catch (error) {
      console.error('Failed to fetch products:', error)
    }
  }

  const handleDeleteProduct = async (productId: string) => {
    if (window.confirm('确定要Delete该产品吗？')) {
      try {
        const url = new URL('/api/products/delete', window.location.origin)
        url.searchParams.set('id', productId)
        await axios.delete(url.toString())
        toast.success('产品Delete成功')
        fetchProducts()
      } catch (error) {
        console.error('Failed to delete product:', error)
        toast.error('产品Delete失败')
      }
    }
  }

  const [editingProduct, setEditingProduct] = useState<Product | null>(null)

  const handleEditProduct = (product: Product) => {
    setEditingProduct(product)
  }

  const handleCancelEdit = () => {
    setEditingProduct(null)
  }

  const handleUpdateProduct = async (updatedProduct: Product) => {
    try {
      const formData = new FormData()
      formData.append('name', updatedProduct.name)
      formData.append('price', updatedProduct.price.toString())
      formData.append('description', updatedProduct.description)
      formData.append('category_id', updatedProduct.category_id.toString())
      if (image) {
        formData.append('image', image, image.name)
      }
  
      
      const url = new URL('/api/products/update', window.location.origin)
      url.searchParams.set('id', updatedProduct.id)
      await axios.put(url.toString(), formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      toast.success('产品更新成功')
      setEditingProduct(null)
      fetchProducts()
    } catch (error) {
      console.error('Failed to update product:', error)
      toast.error('产品更新失败')
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    
    const formData = new FormData()
    formData.append('category_id', selectedCategory?.toString() || '')
    formData.append('name', name)
    formData.append('price', parseFloat(price).toString())
    formData.append('description', description)
    if (image) {
      formData.append('image', image, image.name)
    }

    // Log FormData contents
    console.log('FormData entries:')
    for (const pair of formData.entries()) {
      console.log(pair[0], pair[1])
    }

    try {
      await axios.post('/api/products/create', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Accept': 'application/json'
        },
        transformRequest: (data) => {
          console.log('Transforming request data:', data);
          return data;
        }
      })
      toast.success('产品创建成功')
      // 重置表单字段
      setName('')
      setPrice('')
      setDescription('')
      setImage(null)
      fetchProducts()
    } catch (error) {
      console.error('Failed to create product:', error)
      toast.error('产品创建失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-4">
      {editingProduct && (
        <div className="bg-blue-50 p-6 rounded-lg shadow mb-8 border border-blue-200">
          <h2 className="text-xl font-bold mb-4">Edit Product</h2>
          <EditProductForm 
            product={{
              id: editingProduct.id.toString(),
              name: editingProduct.name,
              price: editingProduct.price,
              description: editingProduct.description,
              image_url: editingProduct.image_url,
              thumbnail_url: editingProduct.thumbnail_url,
              category_id: editingProduct.category_id.toString()
            }}
            onSubmit={handleUpdateProduct}
          />
          <button
            onClick={handleCancelEdit}
            className="w-full bg-gray-200 p-2 rounded hover:bg-gray-300 mt-4"
          >
            Cancel
          </button>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="block text-sm font-medium mb-1">Category</label>
        <select
          value={selectedCategory || ''}
          onChange={(e) => setSelectedCategory(Number(e.target.value))}
          className="w-full p-2 border rounded"
          required
        >
          <option value="">Select Category</option>
          {categories.map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="block text-sm font-medium mb-1">Product Name</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="w-full p-2 border rounded"
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium mb-1">Product Price</label>
        <input
          type="number"
          value={price}
          onChange={(e) => setPrice(e.target.value)}
          className="w-full p-2 border rounded"
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium mb-1">Product Description</label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          className="w-full p-2 border rounded"
          rows={4}
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium mb-1">Product Image</label>
        <input
          type="file"
          accept=".jpg,.jpeg,.png,.gif"
          onChange={(e) => setImage(e.target.files?.[0] || null)}
          className="w-full p-2 border rounded"
          required
        />
        <p className="text-sm text-gray-500 mt-1">Supported formats: jpg, png, gif. The maximum size is 10MB.</p>
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-blue-500 text-white p-2 rounded hover:bg-blue-600 disabled:bg-gray-400"
      >
        {loading ? 'submitting...' : 'Submit'}
      </button>

        <button
          type="button"
          onClick={fetchProducts}
          className="w-full bg-gray-200 p-2 rounded hover:bg-gray-300 mt-4 hidden"
        >
          Refresh
        </button>

      <div className="mt-8">
        <h2 className="text-xl font-bold mb-4">产品列表</h2>
        <div className="space-y-4">
          {products.map((product) => (
            <div key={product.id} className="border p-4 rounded">
              <div className="flex items-center space-x-4">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img 
                  src={product.image_url} 
                  alt={product.name}
                  className="w-20 h-20 object-cover rounded"
                />
                <div>
                  <h3 className="font-medium">{product.name}</h3>
                  <p className="text-gray-600">¥{product.price}</p>
                  <p className="text-sm text-gray-500">category id: {product.category_id}</p>
                  <p className="text-sm text-gray-500">{product.description}</p>
                  <div className="flex space-x-2 mt-2">
                    {!editingProduct && (
                      <button
                        type="button"
                        onClick={() => handleEditProduct(product)}
                        className="text-sm bg-yellow-500 text-white px-2 py-1 rounded hover:bg-yellow-600"
                      >
                        Edit
                      </button>
                    )}
                    <button
                      type="button"
                      onClick={() => handleDeleteProduct(product.id.toString())}
                      className="text-sm bg-red-500 text-white px-2 py-1 rounded hover:bg-red-600"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </form>
    </div>
  )
}
