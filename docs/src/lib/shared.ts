import { docsRoots } from './i18n';

export const appName = 'logictl';
export const siteDescription =
  'logictl is a macOS-native, lightweight, AI-first, agent-first Logitech customization tool.';
export const docsRoute = docsRoots.zh;
export const enDocsRoute = docsRoots.en;
export const docsImageRoute = '/og/docs';
export const docsContentRoute = '/llms.mdx/docs';
export const enDocsContentRoute = '/llms.mdx/en/docs';

export const gitConfig = {
  user: 'RealTong',
  repo: 'logictl-cli',
  branch: 'main',
};

export const githubRepositoryUrl = `https://github.com/${gitConfig.user}/${gitConfig.repo}`;
