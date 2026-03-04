import clsx from 'clsx';
import type {
  ComponentProps,
  CSSProperties,
  JSXElementConstructor,
} from 'react';

export default function Image<
  As extends
    | keyof React.JSX.IntrinsicElements
    | JSXElementConstructor<unknown> = 'div',
>({
  as: Component = 'div' as As,
  className,
  image,
  ...rest
}: {
  as?: As;
  className?: string;
  image?: string;
} & Omit<ComponentProps<As>, 'children'>) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const C = Component as any;

  return (
    <C
      {...rest}
      className={clsx(className, 'image')}
      style={
        {
          '--image': image ? `url(${JSON.stringify(image)})` : undefined,
        } as CSSProperties
      }
      data-image={image || undefined}
    >
      <img className="content" src={image || undefined} />
    </C>
  );
}
