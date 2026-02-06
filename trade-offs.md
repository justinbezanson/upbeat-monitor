# Architectural Trade-Offs: "Always Free" Optimization

This document analyzes the implications of moving away from standard AWS services (S3, API Gateway) to alternative free options (Netlify/GitHub Pages, Lambda Function URLs).

## 1. Frontend Hosting: S3 + CloudFront vs. Netlify/GitHub Pages

### Option A: AWS S3 + CloudFront (Current Plan)
*   **Pros:**
    *   **Unified Infrastructure:** Everything is in one AWS account (IAM, Billing).
    *   **Control:** Granular control over caching headers and security (WAF).
*   **Cons:**
    *   **Cost:** S3 storage is **not free** after 12 months (though cheap, ~$0.02/GB).
    *   **Complexity:** Requires setting up OAI/OAC, Bucket Policies, and CloudFront Distributions.
    *   **SPA Routing:** Requires specific CloudFront error page configuration to handle Vue.js history mode.

### Option B: Netlify (Recommended Alternative)
*   **Pros:**
    *   **Truly Free:** Generous free tier (100GB bandwidth/month) that doesn't expire.
    *   **Developer Experience:** Connects to Git; auto-builds and deploys on push.
    *   **SPA Support:** Native support for SPA routing (redirects `/*` to `index.html` easily).
    *   **Deploy Previews:** Automatically creates live preview URLs for Pull Requests.
*   **Cons:**
    *   **Vendor Split:** Your frontend is on Netlify, backend on AWS.

### Option C: GitHub Pages
*   **Pros:** Free and integrated with your code repository.
*   **Cons:**
    *   **SPA Routing:** Does not natively support SPA history mode (requires a `404.html` hack).
    *   **Config:** Less flexible than Netlify for headers/redirects.

**Verdict:** **Switch to Netlify.** It removes the S3 cost and significantly simplifies the CI/CD pipeline for both the Astro marketing site and Vue SPA.

---

## 2. Backend API: API Gateway vs. Lambda Function URLs

### Option A: AWS API Gateway HTTP API (Current Plan)
*   **Pros:**
    *   **Throttling:** Can limit users to X requests/second. **Critical** for a Public API to prevent abuse.
    *   **Security:** Built-in protection against DDoS (AWS Shield Standard).
    *   **Custom Domains:** Easy SSL certificate management via AWS Certificate Manager.
    *   **Auth:** Built-in JWT Authorizers offload authentication logic from your code.
*   **Cons:**
    *   **Cost:** ~$1.00 per million requests after the 12-month trial.

### Option B: Lambda Function URLs
*   **Pros:**
    *   **Cost:** Completely free (included in Lambda compute pricing).
    *   **Simplicity:** Just a checkbox to enable an HTTP endpoint for a function.
*   **Cons:**
    *   **No Throttling:** If a user writes a bad script that hits your API 10,000 times a second, your Lambda will scale up to process it, potentially skyrocketing your **Compute** bill.
    *   **CORS Complexity:** CORS must be handled carefully in code or basic config.
    *   **Custom Domains:** Function URLs generate ugly domains (`https://<id>.lambda-url.us-east-1.on.aws/`). To use a custom domain (`api.yoursite.com`), you **must** put CloudFront in front of it, which adds complexity back in.

**Verdict:** **Stick with API Gateway.**
For a SaaS that exposes a Public API to customers, the risk of unthrottled traffic (accidental or malicious) outweighs the $1.00/month savings. API Gateway acts as a safety valve.

## 3. Hybrid Recommendation

To maximize the "Free Tier" benefits while maintaining security:

1.  **Frontend:** Host on **Netlify** (Free, Easy).
2.  **Backend:** Keep **API Gateway** (Safe, Scalable).
3.  **Database:** Keep **Turso** (Free, Fast).
4.  **Cron:** Keep **EventBridge + Lambda** (Free, Efficient).

This setup keeps your fixed costs at **$0.00** for the first year, and only ~$1.00/month afterwards, while providing professional-grade security and performance.