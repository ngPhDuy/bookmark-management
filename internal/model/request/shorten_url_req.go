package request

type ShortenURLRequest struct {
	URL string `json:"url" binding:"required,url`
	Exp int    `json:"exp" binding:"required,gt=0"`
}
