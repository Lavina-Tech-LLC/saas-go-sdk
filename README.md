# SaaS Support Go SDK

Official Go SDK for the [SaaS Support](https://saas-support.com) platform. Provides type-safe access to Auth, Billing, and Report modules with zero mandatory dependencies.

```
go get github.com/Lavina-Tech-LLC/saas-support-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    saassupport "github.com/Lavina-Tech-LLC/saas-support-go"
    "github.com/Lavina-Tech-LLC/saas-support-go/auth"
)

func main() {
    client := saassupport.NewClient("pk_live_xxx")

    // Register a user
    result, err := client.Auth.Register(context.Background(), &auth.RegisterParams{
        Email:    "user@example.com",
        Password: "securepassword123",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("User ID: %s\n", result.User.ID)
    fmt.Printf("Token: %s\n", result.AccessToken)
}
```

## Table of Contents

- [Installation](#installation)
- [Client Configuration](#client-configuration)
- [Authentication](#authentication)
  - [Email/Password](#emailpassword)
  - [MFA](#multi-factor-authentication)
  - [OAuth](#oauth)
  - [Magic Link](#magic-link)
  - [Password Reset](#password-reset)
  - [Token Verification](#token-verification)
  - [Session Management](#session-management)
  - [User Profile](#user-profile)
  - [Settings](#settings)
  - [CSRF](#csrf)
- [Organizations](#organizations)
  - [Members](#members)
  - [Invites](#invites)
- [Billing](#billing)
  - [Customers](#customers)
  - [Subscriptions](#subscriptions)
  - [Invoices](#invoices)
  - [Metered Usage](#metered-usage)
  - [Portal](#portal)
  - [Coupons](#coupons)
- [Reports](#reports)
  - [Query Execution](#query-execution)
  - [Saved Queries](#saved-queries)
  - [Dashboards](#dashboards)
  - [Embed Tokens](#embed-tokens)
- [Middleware](#middleware)
  - [net/http](#nethttp-middleware)
  - [Gin](#gin-middleware)
- [Error Handling](#error-handling)
- [Types Reference](#types-reference)

## Installation

```bash
# Core SDK (zero dependencies)
go get github.com/Lavina-Tech-LLC/saas-support-go

# Gin middleware (separate module, only if you use Gin)
go get github.com/Lavina-Tech-LLC/saas-support-go/ginmiddleware
```

**Requirements:** Go 1.21+

The core SDK has zero external dependencies — stdlib only. The `ginmiddleware` package is a separate Go module to avoid pulling Gin as a transitive dependency for non-Gin users.

## Client Configuration

```go
// Default configuration
client := saassupport.NewClient("pk_live_xxx")

// Custom configuration
client := saassupport.NewClient("pk_live_xxx",
    saassupport.WithBaseURL("https://custom-api.example.com/v1"),
    saassupport.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
    saassupport.WithUserAgent("my-app/1.0"),
)
```

| Option | Description | Default |
|--------|-------------|---------|
| `WithBaseURL(url)` | Override the API base URL | `https://api.saas-support.com/v1` |
| `WithHTTPClient(hc)` | Provide a custom `*http.Client` | 30s timeout client |
| `WithUserAgent(ua)` | Set a custom User-Agent header | `saas-support-go/0.1.0` |

The client exposes three service namespaces:

```go
client.Auth    // *auth.Service
client.Billing // *billing.Service
client.Report  // *report.Service
```

### Two-Layer Authentication

The SDK uses two authentication mechanisms simultaneously:

1. **API Key** (`X-API-Key`) — Set once via `NewClient()`. Authenticates your project to the SaaS Support platform. Sent automatically on every request.

2. **User Access Token** (`Authorization: Bearer`) — A short-lived JWT obtained after user login. Required for user-scoped endpoints. Passed explicitly as an `accessToken` parameter to methods that need it.

Both headers coexist on the same request when needed.

---

## Authentication

### Email/Password

**Register**

```go
result, err := client.Auth.Register(ctx, &auth.RegisterParams{
    Email:    "user@example.com",
    Password: "securepassword123",
})
// result.AccessToken, result.RefreshToken, result.User
```

**Login**

```go
result, err := client.Auth.Login(ctx, &auth.LoginParams{
    Email:    "user@example.com",
    Password: "securepassword123",
})
if err != nil {
    log.Fatal(err)
}

if result.MFARequired {
    // User has MFA enabled — prompt for TOTP code, then call LoginMFA
    mfaResult, err := client.Auth.LoginMFA(ctx, &auth.LoginMFAParams{
        MFAToken: result.MFAToken,
        Code:     "123456", // from authenticator app
    })
    // mfaResult.AccessToken, mfaResult.RefreshToken, mfaResult.User
} else {
    // Login complete
    // result.AccessToken, result.RefreshToken, result.User
}
```

### Multi-Factor Authentication

All MFA methods require the user's access token.

**Setup** — Returns a TOTP secret and `otpauth://` URI for QR code generation:

```go
setup, err := client.Auth.MFASetup(ctx, accessToken)
// setup.Secret — base32-encoded TOTP secret
// setup.URI    — otpauth:// URI for QR code rendering
```

**Verify** — Confirms setup with a TOTP code. Returns one-time backup codes:

```go
verified, err := client.Auth.MFAVerify(ctx, accessToken, &auth.MFAVerifyParams{
    Code: "123456",
})
// verified.BackupCodes — []string, store securely
```

**Disable:**

```go
disabled, err := client.Auth.MFADisable(ctx, accessToken, &auth.MFADisableParams{
    Code: "123456",
})
// disabled.Disabled == true
```

### OAuth

**Initiate** — Get the provider's authorization URL:

```go
init, err := client.Auth.OAuthInitiate(ctx, &auth.OAuthInitiateParams{
    Provider:    "google", // or "github"
    RedirectURI: "https://app.example.com/callback",
})
// Redirect user to init.AuthURL
// Store init.State for CSRF validation
```

**Callback** — Exchange the authorization code for tokens:

```go
result, err := client.Auth.OAuthCallback(ctx, &auth.OAuthCallbackParams{
    Provider: "google",
    Code:     code,  // from query params
    State:    state, // validate against stored state
})
// result.AccessToken, result.RefreshToken, result.User
```

### Magic Link

**Send:**

```go
sent, err := client.Auth.MagicLinkSend(ctx, &auth.MagicLinkSendParams{
    Email:       "user@example.com",
    RedirectURL: "https://app.example.com/magic-login",
})
// sent.Sent == true
```

**Verify** — When the user clicks the link:

```go
result, err := client.Auth.MagicLinkVerify(ctx, &auth.MagicLinkVerifyParams{
    Token: tokenFromURL,
})
// result.AccessToken, result.RefreshToken, result.User
```

### Password Reset

**Send reset email:**

```go
sent, err := client.Auth.PasswordResetSend(ctx, &auth.PasswordResetSendParams{
    Email:       "user@example.com",
    RedirectURL: "https://app.example.com/reset-password",
})
```

**Complete reset:**

```go
reset, err := client.Auth.PasswordResetVerify(ctx, &auth.PasswordResetVerifyParams{
    Token:       tokenFromURL,
    NewPassword: "newsecurepassword456",
})
// reset.Reset == true
```

### Token Verification

Server-side verification of user JWTs. This is the primary method for protecting your own API endpoints.

```go
// Convenience method on client
claims, err := client.VerifyToken(ctx, "eyJhbGci...")

// Or via the Auth service directly
claims, err := client.Auth.VerifyToken(ctx, &auth.VerifyTokenParams{
    Token: "eyJhbGci...",
})

if err != nil {
    // Invalid or expired token
}

fmt.Println(claims.Valid)     // true
fmt.Println(claims.UserID)   // "usr_abc123"
fmt.Println(claims.Email)    // "user@example.com"
fmt.Println(claims.ProjectID)// "proj_xyz"
fmt.Println(claims.ExpiresAt)// "2024-01-15T10:30:00Z"
```

### Session Management

**Refresh** — Rotate tokens (single-use refresh tokens):

```go
refreshed, err := client.Auth.Refresh(ctx, &auth.RefreshParams{
    RefreshToken: currentRefreshToken,
})
// refreshed.AccessToken  — new access token
// refreshed.RefreshToken — new refresh token (old one is invalidated)
```

**Logout** — Invalidate the refresh token:

```go
result, err := client.Auth.Logout(ctx, &auth.LogoutParams{
    RefreshToken: refreshToken,
})
// result.LoggedOut == true
```

### User Profile

Requires the user's access token.

```go
// Get current user
user, err := client.Auth.GetMe(ctx, accessToken)
// user.ID, user.Email, user.EmailVerified, user.Provider, user.MFAEnabled, user.Metadata

// Update profile metadata
updated, err := client.Auth.UpdateMe(ctx, accessToken, &auth.UpdateProfileParams{
    Metadata: `{"name": "Jane Doe", "avatar": "https://..."}`,
})
```

### Settings

Public endpoint — returns the project's auth configuration (no user token required).

```go
settings, err := client.Auth.GetSettings(ctx)
// settings.GoogleEnabled      — OAuth with Google
// settings.GitHubEnabled      — OAuth with GitHub
// settings.EmailEnabled       — email/password auth
// settings.MFAEnforced        — MFA required for all users
// settings.PasswordMinLength  — minimum password length
// settings.EmailVerification  — email verification required
```

### CSRF

```go
csrf, err := client.Auth.GetCSRFToken(ctx, accessToken)
// csrf.CSRFToken — include in subsequent state-changing requests
```

---

## Organizations

All organization methods require the user's access token.

```go
// List user's organizations
orgs, err := client.Auth.ListOrgs(ctx, accessToken)

// Create
org, err := client.Auth.CreateOrg(ctx, accessToken, &auth.CreateOrgParams{
    Name: "Acme Corp",
    Slug: "acme-corp",
})

// Get
org, err := client.Auth.GetOrg(ctx, accessToken, "org_abc123")

// Update
org, err := client.Auth.UpdateOrg(ctx, accessToken, "org_abc123", &auth.UpdateOrgParams{
    Name:      "Acme Corporation",
    AvatarURL: "https://example.com/logo.png",
})

// Delete
result, err := client.Auth.DeleteOrg(ctx, accessToken, "org_abc123")
```

### Members

```go
// List members
members, err := client.Auth.ListMembers(ctx, accessToken, "org_abc123")

// Update role
result, err := client.Auth.UpdateMemberRole(ctx, accessToken, "org_abc123", "usr_xyz", &auth.UpdateMemberRoleParams{
    Role: "admin",
})

// Remove member
result, err := client.Auth.RemoveMember(ctx, accessToken, "org_abc123", "usr_xyz")
```

### Invites

```go
// Send invite
invite, err := client.Auth.SendInvite(ctx, accessToken, "org_abc123", &auth.SendInviteParams{
    Email: "colleague@example.com",
    Role:  "member",
})

// Accept invite (by the invited user)
result, err := client.Auth.AcceptInvite(ctx, accessToken, invite.Token)
```

---

## Billing

### Customers

```go
// Create
customer, err := client.Billing.CreateCustomer(ctx, &billing.CreateCustomerParams{
    Email:    "customer@example.com",
    Name:     "Jane Doe",
    Metadata: `{"company": "Acme"}`,
})

// Get
customer, err := client.Billing.GetCustomer(ctx, "cus_abc123")

// Update
customer, err := client.Billing.UpdateCustomer(ctx, "cus_abc123", &billing.UpdateCustomerParams{
    Name: "Jane Smith",
})
```

### Subscriptions

```go
// Subscribe to a plan
sub, err := client.Billing.Subscribe(ctx, "cus_abc123", &billing.SubscribeParams{
    PlanID: "plan_pro_monthly",
})

// Change plan (upgrade/downgrade)
sub, err := client.Billing.ChangePlan(ctx, "cus_abc123", &billing.ChangePlanParams{
    PlanID: "plan_enterprise_monthly",
})

// Cancel (at period end)
result, err := client.Billing.CancelSubscription(ctx, "cus_abc123")
// result.CanceledAtPeriodEnd == true
```

### Invoices

```go
invoices, err := client.Billing.GetCustomerInvoices(ctx, "cus_abc123")
for _, inv := range invoices {
    fmt.Printf("%s — %d cents — %s\n", inv.ID, inv.AmountCents, inv.Status)
    // inv.PdfURL — link to downloadable PDF
}
```

### Metered Usage

```go
// Ingest a usage event
result, err := client.Billing.IngestUsageEvent(ctx, &billing.UsageEventParams{
    CustomerID:     "cus_abc123",
    Metric:         "api_calls",
    Quantity:       150,
    Timestamp:      "2024-01-15T10:30:00Z", // optional, defaults to now
    IdempotencyKey: "evt_unique_123",        // optional, prevents duplicates
})

// Get current period usage
usage, err := client.Billing.GetCurrentUsage(ctx, "cus_abc123")
for _, u := range usage {
    fmt.Printf("%s: %.0f\n", u.Metric, u.Total)
}
```

### Portal

Generate a short-lived token for the Stripe customer portal:

```go
portal, err := client.Billing.CreatePortalToken(ctx, &billing.CreatePortalTokenParams{
    CustomerID: "cus_abc123",
    ExpiresIn:  3600, // seconds, default 3600
})
// portal.PortalToken — embed in portal URL
// portal.ExpiresAt   — expiration time
```

### Coupons

```go
result, err := client.Billing.ApplyCoupon(ctx, "cus_abc123", &billing.ApplyCouponParams{
    Code: "SUMMER20",
})
// result.Applied      — true if successfully applied
// result.DiscountType — "percent" or "fixed"
// result.Amount       — discount amount (20 for 20%, or cents for fixed)
// result.Duration     — "once", "repeating", "forever"
```

---

## Reports

### Query Execution

Execute queries using natural language or raw SQL:

```go
// Natural language query (translated to SQL via AI)
result, err := client.Report.ExecuteQuery(ctx, &report.QueryParams{
    NaturalLanguage: "Show me monthly revenue for the last 6 months",
})

// Raw SQL query
result, err := client.Report.ExecuteQuery(ctx, &report.QueryParams{
    SQL: "SELECT date_trunc('month', created_at) as month, SUM(amount) FROM orders GROUP BY 1 ORDER BY 1",
})

// With filter rules
result, err := client.Report.ExecuteQuery(ctx, &report.QueryParams{
    NaturalLanguage: "Top 10 customers by revenue",
    FilterRules: []report.FilterRule{
        {Table: "orders", Column: "status", Op: "eq", Value: "completed"},
        {Table: "orders", Column: "created_at", Op: "gte", Value: "2024-01-01"},
    },
})

fmt.Printf("SQL: %s\n", result.SQL)
fmt.Printf("Rows: %d (took %dms)\n", result.RowCount, result.ExecutionMs)
fmt.Printf("Chart type: %s\n", result.ChartType)
for _, row := range result.Rows {
    fmt.Println(row) // map[string]interface{}
}
```

**Filter operators:** `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `in`, `between`

### Saved Queries

```go
// List with pagination
page, err := client.Report.ListQueries(ctx, &report.ListParams{
    Page:    1,
    PerPage: 20,
    Sort:    "created_at",
    Order:   "desc",
    Search:  "revenue",
})
// page.Data — []SavedQuery
// page.Meta — {Page, PerPage, Total, TotalPages}

// Save a query
saved, err := client.Report.SaveQuery(ctx, &report.SaveQueryParams{
    Name:            "Monthly Revenue",
    NaturalLanguage: "Show monthly revenue for the last 12 months",
    GeneratedSQL:    "SELECT ...",
    ChartType:       "bar",
})

// Update
saved, err = client.Report.UpdateQuery(ctx, saved.ID, &report.UpdateQueryParams{
    Name:      "Monthly Revenue (Updated)",
    ChartType: "line",
})

// Delete
result, err := client.Report.DeleteQuery(ctx, saved.ID)
```

### Dashboards

```go
// List
page, err := client.Report.ListDashboards(ctx, &report.ListParams{Page: 1, PerPage: 10})

// Create
dashboard, err := client.Report.CreateDashboard(ctx, &report.CreateDashboardParams{
    Name:                   "Executive Dashboard",
    LayoutJSON:             `[{"queryId":"q1","x":0,"y":0,"w":6,"h":4}]`,
    IsPublic:               false,
    RefreshIntervalSeconds: 300,
})

// Get
dashboard, err = client.Report.GetDashboard(ctx, dashboard.ID)

// Update
dashboard, err = client.Report.UpdateDashboard(ctx, dashboard.ID, &report.UpdateDashboardParams{
    Name:     "Executive Dashboard v2",
    IsPublic: true,
})

// Delete
result, err := client.Report.DeleteDashboard(ctx, dashboard.ID)
```

### Embed Tokens

Generate tokens for embedding dashboards in external applications:

```go
// Create embed token with row-level security
embed, err := client.Report.CreateEmbedToken(ctx, &report.CreateEmbedTokenParams{
    DashboardID: "dash_abc123",
    CustomerID:  "cus_xyz",
    ExpiresIn:   3600,
    FilterRules: []report.FilterRule{
        {Table: "orders", Column: "customer_id", Op: "eq", Value: "cus_xyz"},
    },
})
// embed.EmbedToken — use in iframe URL
// embed.ExpiresAt  — token expiration

// List active tokens
tokens, err := client.Report.ListEmbedTokens(ctx)

// Revoke
result, err := client.Report.RevokeEmbedToken(ctx, "tok_abc123")
```

---

## Middleware

Server-side middleware for protecting your API routes using SaaS Support JWT verification.

### net/http Middleware

```go
import (
    saassupport "github.com/Lavina-Tech-LLC/saas-support-go"
    "github.com/Lavina-Tech-LLC/saas-support-go/middleware"
)

client := saassupport.NewClient("pk_live_xxx")

mux := http.NewServeMux()

// Protected route — rejects unauthenticated requests with 401
mux.Handle("/api/profile", middleware.WithAuth(client)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    claims := middleware.GetClaims(r)
    userID := middleware.GetUserID(r)
    email  := middleware.GetEmail(r)

    fmt.Fprintf(w, "Hello %s (ID: %s)", email, userID)
})))

// Optional auth — allows unauthenticated access, populates claims if token present
mux.Handle("/api/feed", middleware.WithOptionalAuth(client)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    claims := middleware.GetClaims(r)
    if claims != nil {
        // Personalized feed
    } else {
        // Public feed
    }
})))
```

**Helpers:**

| Function | Returns | Description |
|----------|---------|-------------|
| `GetClaims(r)` | `*auth.TokenClaims` | Full claims, or `nil` if unauthenticated |
| `GetUserID(r)` | `string` | User ID, or `""` if unauthenticated |
| `GetEmail(r)` | `string` | Email, or `""` if unauthenticated |

### Gin Middleware

Separate module — install independently:

```bash
go get github.com/Lavina-Tech-LLC/saas-support-go/ginmiddleware
```

```go
import (
    saassupport "github.com/Lavina-Tech-LLC/saas-support-go"
    "github.com/Lavina-Tech-LLC/saas-support-go/ginmiddleware"
)

client := saassupport.NewClient("pk_live_xxx")
r := gin.Default()

// Protected group
api := r.Group("/api", ginmiddleware.Auth(client))
{
    api.GET("/profile", func(c *gin.Context) {
        claims := ginmiddleware.GetClaims(c)
        // or: claims := ginmiddleware.MustGetClaims(c) — panics if not authed

        c.JSON(200, gin.H{
            "userId": claims.UserID,
            "email":  claims.Email,
        })
    })
}

// Optional auth
r.GET("/feed", ginmiddleware.OptionalAuth(client), func(c *gin.Context) {
    claims := ginmiddleware.GetClaims(c) // nil if no token
    // ...
})
```

The Gin middleware also sets context values directly:

```go
userID, _ := c.Get("userId") // string
email, _  := c.Get("email")  // string
claims, _ := c.Get("claims") // *auth.TokenClaims
```

**Helpers:**

| Function | Returns | Description |
|----------|---------|-------------|
| `GetClaims(c)` | `*auth.TokenClaims` | Claims or `nil` |
| `MustGetClaims(c)` | `*auth.TokenClaims` | Claims or **panic** (use only behind `Auth`) |

---

## Error Handling

All SDK methods return `(result, error)`. API errors are returned as `*saassupport.Error`:

```go
customer, err := client.Billing.CreateCustomer(ctx, params)
if err != nil {
    // Check specific error types
    if saassupport.IsConflict(err) {
        // 409 — customer already exists
    }
    if saassupport.IsUnauthorized(err) {
        // 401 — bad or expired API key
    }
    if saassupport.IsNotFound(err) {
        // 404 — resource not found
    }
    if saassupport.IsForbidden(err) {
        // 403 — insufficient permissions
    }
    if saassupport.IsRateLimited(err) {
        // 429 — too many requests
    }

    // Access raw error details
    var apiErr *saassupport.Error
    if errors.As(err, &apiErr) {
        fmt.Println(apiErr.Code)            // HTTP status code
        fmt.Println(apiErr.Message)         // API error message
        fmt.Println(string(apiErr.RawBody)) // Full response body
    }

    return err
}
```

Non-API errors (network failures, JSON marshaling) are standard Go errors — use `errors.Is` / `errors.As`.

---

## Types Reference

### Auth Types

| Type | Fields |
|------|--------|
| `User` | `ID`, `Email`, `ProjectID`, `EmailVerified`, `Provider`, `Metadata`, `MFAEnabled` |
| `AuthResult` | `AccessToken`, `RefreshToken`, `User` |
| `LoginResult` | `MFARequired`, `MFAToken`, `AccessToken`, `RefreshToken`, `User` |
| `TokenClaims` | `Valid`, `UserID`, `Email`, `ProjectID`, `ExpiresAt` |
| `Settings` | `GoogleEnabled`, `GitHubEnabled`, `EmailEnabled`, `MFAEnforced`, `PasswordMinLength`, `EmailVerification` |
| `Org` | `ID`, `ProjectID`, `Name`, `Slug`, `AvatarURL`, `Metadata` |
| `Member` | `UserID`, `Email`, `Role` |
| `Invite` | `InviteID`, `Email`, `Role`, `Token`, `ExpiresAt` |
| `MFASetupResult` | `Secret`, `URI` |
| `MFAVerifyResult` | `BackupCodes` |
| `CSRFTokenResult` | `CSRFToken` |

### Billing Types

| Type | Fields |
|------|--------|
| `Customer` | `ID`, `ProjectID`, `Email`, `Name`, `StripeCustomerID`, `BalanceCents`, `Metadata`, `TaxExempt`, `TaxID`, `CreatedAt`, `UpdatedAt` |
| `Subscription` | `ID`, `CustomerID`, `PlanID`, `ProjectID`, `Status`, `StripeSubscriptionID`, `CancelAtPeriodEnd`, `TrialEnd`, `CurrentPeriodStart`, `CurrentPeriodEnd`, `CanceledAt`, `CreatedAt` |
| `Invoice` | `ID`, `ProjectID`, `CustomerID`, `SubscriptionID`, `AmountCents`, `Status`, `StripeInvoiceID`, `PdfURL`, `DueDate`, `PaidAt`, `CreatedAt` |
| `UsageSummary` | `Metric`, `Total` |
| `PortalTokenResult` | `PortalToken`, `ExpiresAt` |
| `ApplyCouponResult` | `Applied`, `DiscountType`, `Amount`, `Duration` |

Subscription statuses: `trialing`, `active`, `past_due`, `paused`, `canceled`

Invoice statuses: `draft`, `open`, `paid`, `void`, `uncollectible`

### Report Types

| Type | Fields |
|------|--------|
| `QueryResult` | `SQL`, `Columns`, `Rows`, `RowCount`, `ExecutionMs`, `ChartType`, `Cached` |
| `SavedQuery` | `ID`, `ProjectID`, `Name`, `NaturalLanguage`, `GeneratedSQL`, `ChartType`, `CreatedAt`, `UpdatedAt` |
| `Dashboard` | `ID`, `ProjectID`, `Name`, `LayoutJSON`, `IsPublic`, `RefreshIntervalSeconds`, `CreatedAt`, `UpdatedAt` |
| `EmbedToken` | `ID`, `ProjectID`, `DashboardID`, `CustomerID`, `FilterRules`, `ExpiresAt`, `CreatedAt`, `RevokedAt` |
| `FilterRule` | `Table`, `Column`, `Op`, `Value` |
| `OffsetPage[T]` | `Data`, `Meta` |
| `OffsetMeta` | `Page`, `PerPage`, `Total`, `TotalPages` |
| `ListParams` | `Page`, `PerPage`, `Sort`, `Order`, `Search` |

---

## Full Example: Backend with Auth

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"

    saassupport "github.com/Lavina-Tech-LLC/saas-support-go"
    "github.com/Lavina-Tech-LLC/saas-support-go/billing"
    "github.com/Lavina-Tech-LLC/saas-support-go/middleware"
)

func main() {
    client := saassupport.NewClient("pk_live_xxx")
    mux := http.NewServeMux()

    // Public: health check
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    // Protected: get user profile
    mux.Handle("/api/me", middleware.WithAuth(client)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims := middleware.GetClaims(r)
        json.NewEncoder(w).Encode(map[string]string{
            "userId": claims.UserID,
            "email":  claims.Email,
        })
    })))

    // Protected: track usage
    mux.Handle("/api/action", middleware.WithAuth(client)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims := middleware.GetClaims(r)

        // Track this API call for metered billing
        _, err := client.Billing.IngestUsageEvent(r.Context(), &billing.UsageEventParams{
            CustomerID: claims.UserID,
            Metric:     "api_calls",
            Quantity:   1,
        })
        if err != nil {
            http.Error(w, "failed to track usage", 500)
            return
        }

        w.Write([]byte("action completed"))
    })))

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## License

MIT
