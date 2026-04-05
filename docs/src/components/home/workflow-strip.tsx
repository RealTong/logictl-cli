import type { HomeWorkflowStep } from '@/lib/home-content';

type WorkflowStripProps = {
  eyebrow: string;
  title: string;
  description: string;
  steps: HomeWorkflowStep[];
  footnote: string;
};

export function WorkflowStrip({
  eyebrow,
  title,
  description,
  steps,
  footnote,
}: WorkflowStripProps) {
  return (
    <section className="docs-surface-card docs-bg-layer docs-reveal overflow-hidden p-5 md:p-6">
      <div className="flex flex-col gap-3 md:max-w-3xl">
        <p className="docs-kicker">{eyebrow}</p>
        <h2 className="docs-title-lg">{title}</h2>
        <p className="docs-lead">{description}</p>
      </div>
      <div className="mt-6 grid gap-4 md:grid-cols-3">
        {steps.map((step) => (
          <article key={step.label} className="docs-surface-card-muted flex flex-col gap-3 p-4">
            <span className="docs-chip w-fit border-transparent bg-[var(--docs-accent-soft)] text-[var(--docs-text-soft)]">
              {step.label}
            </span>
            <div className="space-y-2">
              <h3 className="text-[1rem] font-semibold tracking-[-0.03em] text-[var(--docs-text)]">
                {step.title}
              </h3>
              <p className="docs-caption text-[0.95rem] text-[var(--docs-text-muted)]">
                {step.description}
              </p>
            </div>
          </article>
        ))}
      </div>
      <p className="mt-5 text-sm leading-6 text-[var(--docs-text-soft)]">{footnote}</p>
    </section>
  );
}
