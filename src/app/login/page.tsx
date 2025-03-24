'use client';

import LoginForm from '@/components/LoginForm';
import { useSearchParams } from 'next/navigation';

export default function LoginPage() {
  const searchParams = useSearchParams() ?? new URLSearchParams();
  const from = searchParams.get('from') || '/';

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="w-full max-w-md p-8 bg-white rounded-lg shadow-md">
        <h1 className="text-2xl font-bold mb-6 text-center">
          {from.startsWith('/admin') ? 'Admin Login' : 'User Login'}
        </h1>
        <LoginForm from={from} />
      </div>
    </div>
  );
}
