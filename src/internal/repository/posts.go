package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var (
	ErrFailedFetchingPosts    = errors.New("failed to fetch posts")
	ErrFailedFetchingUserRate = errors.New("failed to fetch user rate")
)

type PostsRepository struct {
	cassandra *gocql.Session
	logger    logger.Logger
}

var _ Post = (*PostsRepository)(nil)

func NewPostsRepository(cassandra *gocql.Session, logger logger.Logger) *PostsRepository {
	return &PostsRepository{
		cassandra: cassandra,
		logger:    logger,
	}
}

func (postsRepository *PostsRepository) CreatePost(
	ctx context.Context,
	ownerId int,
	text string,
	survey string,
	files []string,
) (string, error) {
	const insertPostQuery = `
		INSERT INTO posts
		(id, owner_id, text, survey, created_at, updated_at, files, likes, dislikes)
		VALUES (?, ?, ?, ?, toTimestamp(now()), toTimestamp(now()), ?, 0, 0)
	`

	id := uuid.New()
	query := postsRepository.cassandra.Query(insertPostQuery)
	defer query.Release()

	err := query.WithContext(ctx).Bind(
		id,
		ownerId,
		text,
		survey,
		files,
	).Exec()
	if err != nil {
		postsRepository.logger.Error("Error occurred when creating post: " + err.Error())

		return "", err
	}

	return id.String(), nil
}

func (postsRepository *PostsRepository) getUserRate(ctx context.Context, postID, userID string) (int32, error) {
	const selectRateQuery = `
		SELECT is_like
		FROM posts_rates
		WHERE post_id = ? AND user_id = ?
	`

	query := postsRepository.cassandra.Query(selectRateQuery)
	defer query.Release()

	var isLike bool
	if err := query.Bind(postID, userID).WithContext(ctx).Scan(&isLike); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return 0, nil // Not rated
		}

		postsRepository.logger.Error("Error occurred when fetching user rate: " + err.Error())

		return 0, ErrFailedFetchingUserRate
	}
	if isLike {
		return 1, nil // Liked
	}

	return 2, nil // Disliked
}

func (postsRepository *PostsRepository) GetPostsByOwnerID(
	ctx context.Context,
	ownerID string,
	limit int,
	lastQueriedCreateAt time.Time,
) ([]model.Post, error) {
	const getPostsByOwnerIDQuery = `
		SELECT id, owner_id, text, survey, created_at, updated_at, files, likes, dislikes
		FROM posts
		WHERE owner_id = ? AND created_at < ?
		LIMIT ?
	`

	query := postsRepository.cassandra.Query(getPostsByOwnerIDQuery)
	defer query.Release()

	iter := query.WithContext(ctx).Bind(ownerID, limit, lastQueriedCreateAt).WithContext(ctx).Iter()

	posts := []model.Post{}
	var (
		id        gocql.UUID
		text      string
		files     []string
		survey    string
		likes     int64
		dislikes  int64
		createdAt time.Time
		updatedAt time.Time
	)

	for i := 0; iter.Scan(&id, &text, &files, &survey, &likes, &dislikes, &createdAt, &updatedAt); i++ {
		posts = append(posts, model.Post{})

		posts[i].Id = id.String()
		posts[i].Text = text
		posts[i].Files = files
		posts[i].Survey = survey
		posts[i].Likes = int32(likes)
		posts[i].Dislikes = int32(dislikes)
		posts[i].CreatedAt = timestamppb.New(createdAt)
		posts[i].UpdatedAt = timestamppb.New(updatedAt)
	}

	// Check for errors during iteration
	if err := iter.Close(); err != nil {
		postsRepository.logger.Error(fmt.Sprintf("error fetching posts: %v", err))

		return nil, ErrFailedFetchingPosts
	}

	// Fetch user rates if needed (optional step)
	for i := range posts {
		userRate, err := postsRepository.getUserRate(ctx, posts[i].Id, ownerID)
		if err != nil {
			return nil, err
		}

		posts[i].UserRate = userRate
	}

	return posts, nil
}

func (postsRepository *PostsRepository) DeletePost(ctx context.Context, userID int, postID string) error {
	const (
		deletePostQuery = `
			DELETE FROM posts
			WHERE id = ?
			IF owner_id = ?
		`
		deleteAllRatesQuery = `
			DELETE FROM posts_rates
			WHERE post_id = ?
		`
	)

	query := postsRepository.cassandra.Query(deletePostQuery)
	defer query.Release()

	applied, err := query.WithContext(ctx).Bind(postID, userID).ScanCAS()
	if !applied || err != nil {
		if !applied {
			return ErrNotOwner
		}

		if errors.Is(err, gocql.ErrNotFound) {
			return ErrNotFound
		}

		postsRepository.logger.Error("Error occurred when deleting post: " + err.Error())

		return err
	}

	err = postsRepository.cassandra.Query(deleteAllRatesQuery).WithContext(ctx).Bind(postID).Exec()
	if err != nil {
		postsRepository.logger.Error("Error occurred when deleting rates: " + err.Error())

		return err
	}

	return nil
}

