import type { ReactNode } from 'react';
import { FileCode2 } from 'lucide-react';
import { cn } from '@/lib/cn';

type ConfigBlockProps = {
  code?: string;
  children?: string;
  title?: string;
  description?: ReactNode;
  path?: string;
  language?: string;
  className?: string;
};

export function ConfigBlock({
  code,
  children,
  title = 'Configuration',
  description,
  path = '~/.config/logictl/config.toml',
  language = 'TOML',
  className,
}: ConfigBlockProps) {
  const content = children ?? code ?? '';

  return (
    <figure className={cn('docs-surface-card docs-shell docs-config-shell docs-reveal', className)}>
      <div className="docs-shell-header">
        <div className="docs-shell-heading">
          <FileCode2 className="size-4" />
          <span>{title}</span>
        </div>
        <div className="docs-shell-meta">
          <span className="docs-chip">{language}</span>
          <span className="docs-chip">{path}</span>
        </div>
      </div>
      {description ? <figcaption className="docs-shell-caption">{description}</figcaption> : null}
      <pre className="docs-shell-pre">
        <code>{content}</code>
      </pre>
    </figure>
  );
}
