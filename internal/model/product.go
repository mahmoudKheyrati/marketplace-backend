package model

import (
	"time"
)

type Product struct {
	Id            int64                  `json:"id"`
	CategoryId    int64                  `json:"category_id"`
	Name          string                 `json:"name"`
	Brand         string                 `json:"brand"`
	Description   string                 `json:"description"`
	PictureUrl    string                 `json:"picture_url"`
	Specification map[string]interface{} `json:"specification"` // todo: fix this!  pgtype.JSONB, https://www.alexedwards.net/blog/using-postgresql-jsonb
	CreatedAt     time.Time              `json:"created_at"`
}
