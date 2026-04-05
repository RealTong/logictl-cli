import { LanguageSwitcher } from '@/components/language-switcher';
import type { HomeHero } from '@/lib/home-content';

type HeroShellProps = {
  pathname: string;
  hero: HomeHero;
  cta: {
    primary: {
      label: string;
      href: string;
    };
    secondary: {
      label: string;
      href: string;
    };
    tertiary: {
      label: string;
      href: string;
    };
  };
};

export function HeroShell({ pathname, hero, cta }: HeroShellProps) {
  return (
    <section className="docs-surface-card docs-bg-layer docs-reveal overflow-hidden px-5 py-6 md:px-7 md:py-7">
      <div className="grid gap-7 lg:grid-cols-[minmax(0,1.05fr)_minmax(22rem,0.95fr)] lg:items-center">
        <div className="space-y-5">
          <div className="flex flex-wrap items-center gap-3">
            <span className="docs-chip border-transparent bg-[var(--docs-accent-soft)] text-[var(--docs-text-soft)]">
              {hero.eyebrow}
            </span>
            <LanguageSwitcher pathname={pathname} />
          </div>
          <div className="space-y-4">
            <h1 className="docs-title-xl max-w-4xl">{hero.title}</h1>
            <p className="docs-lead max-w-3xl">{hero.description}</p>
          </div>
          <div className="flex flex-wrap gap-3">
            {hero.chips.map((chip) => (
              <span key={chip} className="docs-chip">
                {chip}
              </span>
            ))}
          </div>
          <p className="docs-caption max-w-2xl text-[0.95rem]">{hero.supportingNote}</p>
          <div className="flex flex-wrap gap-3">
            <a
              href={cta.primary.href}
              className="docs-hover-lift inline-flex items-center rounded-full border border-[var(--docs-border-strong)] bg-[var(--docs-text)] px-4 py-2.5 text-sm font-semibold text-white no-underline shadow-[var(--docs-shadow-sm)]"
            >
              {cta.primary.label}
            </a>
            <a
              href={cta.secondary.href}
              className="docs-hover-lift inline-flex items-center rounded-full border border-[var(--docs-border)] bg-white/72 px-4 py-2.5 text-sm font-semibold text-[var(--docs-text)] no-underline shadow-[var(--docs-shadow-inset)] dark:bg-white/5"
            >
              {cta.secondary.label}
            </a>
            <a
              href={cta.tertiary.href}
              className="inline-flex items-center px-1 py-2 text-sm font-medium text-[var(--docs-text-soft)] no-underline"
            >
              {cta.tertiary.label}
            </a>
          </div>
        </div>

        <div className="grid gap-4 lg:pl-4">
          <div className="docs-surface-card-muted docs-bg-layer p-4">
            <div className="flex items-center justify-between gap-3">
              <span className="docs-kicker">{hero.gestureLabel}</span>
              <span className="docs-chip">Local</span>
            </div>
            <div className="mt-4 rounded-2xl border border-[var(--docs-border)] bg-white/72 px-4 py-3 text-sm font-medium tracking-[-0.02em] text-[var(--docs-text)] shadow-[var(--docs-shadow-inset)] dark:bg-white/5">
              {hero.gestureValue}
            </div>
            <p className="mt-3 text-sm leading-6 text-[var(--docs-text-soft)]">{hero.gestureHint}</p>
          </div>

          <div className="docs-surface-card-muted docs-bg-layer p-4">
            <div className="flex items-center justify-between gap-3">
              <span className="docs-kicker">{hero.configTitle}</span>
              <span className="docs-chip">TOML</span>
            </div>
            <p className="mt-3 text-sm leading-6 text-[var(--docs-text-soft)]">{hero.configCaption}</p>
            <pre className="docs-shell-pre px-0 pb-0 pt-4 text-[0.88rem]">
              <code>
                {hero.configLines.map((line) => (
                  <span key={line} className="block">
                    {line}
                  </span>
                ))}
              </code>
            </pre>
          </div>

          <div className="docs-surface-card-muted docs-bg-layer p-4">
            <span className="docs-kicker">{hero.resultLabel}</span>
            <div className="mt-3 space-y-2">
              <h2 className="text-[1.02rem] font-semibold tracking-[-0.03em] text-[var(--docs-text)]">
                {hero.resultTitle}
              </h2>
              <p className="text-sm leading-6 text-[var(--docs-text-muted)]">
                {hero.resultDescription}
              </p>
            </div>
            <div className="mt-4 flex flex-wrap gap-2">
              {hero.resultMeta.map((item) => (
                <span key={item} className="docs-chip">
                  {item}
                </span>
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
