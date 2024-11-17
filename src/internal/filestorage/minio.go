package filestorage

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Eugune-Usachev/social-network/src/internal/repository"
	"github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	bucketCount = 1 << 8
	bucketMask  = bucketCount - 1
)

var (
	ErrFailedToCheckFileExists      = errors.New("failed to check file existence in database")
	ErrFailedToUploadFile           = errors.New("failed to upload file to MinIO")
	ErrFailedToSaveFileMetadata     = errors.New("failed to save file metadata")
	ErrFailedToGeneratePresignedURL = errors.New("failed to generate presigned URL")
	ErrNotFound                     = errors.New("file not found")
)

// determineBucket determines the bucket name based on the file name.
func determineBucket(fileName string) string {
	hash := sha256.Sum256([]byte(fileName))
	bucketIndex := int(hash[0]) & bucketMask

	return fmt.Sprintf("bucket-%d", bucketIndex)
}

// extractBucketFromFilePath extracts the bucket name from the file path.
func extractBucketFromFilePath(filePath string) string {
	fileName := filepath.Base(filePath)

	return determineBucket(fileName)
}

type MinIOFileStorage struct {
	minioClient *minio.Client
	logger.Logger
	privateFileMetadata repository.PrivateFileMetadata
}

var _ FileStorage = (*MinIOFileStorage)(nil)

func MustNewMinIOFileStorage(
	endpoint,
	accessKeyID,
	secretAccessKey string,
	logger logger.Logger,
	privateFileMetadata repository.PrivateFileMetadata,
) FileStorage {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error occurred when creating minio client: %s", err.Error()))
	}

	fs := &MinIOFileStorage{
		minioClient:         client,
		Logger:              logger,
		privateFileMetadata: privateFileMetadata,
	}
	fs.MustCreateBuckets()

	return fs
}

func (fs MinIOFileStorage) MustCreateBuckets() {
	ctx := context.Background()
	wasCreated := false

	for i := range bucketCount {
		exists, err := fs.minioClient.BucketExists(ctx, fmt.Sprintf("bucket-%d", i))
		if exists && err == nil {
			continue
		}

		if err != nil {
			fs.Logger.Fatal(fmt.Sprintf(
				"Error occurred when checking bucket with number %d existence: %s.",
				i, err.Error(),
			))
		}

		wasCreated = true

		err = fs.minioClient.MakeBucket(ctx, fmt.Sprintf("bucket-%d", i), minio.MakeBucketOptions{})
		if err != nil {
			fs.Logger.Fatal(fmt.Sprintf("Error occurred when creating bucket number %d: %s.", i, err.Error()))
		}
	}

	if !wasCreated {
		fs.Logger.Info("All MinIO buckets have been detected.")
	} else {
		fs.Logger.Info("MinIO buckets have been created successfully.")
	}
}

func (fs MinIOFileStorage) UploadFile(
	ctx context.Context,
	file UploadedFile,
	ownerID int,
	authorizedUserIDs []int,
) (string, error) {
	fileName := filepath.Base(file.FileHeader.Filename)
	if err := validateFileName(fileName); err != nil {
		return "", err
	}

	filePath := generateFileName(fileName, ownerID, len(authorizedUserIDs) > 0)

	bucketName := determineBucket(fileName)

	_, err := fs.minioClient.PutObject(
		ctx,
		bucketName,
		filePath,
		file.File,
		file.FileHeader.Size,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return "", ErrFailedToUploadFile
	}

	// Save metadata for private files
	if len(authorizedUserIDs) > 0 {
		if err := fs.privateFileMetadata.SaveFileMetadata(ctx, filePath, authorizedUserIDs); err != nil {
			errMsg := fmt.Sprintf("failed to save file metadata: %s", err.Error())

			fs.Logger.Error(errMsg)

			return "", ErrFailedToSaveFileMetadata
		}
	}

	return filePath, nil
}

const urlTTL = 5 * time.Minute

func (fs MinIOFileStorage) LoadFile(ctx context.Context, userID int, filePath string) (string, error) {
	isPrivate := strings.HasPrefix(filePath, "private/")

	if isPrivate {
		hasAccess, err := fs.privateFileMetadata.CheckAccess(ctx, filePath, userID)
		if err != nil {
			return "", err
		}

		if !hasAccess {
			return "", ErrForbidden
		}
	}

	// Generate a presigned URL for the file in MinIO
	bucketName := extractBucketFromFilePath(filePath)

	presignedURL, err := fs.minioClient.PresignedGetObject(ctx, bucketName, filePath, urlTTL, nil)
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return "", ErrNotFound
		}

		errMsg := fmt.Sprintf("failed to generate presigned URL: %s", err.Error())

		fs.Logger.Error(errMsg)

		return "", ErrFailedToGeneratePresignedURL
	}

	return presignedURL.String(), nil
}
