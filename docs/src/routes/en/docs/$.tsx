import { createFileRoute, notFound } from '@tanstack/react-router';
import { createServerFn } from '@tanstack/react-start';
import { DocsLayout } from 'fumadocs-ui/layouts/docs';
import { DocsBody, DocsPage } from 'fumadocs-ui/layouts/docs/page';
import { useFumadocsLoader } from 'fumadocs-core/source/client';
import browserCollections from 'collections/browser';
import { Suspense } from 'react';
import { DocsLayoutHeader } from '@/components/docs/docs-layout-header';
import { useMDXComponents } from '@/components/mdx';
import { baseOptions } from '@/lib/layout.shared';
import { gitConfig } from '@/lib/shared';
import { getPageMarkdownUrl, source } from '@/lib/source-en';

export const Route = createFileRoute('/en/docs/$')({
  component: Page,
  loader: async ({ params }) => {
    const slugs = params._splat?.split('/') ?? [];
    const data = await serverLoader({ data: slugs });
    await clientLoader.preload(data.path);
    return data;
  },
});

const serverLoader = createServerFn({
  method: 'GET',
})
  .inputValidator((slugs: string[]) => slugs)
  .handler(async ({ data: slugs }) => {
    const page = source.getPage(slugs);
    if (!page) throw notFound();

    return {
      path: page.path,
      url: page.url,
      markdownUrl: getPageMarkdownUrl(page).url,
      pageTree: await source.serializePageTree(source.getPageTree()),
    };
  });

const clientLoader = browserCollections.enDocs.createClientLoader({
  component(
    { toc, frontmatter, default: MDX },
    {
      markdownUrl,
      path,
      url,
    }: {
      markdownUrl: string;
      path: string;
      url: string;
    },
  ) {
    return (
      <DocsPage toc={toc}>
        <DocsLayoutHeader
          title={frontmatter.title}
          description={frontmatter.description}
          pathname={url}
          markdownUrl={markdownUrl}
          githubUrl={`https://github.com/${gitConfig.user}/${gitConfig.repo}/blob/${gitConfig.branch}/docs/content/en/docs/${path}`}
          eyebrow="Reference"
        />
        <DocsBody className="mt-8">
          <MDX components={useMDXComponents()} />
        </DocsBody>
      </DocsPage>
    );
  },
});

function Page() {
  const { path, pageTree, markdownUrl, url } = useFumadocsLoader(Route.useLoaderData());

  return (
    <DocsLayout {...baseOptions()} tree={pageTree}>
      <Suspense>{clientLoader.useContent(path, { markdownUrl, path, url })}</Suspense>
    </DocsLayout>
  );
}
