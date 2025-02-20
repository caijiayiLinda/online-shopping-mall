'use client';

import { NextPage } from 'next';

interface ErrorProps {
  statusCode?: number;
  title?: string;
}

const CustomError: NextPage<ErrorProps> = ({ statusCode, title }) => {
  return (
    <div>
      {statusCode ? (
        <p>
          {statusCode}: {title}
        </p>
      ) : (
        <p>Unexpected error occurred</p>
      )}
    </div>
  );
};

CustomError.getInitialProps = async (context) => {
  const { err } = context;

  // Check if the error is a hydration error
  if (err?.message.includes('Hydration failed')) {
    // Suppress the hydration error
    return { statusCode: undefined, title: undefined };
  }

  return { title: undefined }; // Suppress all errors
};

export default CustomError;
