package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type ReviewRepo interface {
	CreateReview(ctx context.Context, userId, productId, storeId int64, rate float64, reviewText string) (int64, error)
	UpdateReview(ctx context.Context, userId, reviewId int64, rate float64, reviewText string) (int64, error)
	GetUserAllReviews(ctx context.Context, userId int64) ([]model.Review, error)
	DeleteReview(ctx context.Context, userId, reviewId int64) error
	GetProductReviewsSortedByCreatedAt(ctx context.Context, productId int64) ([]model.Review, error)
	GetProductReviewsSortedByVotes(ctx context.Context, productId int64) ([]model.Review, error)
}

type ReviewRepoImpl struct {
	db *pgxpool.Pool
}

func NewReviewRepoImpl(db *pgxpool.Pool) *ReviewRepoImpl {
	return &ReviewRepoImpl{db: db}
}

func (r *ReviewRepoImpl) CreateReview(ctx context.Context, userId, productId, storeId int64, rate float64, reviewText string) (int64, error) {
	query := `insert into review (product_id, store_id, user_id, rate, review_text)
values ($1, $2, $3, $4, $5) returning id`
	var id int64 = -1

	rows := r.db.QueryRow(ctx, query, productId, storeId, userId, rate, reviewText)
	err := rows.Scan(&id)
	return id, err
}

func (r *ReviewRepoImpl) UpdateReview(ctx context.Context, userId, reviewId int64, rate float64, reviewText string) (int64, error) {
	query := `
update review
set rate = $1, review_text = $2
where id = $3
  and user_id = $4
  and deleted_at is null returning id
`
	var id int64 = -1
	rows := r.db.QueryRow(ctx, query, rate, reviewText, reviewId, userId)

	err := rows.Scan(&id)
	return id, err
}

func (r *ReviewRepoImpl) GetUserAllReviews(ctx context.Context, userId int64) ([]model.Review, error) {
	query := `
select id, product_id, store_id, user_id, rate, review_text, created_at
from review
where user_id = $1
  and deleted_at is null
`
	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	var reviews = make([]model.Review, 0)
	for rows.Next() {
		var review model.Review
		err := rows.Scan(
			&review.Id,
			&review.ProductId,
			&review.StoreId,
			&review.UserId,
			&review.Rate,
			&review.ReviewText,
			&review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *ReviewRepoImpl) DeleteReview(ctx context.Context, userId, reviewId int64) error {
	query := `update review
set deleted_at = now()
where id = $1
  and user_id = $2
  and deleted_at is null
`
	_, err := r.db.Query(ctx, query, reviewId, userId)
	return err
}

func (r *ReviewRepoImpl) GetProductReviewsSortedByCreatedAt(ctx context.Context, productId int64) ([]model.Review, error) {
	query := `
select r.id,
       r.product_id,
       r.store_id,
       r.user_id,
       r.rate,
       r.review_text,
       r.created_at,
       sum(case when v.up_vote = true then 1 else 0 end)   as up_votes,
       sum(case when v.down_vote = true then 1 else 0 end) as down_votes
from review r
         left join votes v on r.id = v.review_id
where product_id = $1
  and deleted_at is null
group by r.id, r.product_id, r.store_id, r.user_id, r.rate, r.review_text, r.created_at
order by created_at desc
`
	rows, err := r.db.Query(ctx, query, productId)
	if err != nil {
		return nil, err
	}

	var reviews = make([]model.Review, 0)
	for rows.Next() {
		var review model.Review
		err := rows.Scan(
			&review.Id,
			&review.ProductId,
			&review.StoreId,
			&review.UserId,
			&review.Rate,
			&review.ReviewText,
			&review.CreatedAt,
			&review.UpVotes,
			&review.DownVotes)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *ReviewRepoImpl) GetProductReviewsSortedByVotes(ctx context.Context, productId int64) ([]model.Review, error) {
	query := `
with cte as (
    select r.id,
           r.product_id,
           r.store_id,
           r.user_id,
           r.rate,
           r.review_text,
           r.created_at,
           r.deleted_at,
           sum(case when v.up_vote = true then 1 else 0 end)   as up_votes,
           sum(case when v.down_vote = true then 1 else 0 end) as down_votes
    from review r
             left join votes v on r.id = v.review_id
    group by r.id, r.product_id, r.store_id, r.user_id, r.rate, r.review_text, r.created_at
)
select id,
       product_id,
       store_id,
       user_id,
       rate,
       review_text,
       created_at,
       up_votes,
       down_votes
from cte
where product_id = $1
  and deleted_at is null
order by up_votes * 2 + down_votes desc;


`
	rows, err := r.db.Query(ctx, query, productId)
	if err != nil {
		return nil, err
	}

	var reviews = make([]model.Review, 0)
	for rows.Next() {
		var review model.Review
		err := rows.Scan(
			&review.Id,
			&review.ProductId,
			&review.StoreId,
			&review.UserId,
			&review.Rate,
			&review.ReviewText,
			&review.CreatedAt,
			&review.UpVotes,
			&review.DownVotes)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}
