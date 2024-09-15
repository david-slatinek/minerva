package models

// SongDto godoc
// @Description Model for a created song
type SongDto struct {
	Id string `json:"id"`
	Song
} //@name SongDto

func (SongDto) TableName() string {
	return "songs"
}
