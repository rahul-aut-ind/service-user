package imageservice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rahul-aut-ind/service-user/domain/errors"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/repositories/dynamorepo"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/repositories/s3repo"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/requestparser"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type (
	UserImageService interface {
		SaveUserImage(uID string, req *requestparser.MultiPartData) (*models.UploadResponse, error)
		GetAllUserImages(req models.PaginatedInput) (*models.PaginatedImageResponse, error)
		GetByUserIDImageID(uID, imageID string) (*models.ImageResponse, error)
		DeleteByUserIDImageID(uID, imageID string) error
		DeleteAllByUserID(uID string) error
	}

	Service struct {
		db  dynamorepo.DataHandler
		s3  s3repo.S3Handler
		log *logger.Logger
	}
)

func New(db dynamorepo.DataHandler, s3 s3repo.S3Handler, l *logger.Logger) *Service {
	return &Service{db, s3, l}
}

func (s *Service) SaveUserImage(uID string, req *requestparser.MultiPartData) (*models.UploadResponse, error) {

	imageID, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("uuid generation failed"))
	}

	s3Path, err := s.s3.Save(uID, imageID, req.Image.Ext, &req.Image.Bytes)
	if err != nil {
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error uploading image to S3"))
	}

	ui := &models.UserImage{
		UserID:    uID,
		ImageID:   imageID.String(),
		Path:      s3Path,
		IsDeleted: false,
		TakenAt:   req.Metadata.TakenAt,
		UpdatedAt: time.Now(),
	}

	err = s.db.CreateOrUpdateImage(ui)
	if err != nil {
		return nil, err
	}

	return &models.UploadResponse{ID: imageID.String()}, nil
}

func (s *Service) GetAllUserImages(req models.PaginatedInput) (*models.PaginatedImageResponse, error) {
	images, err := s.db.GetAllImagesPaginated(req)

	if err != nil {
		return nil, err
	}

	r := make([]models.ImageResponse, 0, len(images.UserImages))
	for i := range images.UserImages {
		res := models.ImageResponse{
			ImageID: images.UserImages[i].ImageID,
			TakenAt: images.UserImages[i].TakenAt,
			Path:    images.UserImages[i].Path,
		}
		r = append(r, res)
	}

	return &models.PaginatedImageResponse{
		Images: r,
		Page:   images.Page,
	}, nil
}

func (s *Service) GetByUserIDImageID(uID, imageID string) (*models.ImageResponse, error) {

	data, err := s.db.GetImage(uID, imageID)
	if err != nil {
		return nil, err
	}

	return &models.ImageResponse{
		ImageID: data.ImageID,
		TakenAt: data.TakenAt,
		Path:    data.Path,
	}, nil
}

func (s *Service) DeleteByUserIDImageID(uID, imageID string) error {
	return s.parallelDeleteTasks(func() error { return s.s3.Delete(uID, imageID) }, func() error { return s.db.DeleteImage(uID, imageID) })
}

func (s *Service) DeleteAllByUserID(uID string) error {
	return s.parallelDeleteTasks(func() error { return s.s3.DeleteAll(uID) }, func() error { return s.db.DeleteAllImages(uID) })
}

func (s *Service) parallelDeleteTasks(deleteFuncs ...func() error) error {
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	for _, deleteFunc := range deleteFuncs {
		func(f func() error) {
			g.Go(func() error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					if err := f(); err != nil {
						return fmt.Errorf("task failed: %w", err)
					}
					return nil
				}
			})
		}(deleteFunc)
	}
	return g.Wait()
}
