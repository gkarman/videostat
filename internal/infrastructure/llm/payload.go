package llm

import "encoding/json"

type AnalysisSummary struct {
	DurationSeconds int    `json:"duration_seconds"`
	Transcript      string `json:"transcript"`
	Summary         any    `json:"summary,omitempty"`
	Chapters        any    `json:"chapters,omitempty"`
	KeyHighlights   any    `json:"key_highlights,omitempty"`
	SentimentPeaks  any    `json:"sentiment_peaks,omitempty"`
	Entities        any    `json:"entities,omitempty"`
}

func ExtractSummary(raw []byte) (*AnalysisSummary, error) {
	var payload struct {
		Text          string `json:"text"`
		AudioDuration int    `json:"audio_duration"`
		Summary       string `json:"summary"`
		Chapters      []struct {
			Start    int    `json:"start"`
			End      int    `json:"end"`
			Headline string `json:"headline"`
			Gist     string `json:"gist"`
			Summary  string `json:"summary"`
		} `json:"chapters"`
		AutoHighlightsResult struct {
			Results []struct {
				Rank  float64 `json:"rank"`
				Text  string  `json:"text"`
				Count int     `json:"count"`
			} `json:"results"`
		} `json:"auto_highlights_result"`
		SentimentAnalysisResults []struct {
			Text       string  `json:"text"`
			Sentiment  string  `json:"sentiment"`
			Confidence float64 `json:"confidence"`
			Start      int     `json:"start"`
			End        int     `json:"end"`
		} `json:"sentiment_analysis_results"`
		Entities []struct {
			Text       string `json:"text"`
			EntityType string `json:"entity_type"`
		} `json:"entities"`
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	s := &AnalysisSummary{
		DurationSeconds: payload.AudioDuration,
		Transcript:      payload.Text,
	}

	if payload.Summary != "" {
		s.Summary = payload.Summary
	}
	if len(payload.Chapters) > 0 {
		s.Chapters = payload.Chapters
	}
	if len(payload.AutoHighlightsResult.Results) > 0 {
		s.KeyHighlights = payload.AutoHighlightsResult.Results
	}
	if len(payload.Entities) > 0 {
		s.Entities = payload.Entities
	}

	var peaks []any
	for _, r := range payload.SentimentAnalysisResults {
		if r.Sentiment != "NEUTRAL" && r.Confidence > 0.8 {
			peaks = append(peaks, r)
		}
	}
	if len(peaks) > 0 {
		s.SentimentPeaks = peaks
	}

	return s, nil
}

const SystemPrompt = `You are an expert content strategist specializing in short-form vertical video (TikTok, YouTube Shorts, Instagram Reels) about retirement planning and personal finance in the USA.

You receive a transcript and analysis of a viral donor video. Your job:
1. Decode the viral formula of the donor — its hook structure, pacing, emotional arc, key facts.
2. Produce a creative brief and a word-for-word avatar script in the same style, adapted for a US retirement finance audience.

Script length rule: target_words = duration_seconds × 2.2 (e.g. 60 s → ~130 words, 90 s → ~200 words).

Output ONLY a valid JSON object — no markdown fences, no explanation, nothing else:

{
  "brief": {
    "core_message": "single most important idea the new video conveys",
    "target_audience": "specific description of who this is for",
    "opening_hook": "exact first 1-2 sentences that stop the scroll",
    "key_talking_points": ["point 1", "point 2", "point 3"],
    "emotional_arc": "how tone and energy shift from start to finish",
    "memorable_phrases": ["phrase 1", "phrase 2"],
    "delivery_style": "pace, tone, energy level — mirror the donor",
    "closing_cta": "final sentence and the one action to drive"
  },
  "script": "Full word-for-word script for the avatar. Plain spoken English only — no bullet points, no markdown, no stage directions, no headers."
}

Script rules:
- Natural spoken language, as if talking directly to camera.
- Mirror the donor's pacing and energy exactly.
- Frame financial insights using US retirement concepts (Social Security, 401(k), IRA, Medicare, RMDs) where relevant.
- Word count must match duration_seconds × 2.2.
- End with a clear, simple call to action.`
