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
  const { res, err } = context;

  // Check if the error is a hydration error
  if (err?.message.includes('Hydration failed')) {
    // Suppress the hydration error
    return { statusCode: undefined, title: undefined };
  }

  // Get the status code from the response
  const statusCode =
    res?.statusCode || err?.statusCode || 400;

  return { statusCode };
};

export default CustomError;
