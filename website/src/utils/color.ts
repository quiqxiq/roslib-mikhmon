/**
 * Resolve a CSS custom property string (e.g. `var(--accent-cyan)`) to its
 * computed value by reading the document root style. Falls back to the
 * input string if it is not a CSS variable.
 */
export function resolveCssColor(color: string): string {
  if (typeof document === 'undefined') return color
  if (color.trim().startsWith('var(')) {
    const match = color.match(/var\((--[^,)]+)\)/)
    if (match) {
      const raw = getComputedStyle(document.documentElement)
        .getPropertyValue(match[1])
        .trim()
      if (raw) return raw
    }
  }
  return color
}

/**
 * Convert any browser-understandable color string into an `rgba(...)` string
 * with the requested alpha. Supports hex (3/4/6/8), rgb, rgba, CSS variables,
 * and named colors.
 *
 * Uses a temporary DOM element so the browser handles parsing.
 */
export function toRgba(color: string, alpha: number): string {
  if (typeof document === 'undefined') return `rgba(0,0,0,${alpha})`

  const resolved = resolveCssColor(color)

  // Fast path for rgba() so we can just swap the alpha
  const rgbaMatch = resolved.match(
    /rgba\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*([\d.]+)\s*\)/,
  )
  if (rgbaMatch) {
    return `rgba(${rgbaMatch[1]}, ${rgbaMatch[2]}, ${rgbaMatch[3]}, ${alpha})`
  }

  // Fast path for rgb()
  const rgbMatch = resolved.match(/rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)/)
  if (rgbMatch) {
    return `rgba(${rgbMatch[1]}, ${rgbMatch[2]}, ${rgbMatch[3]}, ${alpha})`
  }

  // Let the browser parse everything else (hex, named colors, etc.)
  const el = document.createElement('div')
  el.style.color = resolved
  document.body.appendChild(el)
  const computed = getComputedStyle(el).color
  document.body.removeChild(el)

  const m = computed.match(/rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)/)
  if (!m) return `rgba(0,0,0,${alpha})`
  return `rgba(${m[1]}, ${m[2]}, ${m[3]}, ${alpha})`
}
