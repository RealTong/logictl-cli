import defaultMdxComponents from 'fumadocs-ui/mdx';
import type { MDXComponents } from 'mdx/types';
import { AgentCallout } from '@/components/docs/agent-callout';
import { CommandBlock } from '@/components/docs/command-block';
import { ConfigBlock } from '@/components/docs/config-block';
import { DocsPageHeader } from '@/components/docs/docs-page-header';
import { LanguageSwitcher } from '@/components/language-switcher';

export function getMDXComponents(components?: MDXComponents) {
  return {
    ...defaultMdxComponents,
    AgentCallout,
    CommandBlock,
    ConfigBlock,
    DocsPageHeader,
    LanguageSwitcher,
    ...components,
  } satisfies MDXComponents;
}

export const useMDXComponents = getMDXComponents;

declare global {
  type MDXProvidedComponents = ReturnType<typeof getMDXComponents>;
}
