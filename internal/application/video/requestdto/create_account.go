package requestdto

type CreateAccount struct {
	PlatformID int16
	ExternalID string
	Title      string
	URL        string
}
