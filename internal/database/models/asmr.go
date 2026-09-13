package models

// AsmrTrackState 单个音轨的完听状态（契约 §DB）。
// 唯一键 rj_id + track_path：同一作品下不同音轨各自记录。
type AsmrTrackState struct {
	BaseModel
	RjID      string `gorm:"uniqueIndex:idx_asmr_track_rj_track" json:"rjId"`
	TrackPath string `gorm:"uniqueIndex:idx_asmr_track_rj_track" json:"trackPath"`
	Completed bool   `json:"completed"`
}

func (AsmrTrackState) TableName() string { return "asmr_track_states" }

// AsmrWorkState 单个作品的收藏状态（契约 §DB）。
// RjID 唯一索引：一个作品一条收藏记录。
type AsmrWorkState struct {
	BaseModel
	RjID     string `gorm:"uniqueIndex" json:"rjId"`
	Favorite bool   `json:"favorite"`
}

func (AsmrWorkState) TableName() string { return "asmr_work_states" }

// AsmrTrackerSeen 新作跟踪器已评估记录（契约 03.2c / D3）。
// RjID 唯一索引：标记该 RJ 已被 tracker 评估过（含首轮观察窗只记录的情形），防重复通知/评估。
type AsmrTrackerSeen struct {
	BaseModel
	RjID string `gorm:"uniqueIndex" json:"rjId"`
}

func (AsmrTrackerSeen) TableName() string { return "asmr_tracker_seen" }
