import clsx from 'clsx';
import type { ComponentProps } from 'react';
import { Link as RouterLink } from 'react-router';

export default function Link<As extends 'a' | typeof RouterLink>({
  as: Component = RouterLink as As,
  className,
  text,
  ...rest
}: {
  as?: As;
  className?: string;
  text?: string;
} & Omit<ComponentProps<As>, 'children'>) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const C = Component as any;

  return (
    <C
      {...rest}
      className={clsx(className, 'link')}
      data-text={text || undefined}
    >
      <span className="content">{text || undefined}</span>
    </C>
  );
}
