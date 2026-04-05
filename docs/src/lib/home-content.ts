import type { Locale } from './i18n';

type HomeSectionIntro = {
  eyebrow: string;
  title: string;
  description: string;
};

export type HomeFeature = {
  title: string;
  description: string;
  detail: string;
};

export type HomeWorkflowStep = {
  label: string;
  title: string;
  description: string;
};

export type HomeUseCase = {
  title: string;
  binding: string;
  description: string;
  result: string;
};

export type HomeHero = {
  eyebrow: string;
  title: string;
  description: string;
  chips: string[];
  supportingNote: string;
  gestureLabel: string;
  gestureValue: string;
  gestureHint: string;
  configTitle: string;
  configCaption: string;
  configLines: string[];
  resultLabel: string;
  resultTitle: string;
  resultDescription: string;
  resultMeta: string[];
};

type HomeContent = {
  hero: HomeHero;
  cta: {
    primary: {
      label: string;
      href: string;
    };
    secondary: {
      label: string;
      href: string;
    };
    tertiary: {
      label: string;
      href: string;
    };
  };
  featuresSection: HomeSectionIntro;
  features: HomeFeature[];
  workflow: HomeSectionIntro & {
    steps: HomeWorkflowStep[];
    footnote: string;
  };
  useCasesSection: HomeSectionIntro;
  useCases: HomeUseCase[];
  footer: {
    title: string;
    description: string;
  };
};

