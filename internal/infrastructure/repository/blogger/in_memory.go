package blogger

import (
	"context"
	"sync"

	"github.com/gkarman/demo/internal/domain/blogger"
)

type InMemoryRepo struct {
	mu                   sync.Mutex
	bloggers             map[string]*blogger.Blogger // key = url
	videos               map[string]*blogger.Video
	analyses             []*blogger.VideoAnalysis
	prompts              []*blogger.VideoPrompt
	generations          []*blogger.VideoGeneration
	SaveVideoErrFor      map[string]error // externalID -> error
	SaveVideoAnalysisErr error
}

func NewInMemory() *InMemoryRepo {
	return &InMemoryRepo{
		bloggers:        make(map[string]*blogger.Blogger),
		videos:          make(map[string]*blogger.Video),
		SaveVideoErrFor: make(map[string]error),
	}
}

func (r *InMemoryRepo) Save(_ context.Context, b *blogger.Blogger) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bloggers[b.URL] = b
	return nil
}

func (r *InMemoryRepo) ExistByUrl(_ context.Context, url string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.bloggers[url]
	return ok, nil
}

func (r *InMemoryRepo) GetById(ctx context.Context, id string) (*blogger.Blogger, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, b := range r.bloggers {
		if b.ID == id {
			return b, nil
		}
	}

	return nil, blogger.ErrBloggerNotFound
}

func (r *InMemoryRepo) SaveVideo(ctx context.Context, video *blogger.Video) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err, ok := r.SaveVideoErrFor[video.ExternalID]; ok {
		return err
	}

	r.videos[video.ExternalID] = video
	return nil
}

func (r *InMemoryRepo) ListVideosByBlogger(_ context.Context, bloggerID string) ([]*blogger.Video, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var res []*blogger.Video
	for _, v := range r.videos {
		if v.BloggerID == bloggerID {
			res = append(res, v)
		}
	}
	return res, nil
}
func (r *InMemoryRepo) List(ctx context.Context) ([]*blogger.Blogger, error) {
	return nil, nil
}

func (r *InMemoryRepo) UpdateVideoStatus(
	_ context.Context,
	videoID string,
	from blogger.VideoStatus,
	to blogger.VideoStatus,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var targetVideo *blogger.Video
	for _, v := range r.videos {
		if v.ID == videoID {
			targetVideo = v
			break
		}
	}

	if targetVideo == nil {
		return blogger.ErrVideoNotFound
	}

	if targetVideo.Status != from {
		return blogger.ErrConcurrentUpdate
	}

	targetVideo.Status = to
	return nil
}

func (r *InMemoryRepo) GetVideoByUrl(_ context.Context, url string) (*blogger.Video, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range r.videos {
		if v.URL == url {
			copyVideo := *v
			return &copyVideo, nil
		}
	}

	return nil, blogger.ErrVideoNotFound
}

func (r *InMemoryRepo) GetVideoByID(_ context.Context, id string) (*blogger.Video, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range r.videos {
		if v.ID == id {
			copyVideo := *v
			return &copyVideo, nil
		}
	}

	return nil, blogger.ErrVideoNotFound
}

func (r *InMemoryRepo) UpdateVideoState(_ context.Context, v *blogger.Video) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for extID, stored := range r.videos {
		if stored.ID == v.ID {
			copyVideo := *v
			r.videos[extID] = &copyVideo
			return nil
		}
	}

	return blogger.ErrVideoNotFound
}

func (r *InMemoryRepo) SaveVideoAnalysis(_ context.Context, va *blogger.VideoAnalysis) error {
	if r.SaveVideoAnalysisErr != nil {
		return r.SaveVideoAnalysisErr
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.analyses = append(r.analyses, va)
	return nil
}

func (r *InMemoryRepo) GetSavedAnalyses() []*blogger.VideoAnalysis {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.analyses
}

func (r *InMemoryRepo) GetVideoAnalysisByVideoID(_ context.Context, videoID string) (*blogger.VideoAnalysis, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, va := range r.analyses {
		if va.VideoID == videoID {
			return va, nil
		}
	}
	return nil, blogger.ErrVideoNotFound
}

func (r *InMemoryRepo) GetVideoPromptByVideoID(_ context.Context, videoID string) (*blogger.VideoPrompt, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := len(r.prompts) - 1; i >= 0; i-- {
		if r.prompts[i].VideoID == videoID {
			return r.prompts[i], nil
		}
	}
	return nil, blogger.ErrVideoNotFound
}

func (r *InMemoryRepo) SaveVideoPrompt(_ context.Context, vp *blogger.VideoPrompt) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prompts = append(r.prompts, vp)
	return nil
}

