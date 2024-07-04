package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Eugune-Usachev/social-network/src/customErrors"
	"github.com/Eugune-Usachev/social-network/src/internal/repository/cache"
	"github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"
	"strings"
)

type ProfileRepository struct {
	postgres *pgxpool.Pool
	cache    cache.Cache
	logger   logger.Logger
}

var _ Profile = (*ProfileRepository)(nil)

func NewProfileRepository(postgres *pgxpool.Pool, cache cache.Cache, logger logger.Logger) *ProfileRepository {
	return &ProfileRepository{
		postgres: postgres,
		cache:    cache,
		logger:   logger,
	}
}

const (
	redisKeyProfile         = "profile:%d"
	redisKeyNegativeProfile = "nprofile:%d"
	redisKeyInfo            = "info:%d"
)

func (profileRepository *ProfileRepository) GetSmallProfile(ctx context.Context, id int) (profile model.SmallProfile, err error) {
	var (
		isExist bool
		bytes   []byte
	)
	bytes, isExist = profileRepository.cache.GetBytes(ctx, fmt.Sprintf(redisKeyProfile, id))

	if isExist {
		err = proto.Unmarshal(bytes, &profile)
		return
	}

	isNegativeCase := profileRepository.cache.IsNegativeCase(ctx, fmt.Sprintf(redisKeyNegativeProfile, id))
	if isNegativeCase {
		return model.SmallProfile{}, customErrors.NotFound
	}

	const query = `SELECT name, second_name, avatar, description, birthday, gender, email FROM users WHERE id = $1`

	row := profileRepository.postgres.QueryRow(ctx, query, id)
	err = row.Scan(&profile.Name, &profile.SecondName, &profile.Avatar, &profile.Description, &profile.Birthday, &profile.Gender, &profile.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			profileRepository.cache.SetNegativeCase(ctx, fmt.Sprintf(redisKeyNegativeProfile, id))
			return model.SmallProfile{}, customErrors.NotFound
		}
		return
	}

	bytes, err = proto.Marshal(&profile)
	if err != nil {
		return
	}

	profileRepository.cache.SetBytes(ctx, fmt.Sprintf(redisKeyProfile, id), bytes)

	return
}

func (profileRepository *ProfileRepository) UpdateSmallProfile(ctx context.Context, id int, profile *model.UpdateSmallProfile) (err error) {
	smallProfile := model.SmallProfile{}

	query := strings.Builder{}
	returningPartOfQuery := strings.Builder{}

	query.Grow(150)
	returningPartOfQuery.Grow(75)

	args := make([]any, 0, 8)
	returningArgs := make([]any, 0, 8)

	i := 1
	returningArgs = append(returningArgs, &smallProfile.Avatar)

	_, _ = query.Write([]byte("UPDATE users SET "))
	if profile.Name != "" {
		smallProfile.Name = profile.Name
		_, _ = query.Write([]byte("name = $"))
		_, _ = query.Write([]byte(fmt.Sprintf("%d", i)))
		args = append(args, profile.Name)
		i++
	} else {
		returningPartOfQuery.WriteString(", name")
		returningArgs = append(returningArgs, &profile.Name)
	}

	if profile.SecondName != "" {
		smallProfile.SecondName = profile.SecondName
		_, _ = query.Write([]byte(", second_name = $"))
		_, _ = query.Write([]byte(fmt.Sprintf("%d", i)))
		args = append(args, smallProfile.SecondName)
		i++
	} else {
		returningPartOfQuery.WriteString(", second_name")
		returningArgs = append(returningArgs, &profile.SecondName)
	}

	if profile.Description != "" {
		smallProfile.Description = profile.Description
		_, _ = query.Write([]byte(", description = $"))
		_, _ = query.Write([]byte(fmt.Sprintf("%d", i)))
		args = append(args, smallProfile.Description)
		i++
	} else {
		returningPartOfQuery.WriteString(", description")
		returningArgs = append(returningArgs, &profile.Description)
	}

	if profile.Birthday != "" {
		smallProfile.Birthday = profile.Birthday
		_, _ = query.Write([]byte(", birthday = $"))
		_, _ = query.Write([]byte(fmt.Sprintf("%d", i)))
		args = append(args, smallProfile.Birthday)
		i++
	} else {
		returningPartOfQuery.WriteString(", birthday")
		returningArgs = append(returningArgs, &profile.Birthday)
	}

	if profile.Gender != -1 {
		smallProfile.Gender = profile.Gender
		_, _ = query.Write([]byte(", gender = $"))
		_, _ = query.Write([]byte(fmt.Sprintf("%d", i)))
		args = append(args, smallProfile.Gender)
		i++
	} else {
		returningPartOfQuery.WriteString(", gender")
		returningArgs = append(returningArgs, &profile.Gender)
	}

	if profile.Email != "" {
		smallProfile.Email = profile.Email
		_, _ = query.Write([]byte(", email = $"))
		_, _ = query.Write([]byte(fmt.Sprintf("%d", i)))
		args = append(args, smallProfile.Email)
		i++
	} else {
		returningPartOfQuery.WriteString(", email")
		returningArgs = append(returningArgs, &profile.Email)
	}

	_, _ = query.Write([]byte(" WHERE id = $"))
	_, _ = query.Write([]byte(fmt.Sprintf("%d", i)))
	args = append(args, id)

	_, _ = query.Write([]byte(" RETURNING avatar"))

	if len(returningArgs) > 0 {
		_, _ = query.Write([]byte(returningPartOfQuery.String()))
	}

	row := profileRepository.postgres.QueryRow(
		ctx,
		query.String(),
		args...,
	)

	if err = row.Scan(returningArgs...); err != nil {
		return err
	}

	bytes, err := proto.Marshal(&smallProfile)
	if err != nil {
		return err
	}

	profileRepository.cache.SetBytes(ctx, fmt.Sprintf(redisKeyProfile, id), bytes)

	return nil
}

func (profileRepository *ProfileRepository) GetInfo(ctx context.Context, id int) (info string, err error) {
	var (
		isExist    bool
		stringInfo string
	)

	stringInfo, isExist = profileRepository.cache.GetString(ctx, fmt.Sprintf(redisKeyInfo, id))

	if isExist {
		return stringInfo, nil
	}

	const query = `SELECT info FROM users WHERE id = $1`

	row := profileRepository.postgres.QueryRow(ctx, query, id)
	err = row.Scan(&stringInfo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", customErrors.NotFound
		}
		return
	}

	profileRepository.cache.SetString(ctx, fmt.Sprintf(redisKeyInfo, id), stringInfo)

	return
}

func (profileRepository *ProfileRepository) UpdateInfo(ctx context.Context, id int, info string) (err error) {
	const query = `
        UPDATE users 
        SET info = $1
        WHERE id = $2
    `
	_, err = profileRepository.postgres.Exec(ctx, query, info, id)

	if err != nil {
		return err
	}

	profileRepository.cache.SetString(ctx, fmt.Sprintf(redisKeyInfo, id), info)

	return nil
}
