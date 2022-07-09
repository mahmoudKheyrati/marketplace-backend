package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type VoteRepo interface {
	CreateVote(ctx context.Context, userId, reviewId int64, upVote, downVote bool) error
	DeleteVote(ctx context.Context, userId, reviewId int64) error
}

type VoteRepoImpl struct {
	db *pgxpool.Pool
}

func NewVoteRepoImpl(db *pgxpool.Pool) *VoteRepoImpl {
	return &VoteRepoImpl{db: db}
}

func (v *VoteRepoImpl) CreateVote(ctx context.Context, userId, reviewId int64, upVote, downVote bool) error {
	query := `insert into votes (review_id, user_id, up_vote, down_vote)
			   values ($1, $2, $3, $4)`
	_, err := v.db.Query(ctx, query, reviewId, userId, upVote, downVote)
	return err
}
func (v *VoteRepoImpl) DeleteVote(ctx context.Context, userId, reviewId int64) error {
	query := `
	delete
	from votes
	where user_id = $1
	  and review_id = $2
`
	_, err := v.db.Query(ctx, query, userId, reviewId)
	return err
}
