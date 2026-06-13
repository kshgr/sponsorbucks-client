# Lovable Backend Prompt for SponsorBucks Client

Paste this into the existing SponsorBucks Lovable project after the website update.

---

Add backend support for the SponsorBucks local client.

Create Supabase tables and Edge Functions for device linking, ad delivery, signed event ingestion, session tracking, heartbeat qualification, and admin review.

PRIVACY RULE:
Do not collect prompts, model responses, code, filenames, repository names, terminal output, screenshots, clipboard contents, or window titles. The event schema must not contain fields for these.

TABLES:

1. devices
- id
- user_id nullable until linked
- public_key
- device_name
- os
- arch
- client_version
- status: pending, linked, disabled
- created_at
- linked_at
- last_seen_at

2. device_link_codes
- id
- device_id
- link_code
- expires_at
- claimed_by_user_id nullable
- claimed_at nullable
- status: pending, claimed, expired

3. campaigns
- id
- advertiser_id
- status: draft, pending_review, approved, active, paused, exhausted, rejected
- budget_cents
- spend_cents
- targeting_tools text[]
- daily_cap
- created_at

4. creatives
- id
- campaign_id
- line
- destination_url
- status
- created_at
Constraint: line length <= 60.

5. agent_sessions
- id
- user_id
- device_id
- surface
- started_at
- ended_at
- status
- client_version
- human_initiated
- exit_code
- qualification_summary jsonb

6. impression_events
- id
- session_id
- user_id
- device_id
- campaign_id
- creative_id
- surface
- sequence
- visible_ms
- screen_unlocked
- recent_input_bucket
- foreground_supported_surface
- placement_visible
- signature_valid
- qualified
- rejection_reason
- created_at
- raw_event_hash

7. earnings_ledger
- id
- user_id
- event_id
- amount_cents
- status: pending, approved, rejected, paid
- reason
- created_at
- clears_at

8. fraud_flags
- id
- user_id
- device_id
- session_id
- severity
- reason
- status
- created_at
- reviewed_by
- reviewed_at

EDGE FUNCTIONS:

POST /device-start-link
- Accept public key, device name, OS, arch, client version.
- Create pending device and one-time link code.
- Return device_id, link_code, link_url, expires_at.

POST /device-complete-link
- Accept device_id and link_code.
- If not claimed, return pending.
- If claimed, return linked and a device token.

Website /me page:
- If URL has link_code and user is logged in, claim pending link code and attach device to that user.
- If not logged in, show login/signup inline, then claim after auth.

GET /ads-next
- Requires device token.
- Returns one eligible active creative for the surface.
- Enforce active campaigns, remaining budget, and creative status.
- Return no_ad gracefully when no ad is available.

POST /events-session-start
POST /events-heartbeat
POST /events-session-end
- Require device token.
- Require X-SponsorBucks-Signature and verify signature against stored device public key.
- Store events.
- Qualify heartbeat only if:
  - valid device
  - screen_unlocked true
  - foreground_supported_surface true
  - placement_visible true
  - visible_ms >= 5000
  - recent_input_bucket acceptable
  - no concurrent earning placement for the same user
  - user/device/campaign below caps
- Create pending earnings ledger row only for qualified impressions.
- Rejection should be stored with clear rejection_reason.

ADMIN:
- Add admin pages/tables to inspect devices, sessions, impression events, earnings ledger, and fraud flags.
- Allow disabling a device.
- Allow rejecting pending earnings.
- Add CSV export for sessions/impressions/earnings.

RLS:
- Users can only see their own devices, sessions, impressions, and earnings.
- Advertisers can only see their own campaigns and aggregate campaign stats.
- Admins can see everything.

Keep this implementation simple and robust. No live payout automation yet.
