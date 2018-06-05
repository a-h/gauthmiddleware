# gauthmiddleware

Go Middleware to limit access to Web content (not APIs) using Google Sign-in For Websites (Google Authentication).

It renders a basic login screen for non-authenticated users, then issues a session cookie for logged-in users.

# Installation

```
go get github.com/a-h/gauthmiddleware
```

```go
// Load settings from environment variables or use gauthmiddleware.NewWithConfiguration to customise.
handler, err := gauthmiddleware.New()
```

# Usage

Set the required environment variables:

* SESSION_ENCRYPTION_KEY
    * A Base64 encoded key used to encrypt and decrypt cookies. Should be 32 bytes of data.
* COOKIE_NAME
    * The name used for the session cookie generated by the site once Google Authentication is complete, e.g. `auth-session`.
* SET_SECURE_FLAG
    * The site should only be access via HTTPS. When set to true, session cookies are set with the secure flag. The only reason to set this to `false` is during testing.
* GOOGLE_AUTH_CLIENT_ID
    * The ClientID generated by Google which allows your site to request Google Authentication. Configure this at https://developers.google.com/identity/sign-in/web/sign-in
* GOOGLE_ALLOWED_DOMAINS
    * A comma-separated list of GSuite domains which are allowed access to the content, or an asterisk to allow all.