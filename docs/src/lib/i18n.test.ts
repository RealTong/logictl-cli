import { describe, expect, it } from 'vitest';

import { getAlternatePath } from './i18n';

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
