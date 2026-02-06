# Cost Analysis

## Assumptions
*   **Monitors:** 1,000 active monitors.
*   **Check Frequency:** Every 30 seconds (Enterprise/Custom logic) or 1 minute.
*   **API Requests:** 1,000,000 per month (Public + SPA).
*   **Region:** us-east-1 (N. Virginia).
*   **VPC:** Lambda functions are **not** running inside a VPC (Avoiding NAT Gateway costs).

## Monthly Cost Breakdown

| Service | Usage Metric | Estimated Usage | Free Tier Limit | Estimated Cost (Post-Trial) |
| :--- | :--- | :--- | :--- | :--- |
| **AWS Lambda** | Requests | ~1.1M / month | 1M / month | **$0.02** |
| | Compute (GB-s) | ~250k GB-s | 400k GB-s | **$0.00** |
| **API Gateway** | HTTP API Calls | 1M / month | 1M (12 mo only) | **$1.00** |
| **EventBridge** | Scheduler Invokes | ~88k / month | 14M / month | **$0.00** |
| **S3** | Storage | < 1 GB | 5 GB (12 mo only) | **$0.03** |
| **CloudFront** | Data Transfer | < 50 GB | 1 TB | **$0.00** |
| **Turso** | Database Reads | < 1B | 9B | **$0.00** |
| **Total** | | | | **~$1.05 / month** |

## Detailed Notes

### 1. AWS Lambda
*   **Requests:** You have ~1.1M requests (1M API + ~88k Cron triggers). The first 1M are free. You pay $0.20 per million for the excess (~88k), which is negligible.
*   **Compute:** Assuming the Check Runner takes 5 seconds @ 512MB RAM to check a batch of monitors.
    *   `88,000 invocations * 5s * 0.5GB = 220,000 GB-s`.
    *   This is well within the 400,000 GB-s "Always Free" limit.

### 2. API Gateway (HTTP API)
*   **The Main Cost:** API Gateway is only free for the first 12 months. After that, it costs **$1.00 per million requests**.
*   **Optimization:** You can switch to **Lambda Function URLs** (free) for the Public API if you want to save this $1.00, but you lose the throttling and routing features of API Gateway.

### 3. Data Transfer
*   **Outbound:** AWS provides 100GB of free data transfer out to the internet per month.
*   **Monitoring Traffic:** Checking 1,000 sites every minute is roughly `1000 * 2KB * 43200 checks = ~86 GB`. This is close to the limit. If you scale beyond 1,000 monitors, you may start incurring data transfer costs ($0.09/GB).

### 4. Turso (Database)
*   Turso's "Starter" plan is generous (9 Billion reads). Unless you have massive read volume on the Public API, this will remain free.