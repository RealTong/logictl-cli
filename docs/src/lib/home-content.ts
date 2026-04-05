import type { Locale } from './i18n';
import { docsRoute, enDocsRoute } from './shared';

type HomeFeature = {
  title: string;
  description: string;
};

type HomeStep = {
  title: string;
  description: string;
};

type HomeContent = {
  hero: {
    eyebrow: string;
    title: string;
    description: string;
  };
  features: HomeFeature[];
  steps: HomeStep[];
  cta: {
    primary: {
      label: string;
      href: string;
    };
    secondary: {
      label: string;
      href: string;
    };
  };
};

export const homeContent = {
  zh: {
    hero: {
      eyebrow: 'macOS-native Logitech customization',
      title: '把 Logitech 鼠标变成 AI-first、Agent-ready 的工作入口',
      description: '用轻量、局部控制的方式，把手势、按键和工作流绑定到你真正需要的动作。',
    },
    features: [
      {
        title: '轻量',
        description: '本地运行，保留最小的心智负担和最短的响应链路。',
      },
      {
        title: 'AI-first',
        description: '把 AI 调用和快捷工作流当成一等公民，而不是附加功能。',
      },
      {
        title: 'Agent-first',
        description: '面向自动化代理的输出路径，适合持续执行和组合式任务。',
      },
    ],
    steps: [
      {
        title: 'Capture',
        description: '捕获鼠标动作、按键和上下文。',
      },
      {
        title: 'Interpret',
        description: '按配置规则判断意图并选择动作。',
      },
      {
        title: 'Execute',
        description: '在 macOS 上触发最终工作流或系统操作。',
      },
    ],
    cta: {
      primary: {
        label: '开始阅读',
        href: docsRoute,
      },
      secondary: {
        label: '切换到英文',
        href: enDocsRoute,
      },
    },
  },
  en: {
    hero: {
      eyebrow: 'macOS-native Logitech customization',
      title: 'Turn your Logitech mouse into an AI-first, agent-ready workflow entry point',
      description: 'Bind gestures, buttons, and workflows to the actions you actually use, with local control and low overhead.',
    },
    features: [
      {
        title: 'Lightweight',
        description: 'Runs locally with a minimal mental model and a short response path.',
      },
      {
        title: 'AI-first',
        description: 'Treat AI prompts and fast workflows as first-class actions, not add-ons.',
      },
      {
        title: 'Agent-first',
        description: 'Support output paths that are ready for automated agents and repeatable execution.',
      },
    ],
    steps: [
      {
        title: 'Capture',
        description: 'Collect mouse gestures, button events, and context.',
      },
      {
        title: 'Interpret',
        description: 'Match rules against the current state and route the action.',
      },
      {
        title: 'Execute',
        description: 'Launch the final workflow or macOS action immediately.',
      },
    ],
    cta: {
      primary: {
        label: 'Read the docs',
        href: enDocsRoute,
      },
      secondary: {
        label: 'Switch to Chinese',
        href: docsRoute,
      },
    },
  },
} satisfies Record<Locale, HomeContent>;
