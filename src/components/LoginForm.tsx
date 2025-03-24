'use client';

import { useState } from 'react';
import { useAuth } from '../hooks/useAuth';

interface LoginFormProps {
  from: string;
}

export default function LoginForm({ from }: LoginFormProps) {
  const [isRegistering, setIsRegistering] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { login, register, csrfToken } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);


    try {
      let response;
      if (isRegistering) {
        if (password !== confirmPassword) {
          setError('Passwords do not match');
          return;
        }
        response = await register(email, password, confirmPassword, csrfToken, from);
      } else {
        response = await login(email, password, csrfToken, from);
      }
      
      if (!response.ok) {
        const data = await response.json();
        if (response.status === 429) {
          setError('Too many attempts. Please try again later.');
        } else if (response.status === 401) {
          setError('Invalid email or password');
        } else if (response.status === 409) {
          setError('Email already registered');
        } else {
          setError(data.error || 'Request failed. Please try again.');
        }
        return;
      }
    } catch {
      setError('An unexpected error occurred. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="max-w-sm mx-auto mt-8">
      <div className="mb-4">
        <label htmlFor="email" className="block mb-2">Email</label>
        <input
          type="email"
          id="email"
          value={email}
          onChange={(e) => {
            const value = e.target.value;
            if (value.length <= 100) {
              setEmail(value.replace(/[^a-zA-Z0-9@._%+-]/g, ''));
            }
          }}
          required
          pattern="^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$"
          className="w-full px-3 py-2 border rounded"
          disabled={isLoading}
          maxLength={100}
        />
      </div>
      <div className="mb-4">
        <label htmlFor="password" className="block mb-2">Password</label>
        <input
          type="password"
          id="password"
          value={password}
          onChange={(e) => {
            const value = e.target.value;
            if (value.length <= 50) {
              setPassword(value.replace(/[^A-Za-z\d@$!%*?&]/g, ''));
            }
          }}
          required
          pattern="^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$"
          className="w-full px-3 py-2 border rounded"
          disabled={isLoading}
          maxLength={50}
          title="Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character"
        />
      </div>
      {isRegistering && (
        <div className="mb-4">
          <label htmlFor="confirmPassword" className="block mb-2">Confirm Password</label>
          <input
            type="password"
            id="confirmPassword"
            value={confirmPassword}
            onChange={(e) => {
              const value = e.target.value;
              if (value.length <= 50) {
                setConfirmPassword(value.replace(/[^A-Za-z\d@$!%*?&]/g, ''));
              }
            }}
            required
            pattern="^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$"
            className="w-full px-3 py-2 border rounded"
            disabled={isLoading}
            maxLength={50}
            title="Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character"
          />
        </div>
      )}

      {error && (
        <div className="mb-4 p-2 text-red-500 bg-red-50 rounded">
          {error}
        </div>
      )}
      <div className="flex flex-col gap-4">
        <button
          type="submit"
          className="w-full px-4 py-2 text-white bg-blue-500 rounded hover:bg-blue-600 disabled:bg-blue-300"
          disabled={isLoading}
        >
          {isLoading ? (isRegistering ? 'Registering...' : 'Logging in...') : 
           (isRegistering ? 'Register' : 'Login')}
        </button>

        <button
          type="button"
          className="w-full px-4 py-2 text-blue-500 bg-transparent border border-blue-500 rounded hover:bg-blue-50"
          onClick={() => setIsRegistering(!isRegistering)}
          disabled={isLoading}
        >
          {isRegistering ? 'Already have an account? Login' : 'Need an account? Register'}
        </button>
      </div>
    </form>
  );
}