export const homeContent = {
  zh: {
    hero: {
      eyebrow: 'macOS-native Logitech customization',
      title: '把 Logitech 鼠标变成 AI-first、Agent-first 的本地工作入口',
      description:
        'logictl 在 macOS 上捕获手势、判断应用上下文，再把动作交给 Mission Control、桌面切换、脚本或 AI / Agent 工作流。',
      chips: ['本地控制', '低开销', '按应用规则'],
      supportingNote: '适合把常用绑定写进可读、可版本化的配置，而不是埋在难以复用的 GUI 里。',
      gestureLabel: 'Gesture',
      gestureValue: '拇指键 + 上滑',
      gestureHint: '仅在 Code 中命中这条规则',
      configTitle: 'config.toml',
      configCaption: '用简洁规则描述手势、应用条件和回退动作。',
      configLines: [
        '[[bindings]]',
        'trigger = "thumb+swipe-up"',
        'when.app = "Code"',
        'run = "agent.review_selection"',
        'fallback = "mission_control"',
      ],
      resultLabel: 'Result',
      resultTitle: '先触发 Agent 审查，再平滑回退到系统动作',
      resultDescription: '同一个绑定可以先走智能工作流，再在需要时维持稳定的本地系统操作。',
      resultMeta: ['Code', 'thumb+swipe-up', 'local'],
    },
    cta: {
      primary: {
        label: '进入文档',
        href: '/docs',
      },
      secondary: {
        label: '查看配置参考',
        href: '/docs/configuration',
      },
      tertiary: {
        label: '浏览示例',
        href: '/docs/examples',
      },
    },
    featuresSection: {
      eyebrow: 'Why logictl',
      title: '不是新的外壳，而是更贴近桌面工作的控制层',
      description: '重点不在堆更多开关，而在用更短的链路，把 Logitech 输入绑定到你真正频繁执行的动作。',
    },
    features: [
      {
        title: 'Lightweight',
        description: '本地运行，响应路径短，不需要为了一个手势动作绕远路。',
        detail: '更适合高频输入和需要稳定反馈的桌面操作。',
      },
      {
        title: 'AI-first',
        description: '把 AI 提示、摘要、审查或整理动作当成一等绑定目标。',
        detail: '不用跳出当前上下文，就能把动作交给模型工作流。',
      },
      {
        title: 'Agent-first',
        description: '面向脚本、命令和代理任务的输出路径，适合持续执行的自动化。',
        detail: '一个手势可以连到 shell、Raycast、Claude Code 或你自己的 agent entrypoint。',
      },
    ],
    workflow: {
      eyebrow: 'Capture / Interpret / Execute',
      title: '从 HID 输入到 macOS 动作，保持短链路和清晰规则',
      description: '运行时先拿到鼠标输入和上下文，再按配置判断动作，最后交给系统、脚本或 Agent。',
      steps: [
        {
          label: 'Capture',
          title: '捕获手势、按键和前台应用',
          description: '把鼠标输入、按键组合和应用上下文收进统一事件流。',
        },
        {
          label: 'Interpret',
          title: '按规则匹配最合适的动作',
          description: '依据配置判断当前绑定是否命中，并决定是否使用回退逻辑。',
        },
        {
          label: 'Execute',
          title: '立刻触发系统动作、脚本或 Agent',
          description: '把结果交给 Mission Control、桌面切换、命令行或更高层的自动化流程。',
        },
      ],
      footnote: '同一条配置可以同时覆盖系统级快捷动作和更复杂的 AI / Agent 工作流。',
    },
    useCasesSection: {
      eyebrow: 'Real use cases',
      title: '更像真实工作的绑定示例，而不是演示用样板',
      description: '这些场景对应的是桌面操作、窗口管理和日常开发流程里的高频动作。',
    },
    useCases: [
      {
        title: 'Mission Control 总览',
        binding: 'Gesture: thumb + swipe up',
        description: '在任何桌面上快速拉起 Mission Control，看清当前任务分布。',
        result: '结果是系统级视图切换，但仍由统一配置管理。',
      },
      {
        title: '左右切换 Space',
        binding: 'Gesture: thumb + swipe left / right',
        description: '把桌面切换收敛到鼠标上，减少键盘组合键的记忆负担。',
        result: '在多桌面工作流里保持更连续的导航节奏。',
      },
      {
        title: '按应用关闭标签页',
        binding: 'When app = Arc, Code or Terminal',
        description: '同一个侧键在不同应用里映射到关闭标签页、关闭面板或发送命令。',
        result: '规则更贴近上下文，不需要每个应用都重新记手势。',
      },
      {
        title: '调用 AI / Agent 工作流',
        binding: 'Run: agent.review_selection',
        description: '把一条手势直接接到摘要、代码审查、命令生成或你自己的 agent 脚本。',
        result: '从输入设备直接进入 AI / Agent 执行链，而不是先切到另一个启动器。',
      },
    ],
    footer: {
      title: '先进入文档，再决定要把哪些动作交给鼠标',
      description: '快速开始页帮助你完成首次配置，配置参考页负责把规则写清楚。',
    },
  },
  en: {
    hero: {
      eyebrow: 'macOS-native Logitech customization',
      title: 'Turn your Logitech mouse into a local, AI-first, agent-first entry point',
      description:
        'logictl captures gestures on macOS, interprets app context, and routes the result to Mission Control, desktop switching, scripts, or AI and agent workflows.',
      chips: ['Local control', 'Low overhead', 'App-aware rules'],
      supportingNote: 'It is built for readable, versionable configuration instead of burying serious workflow logic inside a GUI.',
      gestureLabel: 'Gesture',
      gestureValue: 'Thumb button + swipe up',
      gestureHint: 'This rule only applies inside Code',
      configTitle: 'config.toml',
      configCaption: 'Describe gestures, app conditions, and fallback actions with a short rule.',
      configLines: [
        '[[bindings]]',
        'trigger = "thumb+swipe-up"',
        'when.app = "Code"',
        'run = "agent.review_selection"',
        'fallback = "mission_control"',
      ],
      resultLabel: 'Result',
      resultTitle: 'Trigger an agent review first, then fall back to a system action',
      resultDescription: 'One binding can route into an intelligent workflow while still keeping a reliable local macOS escape hatch.',
      resultMeta: ['Code', 'thumb+swipe-up', 'local'],
    },
    cta: {
      primary: {
        label: 'Read the docs',
        href: '/en/docs',
      },
      secondary: {
        label: 'Open config reference',
        href: '/en/docs/configuration',
      },
      tertiary: {
        label: 'Explore examples',
        href: '/en/docs/examples',
      },
    },
    featuresSection: {
      eyebrow: 'Why logictl',
      title: 'A focused control layer for desktop work, not another generic wrapper',
      description: 'The point is not more toggles. The point is mapping Logitech input to the actions you repeatedly use with less overhead and more context.',
    },
    features: [
      {
        title: 'Lightweight',
        description: 'Runs locally with a short response path, so a gesture still feels like a direct desktop action.',
        detail: 'That matters when the action is something you trigger dozens of times per day.',
      },
      {
        title: 'AI-first',
        description: 'Treat prompts, summaries, reviews, and command generation as first-class gesture targets.',
        detail: 'You can move from input to model-assisted action without leaving the current workspace.',
      },
      {
        title: 'Agent-first',
        description: 'Route bindings into scripts, commands, and agent entry points that are meant to keep running.',
        detail: 'The same gesture can call shell scripts, Raycast, Claude Code, or your own automation stack.',
      },
    ],
    workflow: {
      eyebrow: 'Capture / Interpret / Execute',
      title: 'From HID input to macOS action with a short, readable runtime path',
      description: 'The runtime captures the event, interprets it against config, and then executes the right local action or workflow.',
      steps: [
        {
          label: 'Capture',
          title: 'Collect gestures, buttons, and foreground app context',
          description: 'Mouse input, button combinations, and app state are pulled into one event stream.',
        },
        {
          label: 'Interpret',
          title: 'Match the event against the rule that should win',
          description: 'Config decides whether the binding applies now and whether a fallback should be used.',
        },
        {
          label: 'Execute',
          title: 'Launch the macOS action, script, or agent immediately',
          description: 'The result can be a system shortcut, a shell command, or a higher-level automation flow.',
        },
      ],
      footnote: 'The same configuration surface can cover both system-level shortcuts and richer AI or agent workflows.',
    },
    useCasesSection: {
      eyebrow: 'Real use cases',
      title: 'Bindings that look like actual desktop work, not demo-only examples',
      description: 'These scenarios map to window management, workspace navigation, and common development routines.',
    },
    useCases: [
      {
        title: 'Mission Control at a thumb gesture',
        binding: 'Gesture: thumb + swipe up',
        description: 'Open Mission Control from anywhere to get an immediate view of active work across spaces.',
        result: 'The result is still a native system action, managed from one config surface.',
      },
      {
        title: 'Move left and right across spaces',
        binding: 'Gesture: thumb + swipe left / right',
        description: 'Keep desktop navigation on the mouse and reduce the cognitive load of keyboard shortcuts.',
        result: 'Workspace movement feels continuous during multi-space workflows.',
      },
      {
        title: 'Close tabs differently per app',
        binding: 'When app = Arc, Code, or Terminal',
        description: 'The same side button can close a browser tab, dismiss a panel, or send an app-specific command.',
        result: 'Rules stay contextual without asking you to memorize a different gesture for every app.',
      },
      {
        title: 'Invoke AI or agent workflows directly',
        binding: 'Run: agent.review_selection',
        description: 'Bind a gesture to summarization, code review, command generation, or your own agent script.',
        result: 'You move straight from mouse input into an AI or agent execution path instead of switching launchers first.',
      },
    ],
    footer: {
      title: 'Start in the docs, then decide which actions belong on the mouse',
      description: 'Quick start gets the first config running. The configuration reference helps you shape the rules precisely.',
    },
  },
} satisfies Record<Locale, HomeContent>;
