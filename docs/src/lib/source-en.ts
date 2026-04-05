import { InferPageType, loader } from 'fumadocs-core/source';
import { lucideIconsPlugin } from 'fumadocs-core/source/lucide-icons';
import { enDocs } from 'collections/server';
import { enDocsContentRoute, enDocsRoute } from './shared';

export const source = loader({
  source: enDocs.toFumadocsSource(),
  baseUrl: enDocsRoute,
  plugins: [lucideIconsPlugin()],
});

export function getPageMarkdownUrl(page: InferPageType<typeof source>) {
  const segments = [...page.slugs, 'content.md'];

  return {
    segments,
    url: `${enDocsContentRoute}/${segments.join('/')}`,
  };
}

export async function getLLMText(page: InferPageType<typeof source>) {
  const processed = await page.data.getText('processed');

  return `# ${page.data.title} (${page.url})

${processed}`;
}
