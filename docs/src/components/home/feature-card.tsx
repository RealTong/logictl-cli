import type { ReactNode } from 'react';

type FeatureCardProps = {
  icon: ReactNode;
  title: string;
  description: string;
  detail: string;
  delayMs?: number;
};

export function FeatureCard({ icon, title, description, detail, delayMs = 0 }: FeatureCardProps) {
  return (
    <article
      className="docs-surface-card docs-bg-layer docs-hover-lift docs-reveal flex h-full flex-col gap-4 p-5 sm:p-6"
      style={{ animationDelay: `${delayMs}ms` }}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex h-11 w-11 items-center justify-center rounded-2xl border border-[var(--docs-border)] bg-white/70 text-[var(--docs-text)] shadow-[var(--docs-shadow-inset)] dark:bg-white/5">
          {icon}
        </div>
        <span className="h-2.5 w-2.5 rounded-full bg-[var(--docs-border-strong)]" />
      </div>
      <div className="space-y-3">
        <h3 className="text-lg font-semibold tracking-[-0.03em] text-[var(--docs-text)]">{title}</h3>
        <p className="text-sm leading-7 text-[var(--docs-text-muted)]">{description}</p>
      </div>
      <p className="mt-auto text-sm leading-6 text-[var(--docs-text-soft)]">{detail}</p>
    </article>
  );
}
