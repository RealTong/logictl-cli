# logictl Docs Site Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Turn the Fumadocs scaffold into a polished bilingual `logictl` docs site with a macOS-native landing page, locale-aware docs routing, and complete first-pass content for the core documentation sections.

**Architecture:** Keep Fumadocs and TanStack Start as the foundation, but split content and source loading by locale so Chinese and English can be maintained as mirrored trees. Build a shared visual system and a small set of reusable product-docs components, then layer a custom landing page and locale-aware docs shells on top of the existing layouts.

**Tech Stack:** TanStack Start, Fumadocs, Fumadocs MDX, React 19, Tailwind CSS v4, TypeScript, `oxlint`, `tsc`, Vite build pipeline.

---

## Scope Check

This is one cohesive docs-site project, but it contains four implementation tracks that should remain clearly separated:

1. bilingual routing and source loading
2. shared visual system and reusable components
3. landing page implementation
4. bilingual content authoring

Do not mix all four tracks in one large edit. Each task below should land independently and leave the docs app in a working state.

Use `@systematic-debugging` if Fumadocs collection generation or TanStack file-route behavior differs from expectations. Use `@verification-before-completion` before claiming the docs site is ready.

## Planned File Structure

- `docs/package.json`: docs app scripts and any new dev dependencies needed for validation
- `docs/source.config.ts`: locale-specific MDX collection definitions
- `docs/src/lib/shared.ts`: app metadata, GitHub metadata, canonical route constants
- `docs/src/lib/layout.shared.tsx`: shared nav options, branding, and locale-aware top-level layout config
- `docs/src/lib/source-zh.ts`: Chinese docs loader and helpers
- `docs/src/lib/source-en.ts`: English docs loader and helpers
- `docs/src/lib/i18n.ts`: locale model, mirrored-route helpers, and language switch behavior
- `docs/src/lib/home-content.ts`: structured landing page copy for both locales
- `docs/src/routes/__root.tsx`: global title/meta, provider wiring, and top-level body styling
- `docs/src/routes/index.tsx`: Chinese landing page
- `docs/src/routes/en/index.tsx`: English landing page
- `docs/src/routes/docs/$.tsx`: Chinese docs route
- `docs/src/routes/en/docs/$.tsx`: English docs route
- `docs/src/components/mdx.tsx`: MDX component registry
- `docs/src/components/language-switcher.tsx`: locale switch control
- `docs/src/components/home/hero-shell.tsx`: landing hero shell
- `docs/src/components/home/feature-card.tsx`: shared landing feature card
- `docs/src/components/home/workflow-strip.tsx`: capture/interpret/execute section
- `docs/src/components/docs/command-block.tsx`: emphasized shell command presentation
- `docs/src/components/docs/config-block.tsx`: emphasized `config.toml` presentation
- `docs/src/components/docs/agent-callout.tsx`: reusable AI/Agent callout component
- `docs/src/components/docs/docs-page-header.tsx`: page toolbar wrapper for markdown copy and view options
- `docs/src/styles/app.css`: design tokens, surface styling, motion, typography, and component utility classes
- `docs/content/docs/index.mdx`: Chinese docs overview page
- `docs/content/docs/quick-start.mdx`: Chinese quick start
- `docs/content/docs/configuration.mdx`: Chinese config reference
- `docs/content/docs/gestures.mdx`: Chinese gestures reference
- `docs/content/docs/troubleshooting.mdx`: Chinese troubleshooting
- `docs/content/docs/examples.mdx`: Chinese examples
- `docs/content/en/docs/index.mdx`: English docs overview page
- `docs/content/en/docs/quick-start.mdx`: English quick start
- `docs/content/en/docs/configuration.mdx`: English config reference
- `docs/content/en/docs/gestures.mdx`: English gestures reference
- `docs/content/en/docs/troubleshooting.mdx`: English troubleshooting
- `docs/content/en/docs/examples.mdx`: English examples
- `docs/superpowers/specs/2026-04-05-fumadocs-docs-design.md`: approved design reference

Keep component boundaries small. Avoid building a giant “marketing page” component or a giant locale switcher that mixes route parsing, presentation, and copy selection in one file.

### Task 1: Rebrand the scaffold and establish locale-aware metadata helpers

