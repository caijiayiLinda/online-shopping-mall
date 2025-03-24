'use client';

import { useState, useEffect } from 'react';

type Category = {
  id: number;
  name: string;
};

export default function CategoryPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [newCategory, setNewCategory] = useState('');
  const [editCategory, setEditCategory] = useState<Category | null>(null);

  // Fetch categories on mount
  useEffect(() => {
    fetchCategories();
  }, []);

  const fetchCategories = async () => {
    try {
      const response = await fetch('/categories');
      if (!response.ok) throw new Error('Failed to fetch categories');
      const data = await response.json();
      setCategories(data);
    } catch (error) {
      console.error('Error fetching categories:', error);
    }
  };

  const handleCreateCategory = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCategory.trim()) return;

    try {
      const response = await fetch('/categories/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ name: newCategory }),
      });

      if (!response.ok) throw new Error('Failed to create category');
      setNewCategory('');
      await fetchCategories();
    } catch (error) {
      console.error('Error creating category:', error);
    }
  };

  const handleUpdateCategory = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editCategory) return;

    try {
      const response = await fetch('/categories/update', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ 
          id: editCategory.id.toString(),
          name: editCategory.name 
        }),
      });

      if (!response.ok) throw new Error('Failed to update category');
      setEditCategory(null);
      await fetchCategories();
    } catch (error) {
      console.error('Error updating category:', error);
    }
  };

  const handleDeleteCategory = async (id: number) => {
    const confirmed = window.confirm('Are you sure you want to delete this category?');
    if (!confirmed) return;
    
    try {
      const response = await fetch('/categories/delete', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ id: id.toString() }),
      });

      if (!response.ok) throw new Error('Failed to delete category');
      await fetchCategories();
    } catch (error) {
      console.error('Error deleting category:', error);
    }
  };

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Category Management</h1>

      {/* Create Category Form */}
      <form onSubmit={handleCreateCategory} className="mb-6">
        <div className="flex gap-2">
          <input
            type="text"
            value={newCategory}
            onChange={(e) => setNewCategory(e.target.value)}
            placeholder="New category name"
            className="px-4 py-2 border rounded"
            required
            minLength={1}
            maxLength={255}
            onInvalid={(e) => {
              (e.target as HTMLInputElement).setCustomValidity('Please enter a category name');
            }}
            onInput={(e) => {
              (e.target as HTMLInputElement).setCustomValidity('');
            }}
          />
          <button 
            type="submit"
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Add Category
          </button>
        </div>
      </form>

      {/* Edit Category Form */}
      {editCategory && (
        <form onSubmit={handleUpdateCategory} className="mb-6">
          <div className="flex gap-2">
            <input
              type="text"
              value={editCategory.name}
              onChange={(e) => setEditCategory({ ...editCategory, name: e.target.value })}
              className="px-4 py-2 border rounded"
              required
              minLength={1}
              maxLength={255}
              onInvalid={(e) => {
                (e.target as HTMLInputElement).setCustomValidity('Please enter a category name');
              }}
              onInput={(e) => {
                (e.target as HTMLInputElement).setCustomValidity('');
              }}
            />
            <button
              type="submit"
              className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
            >
              Update
            </button>
            <button
              type="button"
              onClick={() => setEditCategory(null)}
              className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      {/* Categories List */}
      <div className="space-y-2">
        {categories.map((category) => (
          <div key={category.id} className="flex justify-between items-center p-2 border rounded">
            <span>{category.name}</span>
            <div className="flex gap-2">
              <button
                onClick={() => setEditCategory(category)}
                className="px-2 py-1 bg-yellow-500 text-white rounded hover:bg-yellow-600"
              >
                Edit
              </button>
              <button
                onClick={() => handleDeleteCategory(category.id)}
                className="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600"
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
