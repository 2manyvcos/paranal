import { useDashboardContext } from '@rotorjs/react';
import { useEffect } from 'react';
import { useSearchParams } from 'react-router';

export default function ExposeVars() {
  const { vars } = useDashboardContext();
  const [searchParams, setSearchParams] = useSearchParams();

  useEffect(() => {
    const nextSearchParams = new URLSearchParams();
    Object.entries(vars).forEach(([name, { value, exposed }]) => {
      if (exposed) nextSearchParams.set(name, JSON.stringify(value));
    });

    setSearchParams(nextSearchParams);
  }, [vars, searchParams, setSearchParams]);

  return null;
}
