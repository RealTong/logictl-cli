import type { BaseLayoutProps } from 'fumadocs-ui/layouts/shared';
import { appName, githubRepositoryUrl } from './shared';

export function baseOptions(): BaseLayoutProps {
  return {
    nav: {
      title: appName,
    },
    githubUrl: githubRepositoryUrl,
  };
}
