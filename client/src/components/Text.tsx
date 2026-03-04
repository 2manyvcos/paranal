import clsx from 'clsx';
import type { ComponentProps, JSXElementConstructor } from 'react';

export default function Text<
  As extends
    | keyof React.JSX.IntrinsicElements
    | JSXElementConstructor<unknown> = 'div',
>({
  as: Component = 'div' as As,
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
      className={clsx(className, 'text')}
      data-text={text || undefined}
    >
      <span className="content">{text || undefined}</span>
    </C>
  );
}
