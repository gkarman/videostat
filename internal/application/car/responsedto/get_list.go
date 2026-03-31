package responsedto

type GetList struct {
	Cars []*Car `json:"cars"`
}
