import type { ReactNode } from 'react';
import { Bot, Sparkles } from 'lucide-react';
import { cn } from '@/lib/cn';

type AgentCalloutProps = {
  title?: string;
  eyebrow?: string;
  command?: string;
  children: ReactNode;
  className?: string;
};

export function AgentCallout({
  title = 'Agent workflow',
  eyebrow = 'Agent-first',
  command,
  children,
  className,
}: AgentCalloutProps) {
  return (
    <aside
      className={cn(
        'docs-surface-card docs-bg-layer docs-callout docs-reveal border-[var(--docs-border-strong)]',
        className,
      )}
    >
      <div className="docs-callout-header">
        <div className="space-y-3">
          <div className="inline-flex items-center gap-2">
            <span className="docs-chip border-transparent bg-[var(--docs-accent-soft)] text-[var(--docs-text-soft)]">
              <Sparkles className="size-3.5" />
              <span>{eyebrow}</span>
            </span>
          </div>
          <div className="flex items-start gap-3">
            <span className="mt-0.5 inline-flex size-9 items-center justify-center rounded-2xl border border-[var(--docs-border)] bg-white/60 text-[var(--docs-text-soft)] shadow-[var(--docs-shadow-inset)]">
              <Bot className="size-4.5" />
            </span>
            <div className="space-y-1">
              <h3 className="docs-title-lg text-[1.2rem]">{title}</h3>
              {command ? <p className="docs-caption font-mono">{command}</p> : null}
            </div>
          </div>
        </div>
      </div>
      <div className="docs-callout-body">{children}</div>
    </aside>
  );
}
