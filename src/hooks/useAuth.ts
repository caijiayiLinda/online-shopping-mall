'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const [csrfToken, setCsrfToken] = useState('');
  const [userEmail, setUserEmail] = useState('');
  const router = useRouter();

  const getCSRFToken = useCallback(async () => {
    try {
      const response = await fetch('/api/auth/csrf-token', {
        credentials: 'include',
      });
      
      if (response.ok) {
        const data = await response.json();
        setCsrfToken(data.csrf_token);
      }
    } catch (error) {
      console.error('Failed to get CSRF token:', error);
    }
  }, []);

  const checkAuth = useCallback(async () => {
    try {
      const response = await fetch('/api/auth/check', {
        credentials: 'include',
      });

      if (!response.ok) {
        setIsAuthenticated(false);
        setIsAdmin(false);
        return;
      }

      const data = await response.json();
      console.log('Auth check response:', {
        authenticated: data.authenticated,
        admin: data.admin,
        email: data.email,
        userId: data.userId
      });
      setIsAuthenticated(data.authenticated === true);
      setIsAdmin(data.admin === true);
      setUserEmail(data.email || '');
    } catch (error) {
      console.error('Auth check failed:', error);
      setIsAuthenticated(false);
      setIsAdmin(false);
    }
  }, []);

  useEffect(() => {
    checkAuth();
    getCSRFToken();
  }, [checkAuth, getCSRFToken]);

  const register = async (email: string, password: string, confirmPassword: string, csrfToken: string, from: string) => {
    console.log('Register attempt:', { email, from });
    try {
      const response = await fetch('/api/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-Token': csrfToken,
        },
        body: JSON.stringify({ email, password, confirmPassword }),
        credentials: 'include',
      });

      if (response.ok) {
        const data = await response.json();
        setCsrfToken(data.csrf_token);
        await checkAuth();
        
        // Handle redirect based on login source and admin status
        if (data.admin) {
          router.push('/admin');
        } else if (from === '/login' || from === '/admin/login') {
          router.push('/');
        } else {
          router.push(from);
        }
      }
      return response;
    } catch (error) {
      console.error('Registration failed:', error);
      throw error;
    }
  };

  const login = async (email: string, password: string, csrfToken: string, from: string) => {
    console.log('Login attempt:', { email, from });
    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-Token': csrfToken,
        },
        body: JSON.stringify({ email, password }),
        credentials: 'include',
      });

      if (response.ok) {
        const data = await response.json();
        setCsrfToken(data.csrf_token);
        await checkAuth();
        
        // Handle redirect based on login source and admin status
        if (data.admin) {
          router.push('/admin');
        } else if (from === '/login' || from === '/admin/login') {
          router.push('/');
        } else {
          router.push(from);
        }
      }
      return response;
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    }
  };

  const logout = async () => {
    console.log('Logout attempt');
    try {
      const response = await fetch('/api/auth/logout', {
        method: 'POST',
        credentials: 'include',
      });

      if (response.ok) {
        setIsAuthenticated(false);
        setIsAdmin(false);
        setCsrfToken('');
        router.push('/login');
      }
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  const protectedFetch = async (url: string, options: RequestInit = {}) => {
    const headers = {
      ...options.headers,
      'X-CSRF-Token': csrfToken,
    };

    try {
      const response = await fetch(url, {
        ...options,
        credentials: 'include',
        headers,
      });

      if (response.status === 401) {
        // Token expired, attempt refresh
        const refreshResponse = await fetch('/api/auth/refresh', {
          method: 'POST',
          credentials: 'include',
        });

        if (refreshResponse.ok) {
          return protectedFetch(url, options); // Retry with new token
        }
      }

      return response;
    } catch (error) {
      console.error('Request failed:', error);
      throw error;
    }
  };

  return {
    isAuthenticated,
    isAdmin,
    userEmail,
    csrfToken,
    register,
    login,
    logout,
    checkAuth,
    protectedFetch,
  };
}
