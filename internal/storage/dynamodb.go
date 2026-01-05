package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws-learning-game/internal/game"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDB Table Schema:
// Table: aws-learning-game-sessions
// - PK: GAME#{gameID}
// - SK: STATE
// - Attributes: GameState JSON

const (
	// DefaultTableName 預設 DynamoDB 表格名稱
	DefaultTableName = "aws-learning-game-sessions"
	// PartitionKeyPrefix 分區鍵前綴
	PartitionKeyPrefix = "GAME#"
	// SortKeyState 排序鍵 - 遊戲狀態
	SortKeyState = "STATE"
)

// DynamoDBStorage DynamoDB 儲存實作
type DynamoDBStorage struct {
	client    *dynamodb.Client
	tableName string
}

// DynamoDBConfig DynamoDB 配置
type DynamoDBConfig struct {
	TableName string
	Region    string
	Endpoint  string // 用於本地測試 (LocalStack/DynamoDB Local)
}

// GameItem DynamoDB 項目結構
type GameItem struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	GameID    string `dynamodbav:"GameID"`
	Status    string `dynamodbav:"Status"`
	GameData  string `dynamodbav:"GameData"` // JSON 序列化的遊戲狀態
	CreatedAt string `dynamodbav:"CreatedAt"`
	UpdatedAt string `dynamodbav:"UpdatedAt"`
}

// NewDynamoDBStorage 建立新的 DynamoDB 儲存
func NewDynamoDBStorage(ctx context.Context, cfg DynamoDBConfig) (*DynamoDBStorage, error) {
	// 載入 AWS 配置
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// 建立 DynamoDB 客戶端
	var client *dynamodb.Client
	if cfg.Endpoint != "" {
		// 使用自訂端點 (本地測試)
		client = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		})
	} else {
		client = dynamodb.NewFromConfig(awsCfg)
	}

	tableName := cfg.TableName
	if tableName == "" {
		tableName = DefaultTableName
	}

	return &DynamoDBStorage{
		client:    client,
		tableName: tableName,
	}, nil
}

// NewDynamoDBStorageWithClient 使用現有客戶端建立儲存 (用於測試)
func NewDynamoDBStorageWithClient(client *dynamodb.Client, tableName string) *DynamoDBStorage {
	if tableName == "" {
		tableName = DefaultTableName
	}
	return &DynamoDBStorage{
		client:    client,
		tableName: tableName,
	}
}

// SaveGame 儲存遊戲狀態到 DynamoDB
// Requirements 8.1: 儲存遊戲狀態
// Requirements 8.3: 使用 DynamoDB 儲存
func (s *DynamoDBStorage) SaveGame(ctx context.Context, g *game.Game) error {
	if g == nil {
		return game.ErrGameNotFound
	}

	// 序列化遊戲狀態
	gameData, err := game.SerializeGameState(g)
	if err != nil {
		return fmt.Errorf("%w: %v", game.ErrSaveFailed, err)
	}

	// 建立 DynamoDB 項目
	item := GameItem{
		PK:        PartitionKeyPrefix + g.ID,
		SK:        SortKeyState,
		GameID:    g.ID,
		Status:    string(g.Status),
		GameData:  string(gameData),
		CreatedAt: g.CreatedAt.Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// 轉換為 DynamoDB 屬性
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("%w: failed to marshal item: %v", game.ErrSaveFailed, err)
	}

	// 寫入 DynamoDB
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", game.ErrSaveFailed, err)
	}

	return nil
}

// LoadGame 從 DynamoDB 載入遊戲狀態
// Requirements 8.2: 載入遊戲狀態
// Requirements 8.3: 使用 DynamoDB 儲存
func (s *DynamoDBStorage) LoadGame(ctx context.Context, gameID string) (*game.Game, error) {
	if gameID == "" {
		return nil, game.ErrGameNotFound
	}

	// 建立查詢鍵
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: PartitionKeyPrefix + gameID},
		"SK": &types.AttributeValueMemberS{Value: SortKeyState},
	}

	// 從 DynamoDB 讀取
	result, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", game.ErrLoadFailed, err)
	}

	// 檢查項目是否存在
	if result.Item == nil {
		return nil, game.ErrGameNotFound
	}

	// 解析 DynamoDB 項目
	var item GameItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal item: %v", game.ErrLoadFailed, err)
	}

	// 反序列化遊戲狀態
	g, err := game.DeserializeGameState([]byte(item.GameData))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", game.ErrCorruptedState, err)
	}

	return g, nil
}

// DeleteGame 從 DynamoDB 刪除遊戲狀態
func (s *DynamoDBStorage) DeleteGame(ctx context.Context, gameID string) error {
	if gameID == "" {
		return game.ErrGameNotFound
	}

	// 建立刪除鍵
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: PartitionKeyPrefix + gameID},
		"SK": &types.AttributeValueMemberS{Value: SortKeyState},
	}

	// 從 DynamoDB 刪除
	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.tableName),
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}

	return nil
}

// ListGames 列出所有遊戲
func (s *DynamoDBStorage) ListGames(ctx context.Context) ([]*GameMetadata, error) {
	// 掃描所有遊戲項目
	result, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(s.tableName),
		FilterExpression: aws.String("SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": &types.AttributeValueMemberS{Value: SortKeyState},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	// 解析結果
	games := make([]*GameMetadata, 0, len(result.Items))
	for _, item := range result.Items {
		var gameItem GameItem
		if err := attributevalue.UnmarshalMap(item, &gameItem); err != nil {
			continue // 跳過無法解析的項目
		}

		// 反序列化以獲取玩家數量
		g, err := game.DeserializeGameState([]byte(gameItem.GameData))
		if err != nil {
			continue
		}

		games = append(games, &GameMetadata{
			GameID:      gameItem.GameID,
			Status:      gameItem.Status,
			PlayerCount: len(g.Players),
			CurrentTurn: g.CurrentTurn,
			CreatedAt:   gameItem.CreatedAt,
			UpdatedAt:   gameItem.UpdatedAt,
		})
	}

	return games, nil
}

// CreateTable 建立 DynamoDB 表格 (用於初始化)
func (s *DynamoDBStorage) CreateTable(ctx context.Context) error {
	_, err := s.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(s.tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// TableExists 檢查表格是否存在
func (s *DynamoDBStorage) TableExists(ctx context.Context) (bool, error) {
	_, err := s.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(s.tableName),
	})
	if err != nil {
		// 檢查是否為表格不存在的錯誤
		if isResourceNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// isResourceNotFoundError 檢查是否為資源不存在錯誤
func isResourceNotFoundError(err error) bool {
	// 簡單的錯誤類型檢查
	_, ok := err.(*types.ResourceNotFoundException)
	return ok
}
