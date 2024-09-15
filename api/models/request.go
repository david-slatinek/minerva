package models

// Song godoc
// @Description Model for creating a new song
type Song struct {
	Title    string `json:"title" binding:"required"`
	Duration string `json:"duration" binding:"required"`
	Release  string `json:"release" binding:"required"`
	Author   string `json:"author" binding:"required"`
} //@name Song