**Files:**
- Modify: `docs/src/lib/shared.ts`
- Modify: `docs/src/lib/layout.shared.tsx`
- Modify: `docs/src/routes/__root.tsx`
- Create: `docs/src/lib/i18n.ts`
- Create: `docs/src/lib/home-content.ts`
- Test: `docs/src/lib/i18n.test.ts`

- [ ] **Step 1: Write the failing locale helper test**

```ts
import { describe, expect, it } from 'vitest';
import { getAlternatePath } from './i18n';

describe('getAlternatePath', () => {
  it('maps Chinese docs paths to the mirrored English docs path', () => {
    expect(getAlternatePath('/docs/gestures')).toBe('/en/docs/gestures');
  });

  it('maps English docs paths back to Chinese docs path', () => {
    expect(getAlternatePath('/en/docs/gestures')).toBe('/docs/gestures');
  });
});
```

- [ ] **Step 2: Run the focused test to verify it fails**

Run: `cd docs && bun x vitest run src/lib/i18n.test.ts`
Expected: FAIL with missing `getAlternatePath` or missing test runner setup

- [ ] **Step 3: Add the minimal docs test runner support if needed**

```json
{
  "devDependencies": {
    "vitest": "^3.2.4"
  },
  "scripts": {
    "test": "vitest run"
  }
}
```

If `bun x vitest` works without a script, keep the package diff minimal and only add what is needed.

- [ ] **Step 4: Implement locale constants and mirrored-path helpers**

```ts
export type Locale = 'zh' | 'en';

export function getAlternatePath(pathname: string): string {
  if (pathname === '/en') return '/';
  if (pathname.startsWith('/en/docs')) return pathname.replace(/^\/en/, '');
  if (pathname === '/') return '/en';
  if (pathname.startsWith('/docs')) return `/en${pathname}`;
  return pathname.startsWith('/en') ? pathname.slice(3) || '/' : `/en${pathname}`;
}
```

- [ ] **Step 5: Rebrand shared metadata**

```ts
export const appName = 'logictl';
export const docsRoute = '/docs';
export const enDocsRoute = '/en/docs';

export const gitConfig = {
  user: 'RealTong',
  repo: 'logictl-cli',
  branch: 'main',
};
```

Update `baseOptions()` and the root head metadata so the site title and description describe `logictl`, not the scaffold.

- [ ] **Step 6: Add locale-scoped landing copy primitives**

```ts
export const homeContent = {
  zh: {
    heroTitle: '把 Logitech 鼠标变成 AI-first、Agent-ready 的工作入口',
  },
  en: {
    heroTitle: 'Turn your Logitech mouse into an AI-first, agent-ready workflow entry point',
  },
} as const;
```

- [ ] **Step 7: Run the helper test and type checks**

Run: `cd docs && bun x vitest run src/lib/i18n.test.ts && bun run types:check`
Expected: PASS

- [ ] **Step 8: Commit the metadata and locale foundation**

```bash
git add docs/package.json docs/src/lib/shared.ts docs/src/lib/layout.shared.tsx docs/src/routes/__root.tsx docs/src/lib/i18n.ts docs/src/lib/i18n.test.ts docs/src/lib/home-content.ts
git commit -m "feat: add docs locale helpers and branding"
```

### Task 2: Split content and source loading into mirrored Chinese and English trees

**Files:**
- Modify: `docs/source.config.ts`
- Create: `docs/src/lib/source-zh.ts`
- Create: `docs/src/lib/source-en.ts`
- Modify: `docs/src/start.ts`
- Modify: `docs/src/lib/shared.ts`
- Test: `docs/src/lib/i18n.test.ts`

- [ ] **Step 1: Write the failing mirrored-path coverage test**

```ts
it('maps English docs home back to Chinese docs home', () => {
  expect(getAlternatePath('/en/docs')).toBe('/docs');
});
```

- [ ] **Step 2: Run the helper test to verify the locale behavior still fails or is incomplete**

Run: `cd docs && bun x vitest run src/lib/i18n.test.ts`
Expected: FAIL if mirrored docs-home handling or helper coverage is still missing

- [ ] **Step 3: Define separate Chinese and English MDX collections**

```ts
export const zhDocs = defineDocs({ dir: 'content/docs' });
export const enDocs = defineDocs({ dir: 'content/en/docs' });
```

Keep the generated collection names aligned with the exports from `source.config.ts`.

