import { useApplication } from '@/application';
import Markdown from '@/components/Markdown';

export default function Footer() {
  const { footer } = useApplication();

  return <Markdown className="footer" as="footer" markdown={footer} />;
}
