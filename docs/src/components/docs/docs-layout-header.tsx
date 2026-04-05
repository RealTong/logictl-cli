import { MarkdownCopyButton, ViewOptionsPopover } from 'fumadocs-ui/layouts/docs/page';
import { LanguageSwitcher } from '@/components/language-switcher';

type DocsLayoutHeaderProps = {
  title: string;
  description?: string;
  pathname: string;
  markdownUrl: string;
  githubUrl: string;
  eyebrow?: string;
};

export function DocsLayoutHeader({
  title,
  description,
  pathname,
  markdownUrl,
  githubUrl,
  eyebrow = 'Reference',
}: DocsLayoutHeaderProps) {
  return (
    <div className="docs-surface-card docs-bg-layer docs-reveal flex flex-col gap-4 p-4 md:flex-row md:items-center md:justify-between">
      <div className="min-w-0 space-y-2">
        <p className="docs-kicker">{eyebrow}</p>
        <div className="flex flex-wrap items-center gap-2">
          <span className="docs-chip">{title}</span>
          {description ? (
            <p className="docs-caption max-w-2xl text-[0.86rem] text-[var(--docs-text-muted)]">
              {description}
            </p>
          ) : null}
        </div>
      </div>
      <div className="docs-page-header-actions">
        <LanguageSwitcher pathname={pathname} />
        <MarkdownCopyButton markdownUrl={markdownUrl} />
        <ViewOptionsPopover markdownUrl={markdownUrl} githubUrl={githubUrl} />
      </div>
    </div>
  );
}
