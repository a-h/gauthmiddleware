package tokenverifier

import "encoding/json"

// A Claim represents a subset of fields available in a JWT.
// See (https://tools.ietf.org/html/draft-ietf-oauth-json-web-token-32) and
// https://developers.google.com/identity/sign-in/web/backend-auth#calling-the-tokeninfo-endpoint
type Claim struct {
	// The issuer, should be "https://accounts.google.com" or "accounts.google.com"
	Issuer string `json:"iss"`
	// The expiry, e.g. "1433981953". Should not be in the past.
	Expiry        string `json:"exp"`
	Email         string `json:"email"` // e.g. "testuser@gmail.com",
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`        // e.g. "Test User",
	Picture       string `json:"picture"`     // e.g. "https://lh4.googleusercontent.com/-kYgzyAWpZzJ/ABCDEFGHI/AAAJKLMNOP/tIXL9Ir44LE/s99-c/photo.jpg",
	GivenName     string `json:"given_name"`  // e.g. "Test"
	FamilyName    string `json:"family_name"` // e.g. "User"
	HD            string `json:"hd"`          // e.g. "infinityworks.com" - The GSuite domain.
}

// NewClaim creates an instance of a claim from JWT JSON.
func NewClaim(jwt []byte) (claim *Claim, err error) {
	claim = &Claim{}
	err = json.Unmarshal(jwt, claim)
	return claim, err
}
