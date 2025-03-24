'use client';

import { useState } from 'react';
import { useAuth } from '../hooks/useAuth';

export default function ChangePasswordForm() {
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { csrfToken, protectedFetch, logout } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    if (newPassword !== confirmPassword) {
      setError('New passwords do not match');
      return;
    }

    setIsLoading(true);

    try {
      const response = await protectedFetch('/api/auth/change-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-Token': csrfToken,
        },
        body: JSON.stringify({
          currentPassword,
          newPassword,
        }),
      });

      if (!response.ok) {
        const data = await response.json();
        setError(data.error || 'Failed to change password');
        return;
      }

      setSuccess('Password changed successfully. Please login again.');
      await logout();
    } catch (err) {
      setError('An unexpected error occurred');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="max-w-sm mx-auto mt-8">
      <div className="mb-4">
        <label htmlFor="currentPassword" className="block mb-2">Current Password</label>
        <input
          type="password"
          id="currentPassword"
          value={currentPassword}
          onChange={(e) => setCurrentPassword(e.target.value)}
          required
          className="w-full px-3 py-2 border rounded"
          disabled={isLoading}
        />
      </div>
      <div className="mb-4">
        <label htmlFor="newPassword" className="block mb-2">New Password</label>
        <input
          type="password"
          id="newPassword"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          required
          minLength={8}
          maxLength={50}
          pattern="^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$"
          className="w-full px-3 py-2 border rounded"
          disabled={isLoading}
          title="Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character"
        />
      </div>
      <div className="mb-4">
        <label htmlFor="confirmPassword" className="block mb-2">Confirm New Password</label>
        <input
          type="password"
          id="confirmPassword"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          required
          className="w-full px-3 py-2 border rounded"
          disabled={isLoading}
        />
      </div>
      {error && (
        <div className="mb-4 p-2 text-red-500 bg-red-50 rounded">
          {error}
        </div>
      )}
      {success && (
        <div className="mb-4 p-2 text-green-500 bg-green-50 rounded">
          {success}
        </div>
      )}
      <button
        type="submit"
        className="w-full px-4 py-2 text-white bg-blue-500 rounded hover:bg-blue-600 disabled:bg-blue-300"
        disabled={isLoading}
      >
        {isLoading ? 'Changing...' : 'Change Password'}
      </button>
    </form>
  );
}
