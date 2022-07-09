package api

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/spf13/cast"
)

type ReviewHandler struct {
	reviewRepo repository.ReviewRepo
}

func NewReviewHandler(reviewRepo repository.ReviewRepo) *ReviewHandler {
	return &ReviewHandler{reviewRepo: reviewRepo}
}

type CreateReviewRequest struct {
	ProductId  int64   `json:"product_id"`
	StoreId    int64   `json:"store_id"`
	Rate       float64 `json:"rate"`
	ReviewText string  `json:"review_text"`
}

func (r *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)

	var request CreateReviewRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}
	reviewId, err := r.reviewRepo.CreateReview(ctx, userId, request.ProductId, request.StoreId, request.Rate, request.ReviewText)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not create review"})
	}
	return c.JSON(fiber.Map{"status": "ok", "review_id": reviewId})
}

type UpdateReviewRequest struct {
	Rate       float64 `json:"rate"`
	ReviewText string  `json:"review_text"`
}

func (r *ReviewHandler) UpdateReview(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)
	reviewId := cast.ToInt64(c.Params("reviewId"))

	var request UpdateReviewRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	fmt.Println(request)
	id, err := r.reviewRepo.UpdateReview(ctx, userId, reviewId, request.Rate, request.ReviewText)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not update review"})
	}
	return c.JSON(fiber.Map{"status": "ok", "review_id": id})
}
func (r *ReviewHandler) GetUserAllReviews(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)

	reviews, err := r.reviewRepo.GetUserAllReviews(ctx, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not user reviews"})
	}
	return c.JSON(fiber.Map{"status": "ok", "reviews": reviews})
}
func (r *ReviewHandler) DeleteReview(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)
	reviewId := cast.ToInt64(c.Params("reviewId"))

	err := r.reviewRepo.DeleteReview(ctx, userId, reviewId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not delete review"})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

const (
	SortByDate  = "date"
	SortByVotes = "vote"
)

func (r *ReviewHandler) GetProductReviews(c *fiber.Ctx) error {
	ctx := context.Background()
	productId := cast.ToInt64(c.Params("productId"))
	sortedBy := c.Query("sorted_by", SortByDate)

	var reviews = make([]model.Review, 0)
	var err error
	switch sortedBy {
	case SortByVotes:
		reviews, err = r.reviewRepo.GetProductReviewsSortedByVotes(ctx, productId)

	case SortByDate:
		reviews, err = r.reviewRepo.GetProductReviewsSortedByCreatedAt(ctx, productId)

	}
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get product reviews"})
	}
	return c.JSON(fiber.Map{"status": "ok", "reviews": reviews})
}
