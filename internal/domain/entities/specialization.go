package entities

type Specialization struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Title string `gorm:"unique;not null" json:"title" example:"Терапевт"`
}
