import clsx from 'clsx';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

export default function Markdown({
  as: Component = 'div',
  className,
  markdown,
}: {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  as?: any;
  className?: string;
  markdown?: string;
}) {
  return (
    <Component
      className={clsx(className, 'markdown')}
      data-markdown={markdown || undefined}
    >
      <div className="content">
        <ReactMarkdown remarkPlugins={[remarkGfm]} skipHtml>
          {markdown || ''}
        </ReactMarkdown>
      </div>
    </Component>
  );
}
