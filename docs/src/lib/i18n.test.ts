import { describe, expect, it } from 'vitest';

import { getAlternatePath } from './i18n';
import { homeContent } from './home-content';

describe('getAlternatePath', () => {
  it('switches between the root home pages', () => {
    expect(getAlternatePath('/')).toBe('/en');
    expect(getAlternatePath('/en')).toBe('/');
  });

  it('switches between mirrored docs paths', () => {
    expect(getAlternatePath('/docs')).toBe('/en/docs');
    expect(getAlternatePath('/en/docs')).toBe('/docs');
    expect(getAlternatePath('/docs/gestures')).toBe('/en/docs/gestures');
    expect(getAlternatePath('/en/docs/gestures')).toBe('/docs/gestures');
    expect(getAlternatePath('/docs/configuration')).toBe('/en/docs/configuration');
  });

  it('prefixes unknown paths with the alternate locale', () => {
    expect(getAlternatePath('/examples')).toBe('/en/examples');
    expect(getAlternatePath('/en/examples')).toBe('/examples');
  });

  it('preserves query strings and hashes', () => {
    expect(getAlternatePath('/docs/gestures?panel=1#top')).toBe('/en/docs/gestures?panel=1#top');
    expect(getAlternatePath('/en/docs/gestures?panel=1#top')).toBe('/docs/gestures?panel=1#top');
  });
});

describe('homeContent', () => {
  it('exposes mirrored hero and entry sections for both locales', () => {
    expect(homeContent.zh.hero.title.length).toBeGreaterThan(0);
    expect(homeContent.en.hero.title.length).toBeGreaterThan(0);

    expect(homeContent.zh.cta.primary.href).toBe('/docs');
    expect(homeContent.zh.cta.secondary.href).toBe('/docs/configuration');
    expect(homeContent.zh.cta.tertiary.href).toBe('/docs/examples');
    expect(homeContent.en.cta.primary.href).toBe('/en/docs');
    expect(homeContent.en.cta.secondary.href).toBe('/en/docs/configuration');
    expect(homeContent.en.cta.tertiary.href).toBe('/en/docs/examples');

    expect(homeContent.zh.features).toHaveLength(3);
    expect(homeContent.en.features).toHaveLength(3);
    expect(homeContent.zh.workflow.steps).toHaveLength(3);
    expect(homeContent.en.workflow.steps).toHaveLength(3);
    expect(homeContent.zh.useCases.length).toBeGreaterThanOrEqual(3);
    expect(homeContent.en.useCases.length).toBeGreaterThanOrEqual(3);
  });
});
