import { Link } from '@tanstack/react-router';
import { Languages } from 'lucide-react';
import { cn } from '@/lib/cn';
import { getAlternatePath, getLocaleFromPathname } from '@/lib/i18n';

type LanguageSwitcherProps = {
  pathname: string;
  className?: string;
  zhLabel?: string;
  enLabel?: string;
};

export function LanguageSwitcher({
  pathname,
  className,
  zhLabel = '中文',
  enLabel = 'English',
}: LanguageSwitcherProps) {
  const currentLocale = getLocaleFromPathname(pathname);
  const nextLocale = currentLocale === 'zh' ? 'en' : 'zh';
  const href = getAlternatePath(pathname);
  const nextLabel = nextLocale === 'en' ? enLabel : zhLabel;
  const hint = nextLocale === 'en' ? 'Switch to English' : '切换到中文';

  return (
    <Link
      to={href as never}
      className={cn(
        'docs-surface-card docs-hover-lift inline-flex items-center gap-3 px-3 py-2 text-sm font-medium text-[var(--docs-text)] no-underline',
        className,
      )}
      aria-label={hint}
      title={hint}
    >
      <span className="docs-chip border-transparent bg-[var(--docs-accent-soft)] px-2.5 text-[0.68rem] tracking-[0.18em] text-[var(--docs-text-soft)] uppercase">
        Lang
      </span>
      <span className="inline-flex items-center gap-2">
        <Languages className="size-4 text-[var(--docs-text-soft)]" />
        <span>{nextLabel}</span>
      </span>
    </Link>
  );
}
