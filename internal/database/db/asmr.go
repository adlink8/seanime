package db

import (
	"seanime/internal/database/models"
)

// GetAsmrWorkState 读取作品的本地收藏状态。未找到返回 (nil, false)。
func (db *Database) GetAsmrWorkState(rjID string) (*models.AsmrWorkState, bool) {
	var res models.AsmrWorkState
	err := db.gormdb.Where("rj_id = ?", rjID).First(&res).Error
	if err != nil {
		return nil, false
	}
	return &res, true
}

// UpsertAsmrWorkState 写入/更新作品收藏状态（唯一键 rj_id）。
func (db *Database) UpsertAsmrWorkState(rjID string, favorite bool) error {
	var existing models.AsmrWorkState
	err := db.gormdb.Where("rj_id = ?", rjID).First(&existing).Error
	if err != nil {
		// 不存在 → 新建
		return db.gormdb.Create(&models.AsmrWorkState{RjID: rjID, Favorite: favorite}).Error
	}
	existing.Favorite = favorite
	return db.gormdb.Save(&existing).Error
}

// GetAsmrTrackStates 读取某作品下全部音轨完听记录。
func (db *Database) GetAsmrTrackStates(rjID string) ([]*models.AsmrTrackState, error) {
	var res []*models.AsmrTrackState
	err := db.gormdb.Where("rj_id = ?", rjID).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpsertAsmrTrackState 写入/更新单条音轨完听状态，返回该作品下已完成音轨数（listenedCount）。
// 唯一键 rj_id + track_path。
func (db *Database) UpsertAsmrTrackState(rjID, trackPath string, completed bool) (int, error) {
	var existing models.AsmrTrackState
	err := db.gormdb.Where("rj_id = ? AND track_path = ?", rjID, trackPath).First(&existing).Error
	if err != nil {
		if crerr := db.gormdb.Create(&models.AsmrTrackState{RjID: rjID, TrackPath: trackPath, Completed: completed}).Error; crerr != nil {
			return 0, crerr
		}
	} else {
		existing.Completed = completed
		if saverr := db.gormdb.Save(&existing).Error; saverr != nil {
			return 0, saverr
		}
	}

	var count int64
	if err := db.gormdb.Model(&models.AsmrTrackState{}).
		Where("rj_id = ? AND completed = ?", rjID, true).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// IsAsmrWorkSeen 查询新作跟踪器是否已评估过该作品（契约 03.2c / D3②）。
func (db *Database) IsAsmrWorkSeen(rjID string) (bool, error) {
	var count int64
	err := db.gormdb.Model(&models.AsmrTrackerSeen{}).Where("rj_id = ?", rjID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// MarkAsmrWorkSeen 标记作品已被跟踪器评估，幂等（已存在不新增行、不报错）。
func (db *Database) MarkAsmrWorkSeen(rjID string) error {
	var count int64
	if err := db.gormdb.Model(&models.AsmrTrackerSeen{}).Where("rj_id = ?", rjID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return db.gormdb.Create(&models.AsmrTrackerSeen{RjID: rjID}).Error
}

// CountAsmrTrackerSeen 已评估记录总数（tracker 首跑观察窗判定，契约 03.2c / D5）。
func (db *Database) CountAsmrTrackerSeen() (int64, error) {
	var count int64
	err := db.gormdb.Model(&models.AsmrTrackerSeen{}).Count(&count).Error
	return count, err
}
