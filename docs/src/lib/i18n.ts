export const locales = ['zh', 'en'] as const;

export type Locale = (typeof locales)[number];

export const localeRoots = {
  zh: '/',
  en: '/en',
} as const;

export const docsRoots = {
  zh: '/docs',
  en: '/en/docs',
} as const;

function splitPathSuffix(pathname: string) {
  const hashIndex = pathname.indexOf('#');
  const queryIndex = pathname.indexOf('?');
  const suffixIndex =
    hashIndex === -1 ? queryIndex : queryIndex === -1 ? hashIndex : Math.min(hashIndex, queryIndex);

  if (suffixIndex === -1) {
    return {
      path: pathname,
      suffix: '',
    };
  }

  return {
    path: pathname.slice(0, suffixIndex),
    suffix: pathname.slice(suffixIndex),
  };
}

export function getLocaleFromPathname(pathname: string): Locale {
  const { path } = splitPathSuffix(pathname);

  if (path === '/en' || path.startsWith('/en/')) {
    return 'en';
  }

  return 'zh';
}

export function getAlternatePath(pathname: string): string {
  const { path, suffix } = splitPathSuffix(pathname);

  let alternatePath = path;

  if (path === '/') {
    alternatePath = localeRoots.en;
  } else if (path === localeRoots.en) {
    alternatePath = localeRoots.zh;
  } else if (path === docsRoots.zh || path.startsWith(`${docsRoots.zh}/`)) {
    alternatePath = `${localeRoots.en}${path}`;
  } else if (path === docsRoots.en || path.startsWith(`${docsRoots.en}/`)) {
    alternatePath = path.slice(localeRoots.en.length) || localeRoots.zh;
  } else if (path.startsWith('/en/')) {
    alternatePath = path.slice(localeRoots.en.length) || localeRoots.zh;
  } else if (path.startsWith('/')) {
    alternatePath = `${localeRoots.en}${path}`;
  }

  return `${alternatePath}${suffix}`;
}
