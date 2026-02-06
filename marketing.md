# Marketing Website Implementation Plan

## Overview
The marketing site is the public face of the SaaS. Its primary goals are SEO ranking, user education, and conversion (driving sign-ups to the SPA).

*   **Reference:** Structure and content strategy modeled after [UptimeRobot](https://uptimerobot.com/).
*   **Framework:** [Astro](https://astro.build/) (Chosen for "Zero JS" by default, perfect for static content).
*   **Styling:** Tailwind CSS.
*   **Hosting:** Netlify (Free Tier).

## 1. Architecture

Since the site is 99% static, we will use Netlify's global CDN and Git-based workflow.

*   **Source:** Astro project in Git.
*   **Build:** Netlify (Auto-builds on push).
*   **CDN:** Netlify Edge Network (Handles SSL, Caching, and Global Distribution).
*   **Domain:** Root domain (e.g., `monitor-saas.com`). The SPA lives on `app.monitor-saas.com`.

## 2. Site Structure (Sitemap)

The pages are designed to mirror the feature set defined in `features.md`.

### Core Pages
1.  **Home (`/`)**
    *   **Hero:** "World's leading uptime monitoring service" style value prop.
    *   **CTA:** "Start for Free" (Links to SPA Registration).
    *   **Social Proof:** "Trusted by X companies."
    *   **Quick Features:** Icons for HTTP, Ping, Port monitoring.

2.  **Features (`/features`)**
    *   Deep dive into monitoring types (HTTP, Keyword, Ping, Port).
    *   Explanation of "Status Pages" (Public vs Private).
    *   Alerting capabilities (Email, SMS, Webhooks).

3.  **Pricing (`/pricing`)**
    *   **Interactive Toggle:** Monthly vs Yearly billing.
    *   **Comparison Table:**
        *   **Free:** 50 Monitors, 5-min checks.
        *   **Solo:** 10 Monitors, 1-min checks, SMS.
        *   **Team:** 100 Monitors, White-label Status Pages.
        *   **Enterprise:** 200+ Monitors, 30-sec checks.
    *   **FAQ Section:** "What happens if I exceed limits?", "Cancellation policy".

4.  **Integrations (`/integrations`)**
    *   Grid view of supported alerts: Slack, Discord, Telegram, Microsoft Teams, Webhooks.
    *   SEO pages for each (e.g., `/integrations/slack-uptime-monitoring`).

5.  **Status Page Product (`/status-page`)**
    *   Dedicated landing page selling the "Status Page" feature as a standalone value prop (communicating downtime to customers).

6.  **Blog (`/blog`)**
    *   Built using Astro Content Collections (Markdown/MDX).
    *   Topics: "How to monitor SSL", "Best practices for SRE", "UptimeRobot Alternatives".

## 3. Technical Implementation

### A. Astro Components
*   **`Navbar.astro`**: Responsive nav. Includes "Login" (link to SPA) and "Register" (CTA).
*   **`PricingCard.astro`**: Reusable card for the 4 tiers. Accepts props for price, monitor count, and features list.
*   **`Footer.astro`**: SEO links, Legal (Terms, Privacy), Socials.
*   **`BaseHead.astro`**: SEO meta tags, Open Graph images, Canonical URLs.

### B. Integration with Backend
While mostly static, the marketing site can fetch dynamic stats during build time or client-side.
*   **Total Monitors Count:** Fetch from Public API (`/api/v1/stats/global`) to display "Monitoring 1,000,000+ monitors" in the Hero section.
*   **System Status:** A small indicator in the footer showing if the SaaS itself is operational (fetching the SaaS's own status page JSON).

## 4. SEO Strategy

*   **Performance:** Astro ships zero JavaScript for content-heavy pages, ensuring 95+ Lighthouse scores (Core Web Vitals).
*   **Structured Data:** JSON-LD schemas for "SoftwareApplication" and "PricingPlan".
*   **Sitemap:** Auto-generated `sitemap.xml` via `@astrojs/sitemap`.

## 5. Implementation Phases

### Phase 1: Scaffold & Design
*   Initialize Astro + Tailwind.
*   Build Layouts and Landing Page.

### Phase 2: Content & Pricing
*   Implement Pricing page matching `features.md`.
*   Create "Features" sub-pages.

### Phase 3: Blog & Analytics
*   Set up Content Collections for the blog.
*   Integrate privacy-friendly analytics (e.g., Plausible or GA4).