import clsx from 'clsx';
import {
  Link as RouterLink,
  resolvePath,
  useMatch,
  type LinkProps,
} from 'react-router';

export default function NavItem({
  className,
  to,
  text,
  ...rest
}: { text?: string } & Omit<LinkProps, 'children'>) {
  const match = useMatch(resolvePath(to).pathname);

  return (
    <RouterLink
      className={clsx(className, 'nav-item')}
      to={to}
      {...rest}
      data-active={match != null ? '' : undefined}
      aria-current={match != null ? 'page' : undefined}
      data-text={text || undefined}
    >
      <span className="content">{text || undefined}</span>
    </RouterLink>
  );
}
