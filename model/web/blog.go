package web

type BlogRequest struct {
	ID      uint   `json:"id"`
	Title   string `json:"title" validate:"required"`
	Summary string `json:"summary" validate:"required"`
	Content string `json:"content" validate:"required"`
	UserID  uint   `json:"user_id" validate:"required,gte=0"`
}

type BlogResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
	CreatedBy string `json:"created_by"`
}

type BlogListResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	CreatedBy string `json:"created_by"`
}
