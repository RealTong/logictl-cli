import { createMiddleware, createStart } from '@tanstack/react-start';
import { isMarkdownPreferred, rewritePath } from 'fumadocs-core/negotiation';
import { redirect } from '@tanstack/react-router';
import {
  docsContentRoute,
  docsRoute,
  enDocsContentRoute,
  enDocsRoute,
} from '@/lib/shared';

type RewriteFn = (pathname: string) => string | undefined;

function createRewritePair(route: string, contentRoute: string) {
  const { rewrite: rewriteDocs } = rewritePath(
    `${route}{/*path}`,
    `${contentRoute}{/*path}/content.md`,
  );
  const { rewrite: rewriteSuffix } = rewritePath(
    `${route}{/*path}.mdx`,
    `${contentRoute}{/*path}/content.md`,
  );

  return {
    rewriteDocs: ((pathname: string) => rewriteDocs(pathname) || undefined) as RewriteFn,
    rewriteSuffix: ((pathname: string) => rewriteSuffix(pathname) || undefined) as RewriteFn,
  };
}

const docsRewrites = [
  createRewritePair(docsRoute, docsContentRoute),
  createRewritePair(enDocsRoute, enDocsContentRoute),
];

function rewritePathname(
  pathname: string,
  rewrites: Array<(pathname: string) => string | undefined>,
) {
  for (const rewrite of rewrites) {
    const path = rewrite(pathname);
    if (path) {
      return path;
    }
  }

  return undefined;
}

const llmMiddleware = createMiddleware().server(({ next, request }) => {
  const url = new URL(request.url);
  const path = rewritePathname(
    url.pathname,
    docsRewrites.map(({ rewriteSuffix }) => rewriteSuffix),
  );

  if (path) {
    throw redirect(new URL(path, url));
  }

  if (isMarkdownPreferred(request)) {
    const docsPath = rewritePathname(
      url.pathname,
      docsRewrites.map(({ rewriteDocs }) => rewriteDocs),
    );
    if (docsPath) {
      throw redirect(new URL(docsPath, url));
    }
  }

  return next();
});

export const startInstance = createStart(() => {
  return {
    requestMiddleware: [llmMiddleware],
  };
});
