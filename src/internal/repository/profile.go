package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Eugune-Usachev/social-network/src/internal/repository/cache"
	"github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/Eugune-Usachev/social-network/src/pkg/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"
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

func (profileRepository *ProfileRepository) GetSmallProfile(
	ctx context.Context,
	id int,
) (profile model.SmallProfile, err error) {
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
		return model.SmallProfile{}, ErrNotFound
	}

	const query = `SELECT name, second_name, avatar, description, birthday, gender, email FROM users WHERE id = $1`

	row := profileRepository.postgres.QueryRow(ctx, query, id)

	err = row.Scan(
		&profile.Name,
		&profile.SecondName,
		&profile.Avatar,
		&profile.Description,
		&profile.Birthday,
		&profile.Gender,
		&profile.Email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			profileRepository.cache.SetNegativeCase(ctx, fmt.Sprintf(redisKeyNegativeProfile, id))

			return model.SmallProfile{}, ErrNotFound
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

func (profileRepository *ProfileRepository) UpdateSmallProfile(
	ctx context.Context,
	id int,
	profile *model.UpdateSmallProfile,
) (err error) {
	smallProfile := model.SmallProfile{}

	query := strings.Builder{}
	returningPartOfQuery := strings.Builder{}

	query.Grow(150)
	returningPartOfQuery.Grow(75)

	args := make([]any, 0, 8)
	i := 1

	returningArgs := make([]any, 0, 8)
	returningArgs = append(returningArgs, &smallProfile.Avatar)
	returningArgs = append(returningArgs, &smallProfile.Email)
	wasFirstSet := true

	_, _ = query.WriteString("UPDATE users SET ")

	if profile.GetName() != "" {
		smallProfile.Name = profile.GetName()
		wasFirstSet = false
		_, _ = query.WriteString("name = $")
		_, _ = query.WriteString(fmt.Sprintf("%d", i))

		args = append(args, profile.GetName())
		i++
	} else {
		returningPartOfQuery.WriteString(", name")

		returningArgs = append(returningArgs, &smallProfile.Name)
	}

	if profile.GetSecondName() != "" {
		smallProfile.SecondName = profile.GetSecondName()

		if !wasFirstSet {
			_, _ = query.WriteString(", ")
			wasFirstSet = false
		}

		_, _ = query.WriteString("second_name = $")
		_, _ = query.WriteString(fmt.Sprintf("%d", i))

		args = append(args, smallProfile.GetSecondName())
		i++
	} else {
		returningPartOfQuery.WriteString(", second_name")

		returningArgs = append(returningArgs, &smallProfile.SecondName)
	}

	if profile.GetDescription() != "" {
		smallProfile.Description = profile.GetDescription()

		if !wasFirstSet {
			_, _ = query.WriteString(", ")
			wasFirstSet = false
		}

		_, _ = query.WriteString("description = $")
		_, _ = query.WriteString(fmt.Sprintf("%d", i))

		args = append(args, smallProfile.GetDescription())
		i++
	} else {
		returningPartOfQuery.WriteString(", description")

		returningArgs = append(returningArgs, &smallProfile.Description)
	}

	if profile.GetBirthday() != "" {
		smallProfile.Birthday = profile.GetBirthday()

		if !wasFirstSet {
			_, _ = query.WriteString(", ")
			wasFirstSet = false
		}

		_, _ = query.WriteString("birthday = $")
		_, _ = query.WriteString(fmt.Sprintf("%d", i))

		args = append(args, smallProfile.GetBirthday())
		i++
	} else {
		returningPartOfQuery.WriteString(", birthday")

		returningArgs = append(returningArgs, &smallProfile.Birthday)
	}

	if profile.GetGender() > 0 {
		smallProfile.Gender = profile.GetGender()

		if !wasFirstSet {
			_, _ = query.WriteString(", ")
		}

		_, _ = query.WriteString("gender = $")
		_, _ = query.WriteString(fmt.Sprintf("%d", i))

		args = append(args, smallProfile.GetGender())

		i++
	} else {
		returningPartOfQuery.WriteString(", gender")

		returningArgs = append(returningArgs, &smallProfile.Gender)
	}

	_, _ = query.WriteString(" WHERE id = $")
	_, _ = query.WriteString(fmt.Sprintf("%d", i))

	args = append(args, id)

	_, _ = query.WriteString(" RETURNING avatar, email")

	if len(returningArgs) > 0 {
		_, _ = query.WriteString(returningPartOfQuery.String())
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
		profileRepository.logger.Error(
			fmt.Sprintf("Error has been occurred while marshaling small profile, err: %s, profile: %v",
				err.Error(), &smallProfile),
		)
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
			return "", ErrNotFound
		}

		return
	}

	profileRepository.cache.SetString(ctx, fmt.Sprintf(redisKeyInfo, id), stringInfo)

	return
}

func (profileRepository *ProfileRepository) UpdateInfo(ctx context.Context, id int, info string) error {
	const query = `
        UPDATE users 
        SET info = $1
        WHERE id = $2
    `

	var err error

	_, err = profileRepository.postgres.Exec(ctx, query, info, id)
	if err != nil {
		return err
	}

	profileRepository.cache.SetString(ctx, fmt.Sprintf(redisKeyInfo, id), info)

	return nil
}
