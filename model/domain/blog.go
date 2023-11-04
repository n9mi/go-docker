package domain

type Blog struct {
	ID      uint `gorm:"primaryKey"`
	Title   string
	Summary string
	Content string
	UserID  uint
	User    User
}

type ScanBlogs struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Summary   string
	CreatedBy string
}

type ScanBlog struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Summary   string
	Content   string
	CreatedBy string
}
