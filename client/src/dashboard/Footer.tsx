import Markdown from '@/Markdown';
import { useApplication } from '@/application';

export default function Footer() {
  const { footer } = useApplication();

  if (!footer) return null;

  return (
    <footer className="footer markdown">
      <Markdown>{footer}</Markdown>
    </footer>
  );
}
