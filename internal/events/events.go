package events

type SessionStartEvent struct {
	EventType      string `json:"event_type"`
	SessionID      string `json:"session_id"`
	DeviceID       string `json:"device_id"`
	Surface        string `json:"surface"`
	StartedAt      string `json:"started_at"`
	ClientVersion  string `json:"client_version"`
	BuildID        string `json:"build_id"`
	BuildChannel   string `json:"build_channel"`
	HumanInitiated bool   `json:"human_initiated"`
	NoAd           bool   `json:"no_ad,omitempty"`
}

type HeartbeatEvent struct {
	EventType                  string `json:"event_type"`
	SessionID                  string `json:"session_id"`
	DeviceID                   string `json:"device_id"`
	Surface                    string `json:"surface"`
	CampaignID                 string `json:"campaign_id"`
	CreativeID                 string `json:"creative_id"`
	CreativeHash               string `json:"creative_hash"`
	Sequence                   int    `json:"sequence"`
	VisibleMS                  int    `json:"visible_ms"`
	ScreenUnlocked             bool   `json:"screen_unlocked"`
	RecentInputBucket          string `json:"recent_input_bucket"`
	ForegroundSupportedSurface bool   `json:"foreground_supported_surface"`
	PlacementVisible           bool   `json:"placement_visible"`
	CreatedAt                  string `json:"created_at"`
	ClientVersion              string `json:"client_version"`
	BuildID                    string `json:"build_id"`
	BuildChannel               string `json:"build_channel"`
}

type SessionEndEvent struct {
	EventType     string `json:"event_type"`
	SessionID     string `json:"session_id"`
	DeviceID      string `json:"device_id"`
	Surface       string `json:"surface"`
	EndedAt       string `json:"ended_at"`
	ExitCode      int    `json:"exit_code"`
	ClientVersion string `json:"client_version"`
	BuildID       string `json:"build_id"`
	BuildChannel  string `json:"build_channel"`
}

type ClickEvent struct {
	EventType     string `json:"event_type"`
	SessionID     string `json:"session_id"`
	DeviceID      string `json:"device_id"`
	CampaignID    string `json:"campaign_id"`
	CreativeID    string `json:"creative_id"`
	CreativeHash  string `json:"creative_hash"`
	Surface       string `json:"surface"`
	ClickedAt     string `json:"clicked_at"`
	ClientVersion string `json:"client_version"`
	BuildID       string `json:"build_id"`
	BuildChannel  string `json:"build_channel"`
}