- [ ] **Step 4: Create two locale-specific source loaders**

```ts
export const zhSource = loader({
  source: zhDocs.toFumadocsSource(),
  baseUrl: '/docs',
  plugins: [lucideIconsPlugin()],
});

export const enSource = loader({
  source: enDocs.toFumadocsSource(),
  baseUrl: '/en/docs',
  plugins: [lucideIconsPlugin()],
});
```

- [ ] **Step 5: Update markdown negotiation and route constants for both docs trees**

Make `docs/src/start.ts` rewrite both:

- `/docs/...`
- `/en/docs/...`

to their corresponding markdown content routes.

- [ ] **Step 6: Re-run collection generation, tests, and types**

Run: `cd docs && bun run types:check && bun run build`
Expected: PASS

- [ ] **Step 7: Commit the locale-aware source layer**

```bash
git add docs/source.config.ts docs/src/lib/source-zh.ts docs/src/lib/source-en.ts docs/src/start.ts docs/src/lib/shared.ts docs/src/lib/i18n.test.ts
git commit -m "feat: split docs content sources by locale"
```

### Task 3: Build the shared visual system and reusable docs components

**Files:**
- Modify: `docs/src/styles/app.css`
- Modify: `docs/src/components/mdx.tsx`
- Create: `docs/src/components/language-switcher.tsx`
- Create: `docs/src/components/docs/command-block.tsx`
- Create: `docs/src/components/docs/config-block.tsx`
- Create: `docs/src/components/docs/agent-callout.tsx`
- Create: `docs/src/components/docs/docs-page-header.tsx`
- Test: `docs/src/lib/i18n.test.ts`

- [ ] **Step 1: Write the failing language-switch helper assertion**

```ts
it('keeps docs topic when switching locale', () => {
  expect(getAlternatePath('/docs/configuration')).toBe('/en/docs/configuration');
});
```

- [ ] **Step 2: Run the helper test to verify the switching contract is pinned**

Run: `cd docs && bun x vitest run src/lib/i18n.test.ts`
Expected: PASS or FAIL only if helper logic regressed; do not continue until it is green

- [ ] **Step 3: Introduce docs design tokens and motion primitives**

Add CSS variables and utility classes for:

- background layers
- surface cards
- border and shadow styles
- subtle reveal/hover motion
- typography scale
- highlighted command and config shells

Keep motion lightweight and CSS-first.

- [ ] **Step 4: Build the language switcher and docs presentation components**

```tsx
export function LanguageSwitcher({ pathname }: { pathname: string }) {
  const href = getAlternatePath(pathname);
  return <Link to={href}>English</Link>;
}
```

Also create:

- `CommandBlock`
- `ConfigBlock`
- `AgentCallout`
- `DocsPageHeader`

so content authors can reuse them from MDX.

- [ ] **Step 5: Register the reusable MDX components**

Expose the docs components through `getMDXComponents()` so MDX pages can use them without per-page wiring.

- [ ] **Step 6: Verify type safety and visual build stability**

Run: `cd docs && bun run types:check && bun run build`
Expected: PASS

- [ ] **Step 7: Commit the shared docs visual system**

```bash
git add docs/src/styles/app.css docs/src/components/mdx.tsx docs/src/components/language-switcher.tsx docs/src/components/docs/command-block.tsx docs/src/components/docs/config-block.tsx docs/src/components/docs/agent-callout.tsx docs/src/components/docs/docs-page-header.tsx
git commit -m "feat: add shared docs visual system and components"
```

### Task 4: Implement the bilingual product landing pages

**Files:**
- Modify: `docs/src/routes/index.tsx`
- Create: `docs/src/routes/en/index.tsx`
- Create: `docs/src/components/home/hero-shell.tsx`
- Create: `docs/src/components/home/feature-card.tsx`
- Create: `docs/src/components/home/workflow-strip.tsx`
- Modify: `docs/src/lib/home-content.ts`
- Test: `docs/src/lib/i18n.test.ts`

- [ ] **Step 1: Write the failing landing-copy shape assertion**

```ts
import { homeContent } from './home-content';

it('exposes hero content for both locales', () => {
  expect(homeContent.zh.heroTitle).toBeTruthy();
  expect(homeContent.en.heroTitle).toBeTruthy();
});
```

