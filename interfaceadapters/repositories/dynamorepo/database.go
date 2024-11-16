package dynamorepo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rahul-aut-ind/service-user/domain/errors"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/internal/awsconfig"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
)

type (
	DataHandler interface {
		CreateOrUpdate(p *models.UserImage) error
		GetAllByUserIDPaginated(req models.PaginatedInput) (*models.UserImageResult, error)
		getAllByUserID(uID string) ([]models.UserImage, error)
		GetByImgID(uID, imgID string) (*models.UserImage, error)
		DeleteByImgID(uID, imgID string) error
		DeleteAllImagesByUserID(uID string) error
	}

	DynamoDBRepo struct {
		TableName string
		Client    *dynamodb.Client
		Log       *logger.Logger
	}
)

const (
	HashKey              = "UserID"
	RangeKey             = "ImageID"
	IndexRangeKey        = "TakenAt"
	GlobalSecondaryIndex = "UserIDTakenAtIndex"
)

func New(cfg *awsconfig.AWSConfig, env *config.Env, log *logger.Logger) *DynamoDBRepo {
	return &DynamoDBRepo{TableName: env.DynamoDBTable, Client: createClient(cfg.Config), Log: log}
}

func createClient(cfg *aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(*cfg)
}

// CreateOrUpdate creates or updates an image in the db
func (d *DynamoDBRepo) CreateOrUpdate(req *models.UserImage) error {
	item, err := attributevalue.MarshalMap(req)
	if err != nil {
		d.Log.Error("error marshaling input", err)
		return errors.New(errors.ErrCodeGeneric, fmt.Errorf("error marshaling input"))
	}

	_, err = d.Client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: &d.TableName,
		Item:      item,
	})
	if err != nil {
		d.Log.Errorf("error persisting image %s of user %s to db %v", req.ImageID, req.UserID, err)
		return errors.New(errors.ErrCodeGeneric, fmt.Errorf("error persisting image data"))
	}

	return nil
}

// GetByScanID returns an image by its ID and user ID
func (d *DynamoDBRepo) GetByImgID(uID, imgID string) (*models.UserImage, error) {
	input := &dynamodb.GetItemInput{
		TableName: &d.TableName,
		Key: map[string]types.AttributeValue{
			HashKey:  &types.AttributeValueMemberS{Value: uID},
			RangeKey: &types.AttributeValueMemberS{Value: imgID},
		},
	}

	result, err := d.Client.GetItem(context.Background(), input)
	if err != nil {
		d.Log.Error("error querying db", err)
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error querying db"))
	}
	if result.Item == nil {
		return nil, errors.New(errors.ErrCodeNotFound, fmt.Errorf("image not found"))
	}

	if result.Item["IsDeleted"].(*types.AttributeValueMemberBOOL).Value {
		return nil, errors.New(errors.ErrCodeNotFound, fmt.Errorf("image not found"))
	}

	var imageResult models.UserImage
	err = attributevalue.UnmarshalMap(result.Item, &imageResult)
	if err != nil {
		d.Log.Error("error unmarshaling db response", err)
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error unmarshaling db response"))
	}

	return &imageResult, nil
}

// DeleteByScanID soft deletes an image
func (d *DynamoDBRepo) DeleteByImgID(uID, imageID string) error {
	imageResult, err := d.GetByImgID(uID, imageID)
	if err != nil {
		return err
	}

	imageResult.IsDeleted = true
	imageResult.UpdatedAt = time.Now()
	err = d.CreateOrUpdate(imageResult)
	if err != nil {
		d.Log.Errorf("error deleting image %s of user %s in DB. error :: %v", imageID, uID, err)
		return err
	}

	return nil
}

// GetAllByUserIDPaginated returns all non-deleted scans for a user in paginated manner
func (d *DynamoDBRepo) GetAllByUserIDPaginated(req models.PaginatedInput) (*models.UserImageResult, error) {
	input := &dynamodb.QueryInput{
		TableName:              &d.TableName,
		IndexName:              aws.String(GlobalSecondaryIndex),
		KeyConditionExpression: aws.String("UserID = :uID"),
		FilterExpression:       aws.String("IsDeleted = :isDeleted"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uID":       &types.AttributeValueMemberS{Value: req.UserID},
			":isDeleted": &types.AttributeValueMemberBOOL{Value: false},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(req.Limit),
	}

	if req.LastImageID != "" && req.LastImageTakenAt != "" {
		input.ExclusiveStartKey = map[string]types.AttributeValue{
			HashKey:       &types.AttributeValueMemberS{Value: req.UserID},
			RangeKey:      &types.AttributeValueMemberS{Value: req.LastImageID},
			IndexRangeKey: &types.AttributeValueMemberS{Value: req.LastImageTakenAt},
		}
	}

	result, err := d.Client.Query(context.Background(), input)
	if err != nil {
		d.Log.Error("error querying db", err)
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error querying db"))
	}

	var scanResults []models.UserImage
	err = attributevalue.UnmarshalListOfMaps(result.Items, &scanResults)
	if err != nil {
		d.Log.Error("error unmarshaling db response", err)
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error unmarshaling db response"))
	}
	response := &models.UserImageResult{
		UserImages: scanResults,
	}

	if result.LastEvaluatedKey != nil {
		newMap := make(map[string]string)
		for k, v := range result.LastEvaluatedKey {
			if attrS, ok := v.(*types.AttributeValueMemberS); ok {
				switch k {
				case RangeKey:
					newMap[config.QueryParamLastKey] = attrS.Value
				case IndexRangeKey:
					newMap[config.QueryParamlastKeyDate] = attrS.Value
				}
			}
			response.Page.LastEvaluatedKey = newMap
		}
	}

	return response, nil
}

// GetAllByUserID returns all non-deleted images for a user
func (d *DynamoDBRepo) getAllByUserID(uID string) ([]models.UserImage, error) {
	input := &dynamodb.QueryInput{
		TableName:              &d.TableName,
		IndexName:              aws.String("UserIDTakenAtIndex"),
		KeyConditionExpression: aws.String("UserID = :uID"),
		FilterExpression:       aws.String("IsDeleted = :isDeleted"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uID":       &types.AttributeValueMemberS{Value: uID},
			":isDeleted": &types.AttributeValueMemberBOOL{Value: false},
		},
		ScanIndexForward: aws.Bool(false),
	}

	result, err := d.Client.Query(context.Background(), input)
	if err != nil {
		d.Log.Error("error querying db", err)
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error querying db"))
	}

	var imageResults []models.UserImage
	err = attributevalue.UnmarshalListOfMaps(result.Items, &imageResults)
	if err != nil {
		d.Log.Error("error unmarshaling db response", err)
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error unmarshaling db response"))
	}

	return imageResults, nil
}

// DeleteAllScansByUserID soft deletes all images for a user
func (d *DynamoDBRepo) DeleteAllImagesByUserID(uID string) error {
	imageResults, err := d.getAllByUserID(uID)
	if err != nil {
		return err
	}

	for i := range imageResults {
		imageResult := &imageResults[i]
		imageResult.IsDeleted = true
		imageResult.UpdatedAt = time.Now()
		err := d.CreateOrUpdate(imageResult)
		if err != nil {
			d.Log.Errorf("error deleting images of user %s in DB. error :: %v", uID, err)
			return err
		}
	}

	return nil
}
