import { Link, resolvePath, useMatch, type LinkProps } from 'react-router';

export default function NavLink({ to, ...rest }: LinkProps) {
  const match = useMatch(resolvePath(to).pathname);

  return (
    <Link
      to={to}
      {...rest}
      data-active={match != null ? '' : undefined}
      aria-current={match != null ? 'page' : undefined}
    />
  );
}