func (r *InMemoryRepo) SaveVideoGeneration(_ context.Context, vg *blogger.VideoGeneration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.generations = append(r.generations, vg)
	return nil
}

func (r *InMemoryRepo) GetVideoGenerationByExternalID(_ context.Context, externalID string) (*blogger.VideoGeneration, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, vg := range r.generations {
		if vg.ExternalID == externalID {
			return vg, nil
		}
	}
	return nil, blogger.ErrVideoNotFound
}

func (r *InMemoryRepo) ListPendingVideoGenerations(_ context.Context) ([]*blogger.VideoGeneration, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*blogger.VideoGeneration
	for _, vg := range r.generations {
		if vg.Status == blogger.VideoGenerationPending || vg.Status == blogger.VideoGenerationProcessing {
			result = append(result, vg)
		}
	}
	return result, nil
}

func (r *InMemoryRepo) UpdateVideoGeneration(_ context.Context, vg *blogger.VideoGeneration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, stored := range r.generations {
		if stored.ID == vg.ID {
			r.generations[i] = vg
			return nil
		}
	}
	return blogger.ErrVideoNotFound
}

func (r *InMemoryRepo) SaveVideoWatcher(_ context.Context, w *blogger.VideoWatcher) error {
	return nil
}

func (r *InMemoryRepo) ListVideoWatchers(_ context.Context, videoID string) ([]*blogger.VideoWatcher, error) {
	return nil, nil
}

func (r *InMemoryRepo) SaveBrollSegments(_ context.Context, videoID string, segments []*blogger.BrollSegment) error {
	return nil
}

func (r *InMemoryRepo) ListBrollSegmentsByVideoID(_ context.Context, videoID string) ([]*blogger.BrollSegment, error) {
	return nil, nil
}

func (r *InMemoryRepo) ListPendingBrollSegments(_ context.Context, videoID string) ([]*blogger.BrollSegment, error) {
	return nil, nil
}

func (r *InMemoryRepo) ListProcessingBrollSegments(_ context.Context) ([]*blogger.BrollSegment, error) {
	return nil, nil
}

func (r *InMemoryRepo) UpdateBrollSegment(_ context.Context, s *blogger.BrollSegment) error {
	return nil
}

func (r *InMemoryRepo) CountActiveBrollSegments(_ context.Context, videoID string) (int, error) {
	return 0, nil
}

func (r *InMemoryRepo) GetVideoGenerationByVideoID(_ context.Context, videoID string) (*blogger.VideoGeneration, error) {
	return nil, blogger.ErrVideoNotFound
}

func (r *InMemoryRepo) SaveVideoComposition(_ context.Context, c *blogger.VideoComposition) error {
	return nil
}

func (r *InMemoryRepo) UpdateVideoComposition(_ context.Context, c *blogger.VideoComposition) error {
	return nil
}

func (r *InMemoryRepo) ListProcessingCompositions(_ context.Context) ([]*blogger.VideoComposition, error) {
	return nil, nil
}

func (r *InMemoryRepo) DeleteBrollSegmentsByVideoID(_ context.Context, _ string) error { return nil }
func (r *InMemoryRepo) DeleteCompositionsByVideoID(_ context.Context, _ string) error  { return nil }
func (r *InMemoryRepo) DeleteVideoGenerationsByVideoID(_ context.Context, _ string) error {
	return nil
}

func (r *InMemoryRepo) DeleteVideoWatchersByVideoID(_ context.Context, _ string) error { return nil }

func (r *InMemoryRepo) ExistsVideoWatcher(_ context.Context, _ string, _ int64) (bool, error) {
	return false, nil
}

func (r *InMemoryRepo) ListVideosWithPendingBrollSegments(_ context.Context) ([]string, error) {
	return nil, nil
}

func (r *InMemoryRepo) ListVideosReadyToCompose(_ context.Context) ([]string, error) {
	return nil, nil
}
