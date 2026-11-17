import { useEffect } from 'react';

export type UnexpectedErrorPageProps = {
  error: unknown;
};

export default function UnexpectedErrorPage({
  error,
}: UnexpectedErrorPageProps) {
  useEffect(() => {
    if (error) console.error(error);
  }, [error]);

  // TODO: add UI

  return <div>Something went wrong</div>;
}
