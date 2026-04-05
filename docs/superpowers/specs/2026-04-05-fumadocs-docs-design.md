# logictl Docs Site Design

Date: 2026-04-05
Status: Approved for planning
Primary framework: Fumadocs on TanStack Start
Primary tone: macOS-native, lightweight, AI-first, Agent-first
Language strategy: Full bilingual mirror

## Summary

`logictl` needs a documentation site that does two jobs at once:

- present the product like a polished macOS-native tool
- serve as a fast, trustworthy reference for installation, configuration, gestures, and troubleshooting

The recommended direction is a productized docs site built on top of the existing Fumadocs scaffold. The landing pages should behave like a lightweight product site, while the docs pages should remain documentation-first and preserve the strengths of Fumadocs navigation, search, table of contents, and MDX content flow.

The site should emphasize four traits:

- lightweight runtime and mental model
- macOS-native feel
- AI-first workflows
- Agent-first workflows

The first version should ship as a complete bilingual mirror in Chinese and English, with a shared component system and parallel information architecture.

## Goals

- Make `logictl` look like a serious, polished product rather than a generic docs scaffold
- Give new users a clear path from discovery to first successful setup
- Give existing users a quick way to find configuration, trigger, and troubleshooting details
- Establish a clean bilingual docs structure that can scale without becoming inconsistent
- Highlight AI-first and Agent-first workflows as a core product characteristic rather than a side note

## Non-Goals

- Heavy marketing-style animation systems
- Complex interactive visualizations in V1
- A custom docs engine replacing Fumadocs
- Mixed-language pages with Chinese and English on the same document body
- Exhaustive content for every possible future device or integration in the first pass

## Constraints and Assumptions

- The current scaffold is Fumadocs on TanStack Start and should remain the base
- The design should work well inside the constraints of Fumadocs layouts and MDX-driven content
- Motion should be lightweight and should not depend on a large animation stack
- The docs site should feel visually distinct without fighting the framework
- Chinese and English content should be maintained as mirrored trees, not as inline translation blocks

## Recommended Approach

Three approaches were considered:

1. Content-first customization
2. Productized docs site
3. Marketing-heavy landing page with standard docs behind it

The recommended approach is the productized docs site.

It provides the best balance of:

- visual identity
- implementation risk
- long-term maintainability
- information density
- compatibility with the current Fumadocs scaffold

This approach gives the landing page enough personality to feel like a product while keeping the documentation area efficient, scannable, and framework-friendly.

## Information Architecture

The site should be split into two mirrored trees:

- Chinese root at `/` and `/docs/...`
- English root at `/en` and `/en/docs/...`

Both language trees should use the same:

- visual system
- navigation logic
- component set
- page taxonomy

The first version should ship with these primary sections:

1. Overview
2. Quick Start
3. Configuration
4. Gestures
5. Troubleshooting
6. Examples

The landing page is not the docs home page. It is the product entry point and should guide users into the docs using a few strong calls to action.

Primary entry actions:

- Get Started
- View Config Reference

Secondary entry action:

- Explore Examples

## Landing Page Design

The landing page should feel like a polished macOS utility site rather than a template homepage.

### Core Narrative

The page should tell a simple story:

1. What `logictl` is
2. Why it matters
3. How it works
4. Where to start

### Hero

The hero should state the core positioning in one concise sentence. The message should combine:

- Logitech customization
- local control
- AI-native or Agent-ready workflows
- macOS-native experience

The visual counterpart to the hero should not be a generic illustration. It should resemble a desktop-product composition, such as:

- a mouse gesture state
- a config snippet
- a resulting action or workflow

This helps the page feel product-specific immediately.

### Feature Bands

The landing page should include three top-level feature cards:

- Lightweight
- AI-first
- Agent-first

Each one should tie back to a concrete capability rather than broad product marketing language.

### How It Works

The next section should explain the runtime in three steps:

- Capture
- Interpret
- Execute

This section should visually and verbally connect HID input, rule matching, and output actions without becoming overly technical.

### Real Use Cases

The page should show realistic scenarios such as:

- Mission Control
- desktop switching
- app-specific tab closing
- invoking AI or Agent workflows from gesture bindings

### Footer CTA

The final section should guide the user toward:

- Quick Start
- Config Reference

## Docs Page Design

The docs pages should keep the Fumadocs structural strengths while adding a more intentional product skin.

### Structural Principles

Keep:

- left navigation
- right table of contents
- built-in search
- MDX-driven content

Avoid replacing these with custom structures unless absolutely necessary.

### Visual Language

The docs pages should adopt a macOS-native feel through:

- bright, calm surfaces
- soft gray-blue tonal layering
- subtle translucency
- thin borders
- restrained shadows
- measured spacing

The site should feel premium and quiet rather than loud or highly saturated.

### Content Components

Several content patterns should be elevated into reusable components:

- Hero shell for landing sections
- Feature card
- Command block
- Config block
- Agent tip callout
- Human tip callout
- Troubleshooting block
- Language switcher

These components should make command examples, config snippets, and AI or Agent workflows feel intentional and recognizable throughout the site.

## Bilingual Strategy

The language strategy should be a full mirror rather than partial translation or mixed-language pages.

### Rules

- Chinese and English each get a complete docs tree
- Slugs should match conceptually across languages wherever possible
- Shared components should be language-agnostic
- Content source should be split cleanly by locale rather than combining both languages in one file tree

### Language Switching

The language switcher should try to preserve topic context:

- if the current page has a mirrored equivalent, switch to that page
- otherwise, fall back to the destination language home page

This should be visible and easy to use from both the landing page and the docs shell.

## Content Plan

The first version should cover the following pages in both languages.

### Landing

- Product landing page

### Docs Entry

- Docs home or overview page

### Core Docs

- Quick Start
- Configuration
- Gestures
- Troubleshooting
- Examples

### Content Expectations

Each page should answer practical user questions.

Suggested page rhythm:

- What it is
- Why it matters
- How to use it
- Common pitfalls

The writing should avoid placeholder or template-style documentation. Pages should be grounded in how `logictl` actually behaves.

## Motion Strategy

Motion should be selective and restrained.

Good candidates:

- hero reveal on load
- feature card hover states
- subtle section reveal
- navigation and language-switch transitions
- panel and surface hover refinement

Avoid in V1:

- large animation libraries
- scroll-jacking
- decorative 3D scenes
- motion that slows down reading or navigation

The target is not a flashy site. The target is a product docs experience that feels refined and alive.

## Design Principles

- Product-first on the landing page
- Documentation-first inside docs pages
- macOS-native, not generic SaaS
- lightweight and calm, not visually loud
- AI-first and Agent-first as a built-in story, not a sidebar topic
- bilingual by design, not translated as an afterthought

## Implementation Boundaries

The first version should focus on:

- a polished landing page
- a coherent docs visual system
- a full bilingual docs structure
- strong first-pass content for the core pages

The first version should not try to solve:

- every advanced docs workflow
- every future integration
- elaborate animation systems
- highly bespoke docs navigation beyond what Fumadocs already does well

## Success Criteria

The design is successful when:

- the homepage feels like a real product site rather than a scaffold
- the docs area feels visually aligned with the product while remaining fast to scan
- Chinese and English users can navigate equivalent structures cleanly
- new users can get from landing page to working setup quickly
- existing users can find configuration, gesture, and troubleshooting answers without friction

## Suggested Next Planning Phase

The implementation plan should be split into four tracks:

1. Routing and bilingual source model
2. Shared visual system and docs shell styling
3. Landing page implementation
4. Core bilingual content authoring
