package structs

// GoogleUser is a retrieved and authentiacted google user.
type GoogleUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

// FaceBookUser is a retrieved and authentiacted FaceBook user.
type FaceBookUser struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