func (postsRepository *PostsRepository) deleteLike(ctx context.Context, userId int, postID string) error {
	const (
		deleteRateQuery = `
			DELETE FROM posts_rates
			WHERE post_id = ? AND user_id = ?
		`
		decLikeQuery = `
			UPDATE posts
			SET likes = likes - 1
			WHERE id = ?
		`
	)

	err := postsRepository.cassandra.Query(deleteRateQuery).WithContext(ctx).Bind(postID, userId).Exec()
	if err != nil {
		postsRepository.logger.Error("Error occurred when unrating post (delete rate): " + err.Error())

		return err
	}

	err = postsRepository.cassandra.Query(decLikeQuery).WithContext(ctx).Bind(postID).Exec()
	if err != nil {
		postsRepository.logger.Error("Error occurred when unrating post (dec like): " + err.Error())

		return err
	}

	return nil
}

func (postsRepository *PostsRepository) deleteDislike(ctx context.Context, userId int, postID string) error {
	const (
		deleteRateQuery = `
			DELETE FROM posts_rates
			WHERE post_id = ? AND user_id = ?
		`
		decDislikeQuery = `
			UPDATE posts
			SET dislikes = dislikes - 1
			WHERE id = ?
		`
	)

	err := postsRepository.cassandra.Query(deleteRateQuery).WithContext(ctx).Bind(postID, userId).Exec()
	if err != nil {
		postsRepository.logger.Error("Error occurred when unrating post (delete rate): " + err.Error())

		return err
	}

	err = postsRepository.cassandra.Query(decDislikeQuery).WithContext(ctx).Bind(postID).Exec()
	if err != nil {
		postsRepository.logger.Error("Error occurred when unrating post (dec dislike): " + err.Error())

		return err
	}

	return nil
}

func (postsRepository *PostsRepository) UnratePost(ctx context.Context, userID int, postID string) error {
	const (
		selectPreviousRateQuery = `
			SELECT is_like
			FROM posts_rates
			WHERE post_id = ? AND user_id = ?
		`
	)

	query := postsRepository.cassandra.Query(selectPreviousRateQuery)
	defer query.Release()

	var previousRate map[string]any
	err := query.WithContext(ctx).Bind(postID, userID).Scan(previousRate)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return ErrNotFound
		}

		postsRepository.logger.Error("Error occurred when unrating post: " + err.Error())

		return err
	}

	isLike := previousRate["is_like"].(bool)
	if isLike {
		err = postsRepository.deleteLike(ctx, userID, postID)
	} else {
		err = postsRepository.deleteDislike(ctx, userID, postID)
	}
	if err != nil {
		postsRepository.logger.Error("Error occurred when unrating post: " + err.Error())

		return err
	}

	return nil
}

func (postsRepository *PostsRepository) RatePost(ctx context.Context, userID int, postID string, isLike bool) error {
	const (
		insertRateQuery = `
			INSERT INTO posts_rates
			(post_id, user_id, is_like)
			VALUES (?, ?, ?)
			IF NOT EXISTS
		`
		incLikeQuery = `
			UPDATE posts
			SET likes = likes + 1
			WHERE id = ?
		`
		incDislikeQuery = `
			UPDATE posts
			SET dislikes = dislikes + 1
			WHERE id = ?
		`
	)

	query := postsRepository.cassandra.Query(insertRateQuery)
	defer query.Release()

	var rateRow map[string]interface{}

	applied, err := query.WithContext(ctx).Bind(postID, userID, isLike).MapScanCAS(rateRow)
	if err != nil {
		postsRepository.logger.Error("Error occurred when rating post (try insert rate): " + err.Error())

		return err
	}

	if !applied {
		if rateRow["is_like"].(bool) == isLike {
			return nil
		}

		if rateRow["is_like"].(bool) == true {
			err = postsRepository.deleteLike(ctx, userID, postID)
		} else {
			err = postsRepository.deleteDislike(ctx, userID, postID)
		}

		if err != nil {
			postsRepository.logger.Error("Error occurred when rating post (try delete previous rate): " + err.Error())

			return err
		}

		err = query.WithContext(ctx).Bind(postID, userID, isLike).Exec()
		if err != nil {
			postsRepository.logger.Error("Error occurred when rating post (another try insert rate): " + err.Error())

			return err
		}
	}

	if isLike {
		err = postsRepository.cassandra.Query(incLikeQuery).WithContext(ctx).Bind(postID).Exec()
	} else {
		err = postsRepository.cassandra.Query(incDislikeQuery).WithContext(ctx).Bind(postID).Exec()
	}
	if err != nil {
		postsRepository.logger.Error("Error occurred when rating post: " + err.Error())

		return err
	}

	return nil
}
