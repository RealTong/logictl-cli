import type { ReactNode } from 'react';
import { Terminal } from 'lucide-react';
import { cn } from '@/lib/cn';

type CommandBlockProps = {
  command?: string | string[];
  children?: string;
  title?: string;
  description?: ReactNode;
  prompt?: string;
  badge?: string;
  className?: string;
};

function getCommandContent(command?: string | string[], children?: string) {
  if (typeof children === 'string' && children.length > 0) {
    return children;
  }

  if (Array.isArray(command)) {
    return command.join('\n');
  }

  return command ?? '';
}

export function CommandBlock({
  command,
  children,
  title = 'Command',
  description,
  prompt = '$',
  badge = 'Shell',
  className,
}: CommandBlockProps) {
  const content = getCommandContent(command, children);
  const lines = content.split('\n');

  return (
    <figure className={cn('docs-surface-card docs-shell docs-command-shell docs-reveal', className)}>
      <div className="docs-shell-header">
        <div className="docs-shell-heading">
          <Terminal className="size-4" />
          <span>{title}</span>
        </div>
        <div className="docs-shell-meta">
          <span className="docs-chip">{badge}</span>
        </div>
      </div>
      {description ? <figcaption className="docs-shell-caption">{description}</figcaption> : null}
      <pre className="docs-shell-pre">
        <code>
          {lines.map((line, index) => (
            <span key={`${line}-${index}`} className="docs-shell-line">
              <span className="docs-shell-prefix">{prompt}</span>
              <span>{line}</span>
            </span>
          ))}
        </code>
      </pre>
    </figure>
  );
}
