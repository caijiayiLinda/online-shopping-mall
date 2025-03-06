'use client'
import { useState, useEffect } from 'react'
import axios from 'axios'
import { toast } from 'react-hot-toast'

interface Category {
  id: number
  name: string
  createdAt: string
  updatedAt: string
}

interface Product {
  id: number
  name: string
  price: number
  description: string
  image_url: string
  category_id: number
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
        const formattedCategories = response.data.map((category: any) => ({
          id: category.id,
          name: category.name,
          createdAt: category.created_at,
          updatedAt: category.updated_at
        }))
        console.log('Formatted categories:', formattedCategories)
        setCategories(formattedCategories)
      } catch (error) {
        console.error('Failed to fetch categories:', error)
        toast.error('获取类别失败')
      }
    }
    fetchCategories()
    fetchProducts()
  }, [])

  const fetchProducts = async () => {
    try {
      const response = await axios.get('/api/products')
      const formattedProducts = response.data.map((product: any) => ({
        ...product,
        image_url: product.image_url,
        category_id: product.category_id
      }))
      setProducts(formattedProducts)
    } catch (error) {
      console.error('Failed to fetch products:', error)
    }
  }

  const handleDeleteProduct = async (productId: number) => {
    if (window.confirm('确定要删除该产品吗？')) {
      try {
        await axios.delete(`/api/products/delete?id=${productId}`)
        toast.success('产品删除成功')
        fetchProducts()
      } catch (error) {
        console.error('Failed to delete product:', error)
        toast.error('产品删除失败')
      }
    }
  }

  const handleEditProduct = async (product: Product) => {
    const newName = prompt('请输入产品名称', product.name)
    const newPrice = prompt('请输入产品价格', product.price.toString())
    const newDescription = prompt('请输入产品描述', product.description)

    if (newName && newPrice && newDescription) {
      try {
        const formData = new FormData()
        formData.append('name', newName)
        formData.append('price', newPrice)
        formData.append('description', newDescription)
        formData.append('category_id', product.category_id.toString())
        
        await axios.put(`/api/products/update?id=${product.id}`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        toast.success('产品更新成功')
        fetchProducts()
      } catch (error) {
        console.error('Failed to update product:', error)
        toast.error('产品更新失败')
      }
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
      fetchProducts()
    } catch (error) {
      console.error('Failed to create product:', error)
      toast.error('产品创建失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="block text-sm font-medium mb-1">类别</label>
        <select
          value={selectedCategory || ''}
          onChange={(e) => setSelectedCategory(Number(e.target.value))}
          className="w-full p-2 border rounded"
          required
        >
          <option value="">选择类别</option>
          {categories.map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="block text-sm font-medium mb-1">产品名称</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="w-full p-2 border rounded"
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium mb-1">价格</label>
        <input
          type="number"
          value={price}
          onChange={(e) => setPrice(e.target.value)}
          className="w-full p-2 border rounded"
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium mb-1">描述</label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          className="w-full p-2 border rounded"
          rows={4}
          required
        />
      </div>

      <div>
        <label className="block text-sm font-medium mb-1">产品图片</label>
        <input
          type="file"
          accept=".jpg,.jpeg,.png,.gif"
          onChange={(e) => setImage(e.target.files?.[0] || null)}
          className="w-full p-2 border rounded"
          required
        />
        <p className="text-sm text-gray-500 mt-1">支持格式：jpg, png, gif，最大10MB</p>
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-blue-500 text-white p-2 rounded hover:bg-blue-600 disabled:bg-gray-400"
      >
        {loading ? '提交中...' : '提交'}
      </button>

      <button
        type="button"
        onClick={fetchProducts}
        className="w-full bg-gray-200 p-2 rounded hover:bg-gray-300 mt-4"
      >
        刷新产品列表
      </button>

      <div className="mt-8">
        <h2 className="text-xl font-bold mb-4">产品列表</h2>
        <div className="space-y-4">
          {products.map((product) => (
            <div key={product.id} className="border p-4 rounded">
              <div className="flex items-center space-x-4">
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
                    <button
                      type="button"
                      onClick={() => handleEditProduct(product)}
                      className="text-sm bg-yellow-500 text-white px-2 py-1 rounded hover:bg-yellow-600"
                    >
                      编辑
                    </button>
                    <button
                      type="button"
                      onClick={() => handleDeleteProduct(product.id)}
                      className="text-sm bg-red-500 text-white px-2 py-1 rounded hover:bg-red-600"
                    >
                      删除
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </form>
  )
}
