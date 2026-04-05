import type { ReactNode } from 'react';
import { cn } from '@/lib/cn';

type DocsPageHeaderProps = {
  title: string;
  description?: ReactNode;
  eyebrow?: string;
  labels?: string[];
  actions?: ReactNode;
  children?: ReactNode;
  className?: string;
};

export function DocsPageHeader({
  title,
  description,
  eyebrow,
  labels,
  actions,
  children,
  className,
}: DocsPageHeaderProps) {
  return (
    <header className={cn('docs-surface-card docs-bg-layer docs-page-header docs-reveal', className)}>
      <div className="space-y-4">
        {eyebrow ? <p className="docs-kicker">{eyebrow}</p> : null}
        <div className="space-y-3">
          <h1 className="docs-title-xl">{title}</h1>
          {description ? <div className="docs-lead max-w-3xl">{description}</div> : null}
        </div>
      </div>
      {labels?.length ? (
        <div className="docs-page-header-meta">
          {labels.map((label) => (
            <span key={label} className="docs-chip">
              {label}
            </span>
          ))}
        </div>
      ) : null}
      {actions ? <div className="docs-page-header-actions">{actions}</div> : null}
      {children ? <div className="docs-callout-body">{children}</div> : null}
    </header>
  );
}