- [ ] **Step 2: Run the focused test to verify content structure is pinned**

Run: `cd docs && bun x vitest run src/lib/i18n.test.ts`
Expected: PASS if Task 1 already covered the file, otherwise extend the suite until the content contract is enforced

- [ ] **Step 3: Build the shared landing components**

Create focused home components for:

- hero shell
- feature cards
- workflow strip

Do not build a single oversized landing component.

- [ ] **Step 4: Replace the scaffold homepage with the Chinese landing page**

The page should include:

- hero
- primary CTA to `/docs`
- secondary CTA to `/docs/configuration`
- Lightweight / AI-first / Agent-first cards
- Capture / Interpret / Execute section
- realistic use-case section

- [ ] **Step 5: Add the mirrored English landing page at `/en`**

Reuse the same components and structure with English copy.

- [ ] **Step 6: Verify the landing pages build correctly**

Run: `cd docs && bun run types:check && bun run build`
Expected: PASS

- [ ] **Step 7: Commit the landing pages**

```bash
git add docs/src/routes/index.tsx docs/src/routes/en/index.tsx docs/src/components/home/hero-shell.tsx docs/src/components/home/feature-card.tsx docs/src/components/home/workflow-strip.tsx docs/src/lib/home-content.ts
git commit -m "feat: add bilingual product landing pages"
```

### Task 5: Implement locale-aware docs routes and polish the docs shell

**Files:**
- Modify: `docs/src/routes/docs/$.tsx`
- Create: `docs/src/routes/en/docs/$.tsx`
- Modify: `docs/src/lib/layout.shared.tsx`
- Modify: `docs/src/lib/source-zh.ts`
- Modify: `docs/src/lib/source-en.ts`
- Modify: `docs/src/components/docs/docs-page-header.tsx`
- Create: `docs/src/components/docs/docs-layout-header.tsx`
- Test: `docs/src/lib/i18n.test.ts`

- [ ] **Step 1: Write the failing alternate-path assertion for an English docs page**

```ts
it('maps English troubleshooting to Chinese troubleshooting', () => {
  expect(getAlternatePath('/en/docs/troubleshooting')).toBe('/docs/troubleshooting');
});
```

- [ ] **Step 2: Run the helper test to confirm route mirroring still passes before wiring routes**

Run: `cd docs && bun x vitest run src/lib/i18n.test.ts`
Expected: PASS

- [ ] **Step 3: Extract a shared docs-page renderer**

Factor the common docs page chrome so Chinese and English routes only differ by source loader and GitHub content path prefix.

- [ ] **Step 4: Wire the English docs route**

Create `docs/src/routes/en/docs/$.tsx` mirroring the existing Chinese docs route but using the English source loader and English base URL.

- [ ] **Step 5: Add locale-aware shell details**

Update the docs shell to include:

- language switcher
- refined nav branding
- locale-aware markdown/view-source links

- [ ] **Step 6: Verify route loading and production build**

Run: `cd docs && bun run types:check && bun run build`
Expected: PASS

- [ ] **Step 7: Commit the locale-aware docs routes**

```bash
git add docs/src/routes/docs/$.tsx docs/src/routes/en/docs/$.tsx docs/src/lib/layout.shared.tsx docs/src/lib/source-zh.ts docs/src/lib/source-en.ts docs/src/components/docs/docs-page-header.tsx docs/src/components/docs/docs-layout-header.tsx
git commit -m "feat: add locale-aware docs routes and shell"
```

### Task 6: Author the Chinese core docs set

**Files:**
- Modify: `docs/content/docs/index.mdx`
- Create: `docs/content/docs/quick-start.mdx`
- Create: `docs/content/docs/configuration.mdx`
- Create: `docs/content/docs/gestures.mdx`
- Create: `docs/content/docs/troubleshooting.mdx`
- Create: `docs/content/docs/examples.mdx`
- Test: `docs/src/routes/docs/$.tsx`

- [ ] **Step 1: Replace the scaffold docs overview with a Chinese docs overview page**

Make `index.mdx` describe the docs sections, expected setup flow, and where to go next.

- [ ] **Step 2: Write the Quick Start page**

Cover:

- install dependencies or binary expectations
- macOS permissions
- daemon start
- first verification flow

- [ ] **Step 3: Write the Configuration page**

Cover:

