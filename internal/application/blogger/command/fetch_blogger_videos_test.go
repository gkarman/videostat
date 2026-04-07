package command

//func TestFetchBloggerVideos_Execute(t *testing.T) {
//	ctx := context.Background()
//
//	repo := blogger.NewInMemory()
//	bloggerID := "b1"
//	repo.Bloggers[bloggerID] = &blogger.Blogger{
//		ID:         bloggerID,
//		PlatformID: 1,
//		URL:        "https://youtube.com/@test",
//	}
//
//	// Подготовка FakeVideoSearcher
//	videos := []*blogger.Video{
//		{ExternalID: "v1", URL: "https://youtube.com/v1", Title: "Video 1", Views: 1000},
//		{ExternalID: "v2", URL: "https://youtube.com/v2", Title: "Video 2", Views: 5000},
//	}
//	searcher := blogger.NewFakeVideoSearcher(videos, nil)
//
//	// Команда
//	cmd := command.NewFetchBloggerVideos(repo, searcher)
//
//	// Запуск
//	err := cmd.Execute(ctx, reqdto.FetchBloggerVideos{BloggerID: bloggerID})
//	require.NoError(t, err)
//
//	// Проверяем, что видео сохранены в InMemory
//	savedVideos := repo.ListVideosByBlogger(bloggerID)
//	require.Len(t, savedVideos, 2)
//
//	require.Equal(t, "v1", savedVideos[0].ExternalID)
//	require.Equal(t, "v2", savedVideos[1].ExternalID)
//}
