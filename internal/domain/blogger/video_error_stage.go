package blogger

type VideoErrorStage string

const (
	ErrorStageFileFetch VideoErrorStage = "file_fetch"
	ErrorStageAnalysis  VideoErrorStage = "analysis"
	ErrorStageInsights  VideoErrorStage = "insights"
)
