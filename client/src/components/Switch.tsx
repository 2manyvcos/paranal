import { Switch as HeadlessSwitch } from '@headlessui/react';
import clsx from 'clsx';
import type { ComponentProps, JSXElementConstructor } from 'react';

export default function Switch<
  As extends
    | keyof React.JSX.IntrinsicElements
    | JSXElementConstructor<unknown>
    | typeof HeadlessSwitch = typeof HeadlessSwitch,
>({
  as: Component = HeadlessSwitch as As,
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
    <C {...rest} className={clsx(className, 'switch')}>
      <span className="content" />
    </C>
  );
}
