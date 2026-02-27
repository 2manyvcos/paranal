import clsx from 'clsx';
import type { CSSProperties } from 'react';

export default function Image({
  as: Component = 'div',
  className,
  image,
}: {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  as?: any;
  className?: string;
  image?: string;
}) {
  return (
    <Component
      className={clsx(className, 'image')}
      style={
        {
          '--image': image ? `url(${JSON.stringify(image)})` : undefined,
        } as CSSProperties
      }
      data-image={image || undefined}
    >
      <img className="content" src={image || undefined} />
    </Component>
  );
}
