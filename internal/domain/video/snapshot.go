package video

type AccountSnapshot struct {
	Account      *Account
	AccountStats *AccountStats

	Contents []*ContentSnapshot
}

type ContentSnapshot struct {
	Content *Content
	Stats   *ContentStats
}