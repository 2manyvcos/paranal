import clsx from 'clsx';
import type { ComponentProps, JSXElementConstructor } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

export default function Markdown<
  As extends
    | keyof React.JSX.IntrinsicElements
    | JSXElementConstructor<unknown> = 'div',
>({
  as: Component = 'div' as As,
  className,
  markdown,
  ...rest
}: {
  as?: As;
  className?: string;
  markdown?: string;
} & Omit<ComponentProps<As>, 'children'>) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const C = Component as any;

  return (
    <C
      {...rest}
      className={clsx(className, 'markdown')}
      data-markdown={markdown || undefined}
    >
      <div className="content">
        <ReactMarkdown remarkPlugins={[remarkGfm]} skipHtml>
          {markdown || ''}
        </ReactMarkdown>
      </div>
    </C>
  );
}
