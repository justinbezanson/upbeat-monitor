# Documentation Site Implementation Plan

## Overview
To make the research and planning markdown files easily viewable as a hosted website, we will use **VitePress**. It is a static site generator powered by Vue.js (aligning with your stack) that is optimized for technical documentation.

## 1. Why VitePress?
*   **Zero Config (mostly):** Works with existing Markdown files.
*   **Vue Integration:** You can use Vue components inside your markdown if needed.
*   **Performance:** Extremely fast static site generation.
*   **GitHub Pages:** First-class support for deployment via Actions.

## 2. Setup Instructions

### A. Initialization
Run the following in the repository root:

```bash
npm init -y
npm install -D vitepress vue
npx vitepress init
```

### B. Configuration (`.vitepress/config.mts`)
We will configure the sidebar to link to your existing research files.

```typescript
import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "Server Monitor Research",
  description: "Architecture and Implementation Plans",
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Architecture', link: '/architecture' }
    ],

    sidebar: [
      {
        text: 'Research & Plans',
        items: [
          { text: 'Architecture', link: '/architecture' },
          { text: 'Features & Pricing', link: '/features' },
          { text: 'Trade-offs', link: '/trade-offs' },
          { text: 'Complexity Analysis', link: '/complexity-analysis' },
          { text: 'Cost Analysis', link: '/cost-analysis' },
        ]
      },
      {
        text: 'Implementation',
        items: [
          { text: 'Backend API', link: '/backend-api' },
          { text: 'Frontend SPA', link: '/frontend-spa' },
          { text: 'Cron & Workers', link: '/cron' },
          { text: 'Marketing Site', link: '/marketing' },
          { text: 'Docs Setup', link: '/docs-site' },
          { text: 'Self-Hosting Fizzy', link: '/self-hosting-fizzy' },
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/your-username/your-repo' }
    ]
  }
})
```

## 3. Deployment (GitHub Actions)

Create `.github/workflows/deploy-docs.yml`:

```yaml
name: Deploy Documentation

on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm ci
      - run: npx vitepress build
      - uses: actions/upload-pages-artifact@v3
        with:
          path: .vitepress/dist

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## 4. Final Step
1.  Go to your GitHub Repository Settings -> **Pages**.
2.  Under "Build and deployment", select **GitHub Actions** as the source.
3.  Push the code.