import { BrainCircuit, Feather, MousePointer2 } from 'lucide-react';
import { createFileRoute } from '@tanstack/react-router';
import { FeatureCard } from '@/components/home/feature-card';
import { HeroShell } from '@/components/home/hero-shell';
import { WorkflowStrip } from '@/components/home/workflow-strip';
import { baseOptions } from '@/lib/layout.shared';
import { homeContent } from '@/lib/home-content';
import { HomeLayout } from 'fumadocs-ui/layouts/home';

export const Route = createFileRoute('/')({
  component: Home,
});

function Home() {
  const content = homeContent.zh;

  return (
    <HomeLayout {...baseOptions()}>
      <main className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-6 px-4 pb-16 pt-6 md:px-6 md:pb-20">
        <HeroShell pathname="/" hero={content.hero} cta={content.cta} />

        <section className="grid gap-4 lg:grid-cols-[minmax(0,0.84fr)_minmax(0,1.16fr)]">
          <div className="docs-surface-card docs-bg-layer docs-reveal flex flex-col gap-4 p-5 md:p-6">
            <p className="docs-kicker">{content.featuresSection.eyebrow}</p>
            <h2 className="docs-title-lg">{content.featuresSection.title}</h2>
            <p className="docs-lead">{content.featuresSection.description}</p>
          </div>
          <div className="grid gap-4 md:grid-cols-3">
            <FeatureCard icon={<Feather className="size-5" />} delayMs={0} {...content.features[0]} />
            <FeatureCard
              icon={<BrainCircuit className="size-5" />}
              delayMs={80}
              {...content.features[1]}
            />
            <FeatureCard
              icon={<MousePointer2 className="size-5" />}
              delayMs={160}
              {...content.features[2]}
            />
          </div>
        </section>

        <WorkflowStrip
          eyebrow={content.workflow.eyebrow}
          title={content.workflow.title}
          description={content.workflow.description}
          steps={content.workflow.steps}
          footnote={content.workflow.footnote}
        />

        <section className="grid gap-4 lg:grid-cols-[minmax(0,0.86fr)_minmax(0,1.14fr)]">
          <div className="docs-surface-card docs-bg-layer docs-reveal flex flex-col gap-4 p-5 md:p-6">
            <p className="docs-kicker">{content.useCasesSection.eyebrow}</p>
            <h2 className="docs-title-lg">{content.useCasesSection.title}</h2>
            <p className="docs-lead">{content.useCasesSection.description}</p>
          </div>
          <div className="grid gap-4 md:grid-cols-2">
            {content.useCases.map((item) => (
              <article
                key={item.title}
                className="docs-surface-card-muted docs-reveal flex flex-col gap-3 p-4"
              >
                <div className="flex flex-wrap items-center gap-3">
                  <h3 className="text-[1rem] font-semibold tracking-[-0.03em] text-[var(--docs-text)]">
                    {item.title}
                  </h3>
                  <span className="docs-chip font-mono text-[0.72rem]">{item.binding}</span>
                </div>
                <p className="docs-caption text-[0.95rem] text-[var(--docs-text-muted)]">
                  {item.description}
                </p>
                <p className="mt-auto text-sm leading-6 text-[var(--docs-text-soft)]">{item.result}</p>
              </article>
            ))}
          </div>
        </section>

        <section className="docs-surface-card docs-bg-layer docs-reveal flex flex-col gap-4 px-5 py-6 md:flex-row md:items-center md:justify-between md:px-6">
          <div className="max-w-2xl space-y-3">
            <p className="docs-kicker">开始使用</p>
            <h2 className="docs-title-lg">{content.footer.title}</h2>
            <p className="docs-lead">{content.footer.description}</p>
          </div>
          <div className="flex flex-wrap gap-3">
            <a
              href={content.cta.primary.href}
              className="docs-hover-lift inline-flex items-center rounded-full border border-[var(--docs-border-strong)] bg-[var(--docs-text)] px-4 py-2.5 text-sm font-semibold text-white no-underline shadow-[var(--docs-shadow-sm)]"
            >
              {content.cta.primary.label}
            </a>
            <a
              href={content.cta.secondary.href}
              className="docs-hover-lift inline-flex items-center rounded-full border border-[var(--docs-border)] bg-white/72 px-4 py-2.5 text-sm font-semibold text-[var(--docs-text)] no-underline shadow-[var(--docs-shadow-inset)] dark:bg-white/5"
            >
              {content.cta.secondary.label}
            </a>
          </div>
        </section>
      </main>
    </HomeLayout>
  );
}
