<picture>
  <source media="(prefers-color-scheme: dark)" srcset="assets/tessera-wordmark-dark.svg">
  <img src="assets/tessera-wordmark.svg" alt="Tessera" width="240">
</picture>

A page builder for Vue where pages are edited by your users and by AI — and every change passes the same validation.

> [!WARNING]
> Early stage — not usable yet. Nothing is published, and everything here may change.

## Why

Letting end users edit a page means accepting page data you did not write. The same is true for AI-generated pages and for imported ones. Tessera treats all three as untrusted input: a page is plain JSON, and it is checked against a schema of allowed blocks and props before it is saved or rendered.

## Planned

- **core** — page schema, validation, sanitization and schema migrations
- **vue** — the editor and the renderer
- **server** — server-side validation of pages before they are stored

## License

[MIT](LICENSE)
