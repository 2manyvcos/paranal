import { Button as HeadlessButton } from '@headlessui/react';
import clsx from 'clsx';
import type { ComponentProps } from 'react';

export default function Button({
  as: Component = HeadlessButton,
  className,
  text,
  ...rest
}: {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  as?: any;
  className?: string;
  text?: string;
} & Omit<ComponentProps<'button'>, 'children'>) {
  return (
    <Component
      {...rest}
      className={clsx(className, 'button')}
      data-text={text || undefined}
    >
      <span className="content">{text || undefined}</span>
    </Component>
  );
}
