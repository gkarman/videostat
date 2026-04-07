package blogger

type Blogger struct {
	ID         string
	PlatformID int
	URL        string
	events     []any
}

func Create(req CreateBloggerDto) (*Blogger, error) {
	if req.URL == "" {
		return nil, ErrUrlInvalid
	}

	blogger := &Blogger{
		ID:         req.ID,
		PlatformID: req.PlatformID,
		URL:        req.URL,
	}

	event := &Created{
		ID:   blogger.ID,
		URL: blogger.URL,
	}
	blogger.addEvent(event)

	return blogger, nil
}

func (b *Blogger) addEvent(e any) {
	b.events = append(b.events, e)
}

func (b *Blogger) PullEvents() []any {
	evs := b.events
	b.events = nil
	return evs
}
