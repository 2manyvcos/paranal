import clsx from 'clsx';

export default function Text({
  as: Component = 'div',
  className,
  text,
}: {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  as?: any;
  className?: string;
  text?: string;
}) {
  return (
    <Component
      className={clsx(className, 'text')}
      data-text={text || undefined}
    >
      <span className="content">{text || undefined}</span>
    </Component>
  );
}
