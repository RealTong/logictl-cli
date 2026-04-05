import { defineConfig, defineDocs } from 'fumadocs-mdx/config';

const docsOptions = {
  postprocess: {
    includeProcessedMarkdown: true,
  },
};

export const docs = defineDocs({
  dir: 'content/docs',
  docs: docsOptions,
});

export const enDocs = defineDocs({
  dir: 'content/en/docs',
  docs: docsOptions,
});

export default defineConfig();
