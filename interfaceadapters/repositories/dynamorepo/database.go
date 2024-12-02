package dynamorepo

import (
	"context"
	"fmt"
	"sync"
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
		AddImage(p *models.UserImage) error
		GetAllImagesPaginated(req models.PaginatedInput) (*models.UserImageResult, error)
		GetImage(uID, imgID string) (*models.UserImage, error)
		DeleteImage(uID, imgID string) error
		DeleteAllImages(uID string) error
		getAllItems(uID string) ([]models.UserImage, error)
		softDeleteItem(p *models.UserImage) error
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

func (d *DynamoDBRepo) AddImage(req *models.UserImage) error {
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

func (d *DynamoDBRepo) GetImage(uID, imgID string) (*models.UserImage, error) {
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

func (d *DynamoDBRepo) DeleteImage(uID, imageID string) error {
	imageResult, err := d.GetImage(uID, imageID)
	if err != nil {
		return err
	}

	err = d.softDeleteItem(imageResult)
	if err != nil {
		d.Log.Errorf("error deleting image %s of user %s in DB. error :: %v", imageID, uID, err)
		return err
	}

	return nil
}

func (d *DynamoDBRepo) GetAllImagesPaginated(req models.PaginatedInput) (*models.UserImageResult, error) {
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

	var imageResults []models.UserImage
	err = attributevalue.UnmarshalListOfMaps(result.Items, &imageResults)
	if err != nil {
		d.Log.Error("error unmarshaling db response", err)
		return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error unmarshaling db response"))
	}
	response := &models.UserImageResult{
		UserImages: imageResults,
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

func (d *DynamoDBRepo) getAllItems(uID string) ([]models.UserImage, error) {
	var lastEvaluatedKey map[string]types.AttributeValue
	var allImages []models.UserImage

	for {
		var imageResults []models.UserImage
		input := &dynamodb.QueryInput{
			TableName:              &d.TableName,
			IndexName:              aws.String("UserIDTakenAtIndex"),
			KeyConditionExpression: aws.String("UserID = :uID"),
			FilterExpression:       aws.String("IsDeleted = :isDeleted"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":uID":       &types.AttributeValueMemberS{Value: uID},
				":isDeleted": &types.AttributeValueMemberBOOL{Value: false},
			},
			ScanIndexForward:  aws.Bool(false),
			ExclusiveStartKey: lastEvaluatedKey,
		}

		result, err := d.Client.Query(context.Background(), input)
		if err != nil {
			d.Log.Error("error querying db", err)
			return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error querying db"))
		}

		err = attributevalue.UnmarshalListOfMaps(result.Items, &imageResults)
		if err != nil {
			d.Log.Error("error unmarshaling db response", err)
			return nil, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error unmarshaling db response"))
		}
		allImages = append(allImages, imageResults...)

		if result.LastEvaluatedKey == nil {
			break
		}

		lastEvaluatedKey = result.LastEvaluatedKey
	}

	return allImages, nil
}

func (d *DynamoDBRepo) softDeleteItem(req *models.UserImage) error {

	_, err := d.Client.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: &d.TableName,
		Key: map[string]types.AttributeValue{
			HashKey:  &types.AttributeValueMemberS{Value: req.UserID},
			RangeKey: &types.AttributeValueMemberS{Value: req.ImageID},
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":isDeleted": &types.AttributeValueMemberBOOL{Value: true},
			":updatedAt": &types.AttributeValueMemberS{Value: time.Now().String()},
		},
		UpdateExpression: aws.String("SET IsDeleted = :isDeleted, UpdatedAt = :updatedAt"),
	})
	if err != nil {
		d.Log.Errorf("error persisting scan %s of user %s. error %v", req.ImageID, req.UserID, err)
		return errors.New(errors.ErrCodeGeneric, fmt.Errorf("error processing image"))
	}

	return nil
}

func (d *DynamoDBRepo) DeleteAllImages(uID string) error {
	imageResults, err := d.getAllItems(uID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(imageResults))

	for _, item := range imageResults {
		wg.Add(1)
		go func(item models.UserImage) {
			defer wg.Done()
			err := d.softDeleteItem(&item)
			if err != nil {
				errChan <- err
			}
		}(item)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			d.Log.Errorf("error deleting images of user %s in DB. error :: %v", uID, err)
			return err
		}
	}

	return nil
}
