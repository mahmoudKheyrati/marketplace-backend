package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/spf13/cast"
)

type VoteHandler struct {
	voteRepo repository.VoteRepo
}

func NewVoteHandler(voteRepo repository.VoteRepo) *VoteHandler {
	return &VoteHandler{voteRepo}
}

type VoteRequest struct {
	ReviewId int64 `json:"review_id"`
	UpVote   bool  `json:"up_vote"`
	DownVote bool  `json:"down_vote"`
}

func (v *VoteHandler) CreateVote(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)

	var request VoteRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	err = v.voteRepo.CreateVote(ctx, userId, request.ReviewId, request.UpVote, request.DownVote)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not vote for the review"})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

func (v *VoteHandler) DeleteVote(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)
	reviewId := cast.ToInt64(c.Params("reviewId"))

	err := v.voteRepo.DeleteVote(ctx, userId, reviewId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not vote for review"})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