- `devices`
- `actions`
- `profiles`
- binding semantics
- representative `config.toml` examples using `ConfigBlock`

- [ ] **Step 4: Write the Gestures page**

Cover:

- `gesture_button_down`
- release-time directional gestures
- supported directions
- how gesture recognition behaves

- [ ] **Step 5: Write the Troubleshooting page**

Cover:

- Input Monitoring
- Accessibility
- daemon install vs local build
- BLE HID capture
- `test event` interpretation

- [ ] **Step 6: Write the Examples page**

Add realistic sample configurations for:

- Mission Control
- desktop switching
- Chrome tab close
- Agent-triggered workflow examples

- [ ] **Step 7: Verify the Chinese docs render in build output**

Run: `cd docs && bun run build`
Expected: PASS

- [ ] **Step 8: Commit the Chinese docs content**

```bash
git add docs/content/docs/index.mdx docs/content/docs/quick-start.mdx docs/content/docs/configuration.mdx docs/content/docs/gestures.mdx docs/content/docs/troubleshooting.mdx docs/content/docs/examples.mdx
git commit -m "docs: add Chinese core docs content"
```

### Task 7: Author the mirrored English core docs set

**Files:**
- Create: `docs/content/en/docs/index.mdx`
- Create: `docs/content/en/docs/quick-start.mdx`
- Create: `docs/content/en/docs/configuration.mdx`
- Create: `docs/content/en/docs/gestures.mdx`
- Create: `docs/content/en/docs/troubleshooting.mdx`
- Create: `docs/content/en/docs/examples.mdx`
- Test: `docs/src/routes/en/docs/$.tsx`

- [ ] **Step 1: Create the English docs overview page**

Mirror the Chinese structure but rewrite naturally in English rather than translating line by line.

- [ ] **Step 2: Write the English Quick Start page**

Keep the structure aligned with the Chinese version and preserve the same install, permission, and daemon flow.

- [ ] **Step 3: Write the English Configuration page**

Mirror the same concepts, examples, and config semantics as the Chinese page.

- [ ] **Step 4: Write the English Gestures page**

Mirror the same directional gesture behavior and examples.

- [ ] **Step 5: Write the English Troubleshooting page**

Mirror the same support surface:

- permissions
- BLE
- daemon behavior
- event testing

- [ ] **Step 6: Write the English Examples page**

Keep example parity with the Chinese docs so language switching does not drop important guidance.

- [ ] **Step 7: Verify the English docs render in build output**

Run: `cd docs && bun run build`
Expected: PASS

- [ ] **Step 8: Commit the English docs content**

```bash
git add docs/content/en/docs/index.mdx docs/content/en/docs/quick-start.mdx docs/content/en/docs/configuration.mdx docs/content/en/docs/gestures.mdx docs/content/en/docs/troubleshooting.mdx docs/content/en/docs/examples.mdx
git commit -m "docs: add English core docs content"
```

### Task 8: Final polish, QA, and delivery verification

**Files:**
- Modify: `docs/src/styles/app.css`
- Modify: `docs/src/routes/index.tsx`
- Modify: `docs/src/routes/en/index.tsx`
- Modify: `docs/src/routes/docs/$.tsx`
- Modify: `docs/src/routes/en/docs/$.tsx`
- Modify: `docs/src/components/*` as needed for final polish

- [ ] **Step 1: Run the full docs validation suite**

Run: `cd docs && bun run lint && bun run types:check && bun run build`
Expected: PASS

- [ ] **Step 2: Launch the docs app locally and manually verify both locales**

Run: `cd docs && bun run dev`
Expected:

- `/` shows the Chinese product landing page
- `/en` shows the English product landing page
- `/docs/...` and `/en/docs/...` both render
- language switching preserves topic when possible

- [ ] **Step 3: Perform visual QA checks**

Check:

- landing-page hierarchy
- docs readability
- command and config block styling
- mobile nav behavior
- search and TOC layout
- hover and reveal motion restraint

- [ ] **Step 4: Fix only polish issues surfaced by verification**

Do not add new features in this step.

- [ ] **Step 5: Re-run the full validation suite**

Run: `cd docs && bun run lint && bun run types:check && bun run build`
Expected: PASS

- [ ] **Step 6: Commit the final polish**

```bash
git add docs
git commit -m "polish: finalize bilingual docs site"
```
