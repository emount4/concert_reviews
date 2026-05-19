package favorites_transport_http

type FavoriteRequest struct {
	TargetType string `json:"target_type" validate:"required,oneof=artist venue concert"`
	TargetID   string `json:"target_id" validate:"required"`
}

type FavoriteResponse struct {
	ID         int     `json:"id"`
	TargetType string  `json:"target_type"`
	TargetID   string  `json:"target_id"`
	Name       string  `json:"name"`
	ImageURL   *string `json:"image_url,omitempty"`
	CreatedAt  string  `json:"created_at"`
}

type ListFavoritesResponse struct {
	Items []FavoriteResponse `json:"items"`
}
