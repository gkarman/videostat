package command

//type FetchBloggerVideos struct {
//	repo      blogger.Repo
//	vSearcher blogger2.VideoSearcher
//}
//
//func NewFetchBloggerVideos(repo blogger.Repo, vSearchers blogger2.VideoSearcher) *FetchBloggerVideos {
//	return &FetchBloggerVideos{
//		repo:      repo,
//		vSearcher: vSearchers,
//	}
//}
//
//func (c *FetchBloggerVideos) Execute(ctx context.Context, req reqdto.FetchBloggerVideos) error {
//	log := logger.FromContext(ctx).With(
//		"component", "FetchBloggerVideos",
//		"BloggerId", req.BloggerID,
//	)
//
//	b, err := c.repo.GetById(ctx, req.BloggerID)
//	if err != nil {
//		log.Error("get blogger by id failed", "err", err)
//		return fmt.Errorf("get blogger by id: %w", err)
//	}
//
//	videos, err := c.vSearcher.Search(ctx, b)
//	if err != nil {
//		log.Error("video search failed", "err", err)
//		return fmt.Errorf("video search: %w", err)
//	}
//
//	for _, v := range videos {
//		err := c.repo.SaveVideo(ctx, v)
//		if err != nil {
//			log.Error("save video failed", "videoID", v.ExternalID, "err", err)
//			continue
//		}
//	}
//
//	log.Info("videos fetched successfully", "count", len(videos))
//	return nil
//}
