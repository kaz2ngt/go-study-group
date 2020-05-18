package chapter6

//go:generate mockgen -source=$GOFILE -destination=kadai_mock.go -package=$GOPACKAGE -self_package=github.com/apbgo/go-study-group/$GOPACKAGE

import (
	"context"
	"database/sql"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// [課題内容]
// 以下の2つのInterface(IFUserItemService, IFUserItemRepository)を満たすstructを実装してください。
// 今回の課題ではトランザクション境界はProvide()内で構いません。
// gomockを利用してRepositoryのMockファイルを自動生成しています。
// テストではIFUserItemRepositoryをモックしたUserItemServiceを使ってみましょう。
// TransactionのBegin, Commit, Rollbackも本来Mockしたいですが、実装が多くなってしまうので
// 今回はする必要がありません。

// Reward 報酬モデル
type Reward struct {
	ItemID int64
	Count  int64
}

// IFUserItemService 報酬の付与の機能を表すインターフェイス
type IFUserItemService interface {
	// 対象のUserIDに引数で渡された報酬を付与します.
	Provide(ctx context.Context, userID int64, rewards ...Reward) error
}

// IFUserItemRepository i_user_itemテーブルへの操作を行うインターフェイス
type IFUserItemRepository interface {
	// FindByUserIdAndItemIDs 一致するモデルを複数返却する.
	FindByUserIdAndItemIDs(
		ctx context.Context,
		tx *sql.Tx,
		userID int64,
		itemIDs []int64,
	) (iUserItems []*IUserItem, err error)

	// Insert 対象のモデルから1件Insertを実行する
	Insert(
		ctx context.Context,
		tx *sql.Tx,
		iUserItem *IUserItem,
	) error

	// Update対象のモデルから1件Updateを実行する
	// Update対象レコードが0件の場合、okはfalseになる
	Update(
		ctx context.Context,
		tx *sql.Tx,
		iUserItem *IUserItem,
	) (ok bool, err error)
}

// UserItemService [実装対象]
type UserItemService struct {
	db                 *sql.DB
	userItemRepository IFUserItemRepository
}

func (s *UserItemService) Provide(ctx context.Context, userID int64, rewards ...Reward) error {
	// Transactionを開始
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}
	// db変わらないようにCommitではなくRollbackにしておく
	defer tx.Rollback()

	// itemID毎にrewardの個数をまとめたmapを形成
	rewardMap := make(map[int64]*Reward, len(rewards))
	for _, reward := range rewards {
		if _, ok := rewardMap[reward.ItemID]; !ok {
			rewardMap[reward.ItemID] = &Reward{ItemID: reward.ItemID}
		}
		rewardMap[reward.ItemID].Count += reward.Count
	}

	// itemIDsの抽出
	itemIDs := make([]int64, 0, len(rewardMap))
	for itemID, _ := range rewardMap {
		itemIDs = append(itemIDs, itemID)
	}

	// userItemsの取得とusertItemMapへの変換
	userItems, err := s.userItemRepository.FindByUserIdAndItemIDs(ctx, tx, userID, itemIDs)
	userItemMap := make(map[int64]*IUserItem, len(userItems))
	for _, userItem := range userItems {
		userItemMap[userItem.ItemID] = userItem
	}

	// rewardMapからrewardを付与
	now := time.Now()
	for itemID, reward := range rewardMap {
		if userItem, ok := userItemMap[itemID]; !ok {
			// userItemを所持していない場合insert
			insertUserItem := &IUserItem{
				UserID:    userID,
				ItemID:    itemID,
				Count:     reward.Count,
				CreatedAt: now,
				UpdatedAt: now,
			}
			err := s.userItemRepository.Insert(ctx, tx, insertUserItem)
			if err != nil {
				return err
			}
		} else {
			// userItemを所持している場合update
			updateUserItem := &IUserItem{
				UserID:    userID,
				ItemID:    itemID,
				Count:     userItem.Count + reward.Count,
				CreatedAt: userItem.CreatedAt,
				UpdatedAt: now,
			}
			_, err := s.userItemRepository.Update(ctx, tx, updateUserItem)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// NewUserItemService コンストラクタ [実装対象]
func NewUserItemService(
	db *sql.DB,
	userItemRepository IFUserItemRepository,
) *UserItemService {
	return &UserItemService{
		db:                 db,
		userItemRepository: userItemRepository,
	}
}

// UserItemRepository [実装対象]
type UserItemRepository struct {
}

func (r *UserItemRepository) FindByUserIdAndItemIDs(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	itemIDs []int64,
) (iUserItems []*IUserItem, err error) {
	if len(itemIDs) == 0 {
		return nil, nil
	}
	query := "SELECT * FROM i_user_item WHERE user_id = ?, item_id IN (?" + strings.Repeat(", ?", len(itemIDs)-1) + ")"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	txStmt := tx.Stmt(stmt)
	defer txStmt.Close()

	args := make([]interface{}, 0, len(itemIDs)+1)
	args = append(args, userID)
	for _, itemID := range itemIDs {
		args = append(args, itemID)
	}
	rows, err := txStmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 1レコードずつ、処理する
	for rows.Next() {
		// モデルを作成して、カラムへのポインタをScan()に渡す
		iUserItem := IUserItem{}
		if err := rows.Scan(
			&iUserItem.UserID,
			&iUserItem.ItemID,
			&iUserItem.Count,
			&iUserItem.CreatedAt,
			&iUserItem.UpdatedAt,
			&iUserItem.DeletedAt,
		); err != nil {
			return nil, err
		}
		iUserItems = append(iUserItems, &iUserItem)
	}

	return iUserItems, nil
}

func (r *UserItemRepository) Insert(
	ctx context.Context,
	tx *sql.Tx,
	iUserItem *IUserItem,
) error {
	query := "INSERT INTO i_user_item VALUES (?, ?, ?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	txStmt := tx.Stmt(stmt)
	defer txStmt.Close()

	_, err = txStmt.ExecContext(
		ctx,
		iUserItem.UserID,
		iUserItem.ItemID,
		iUserItem.Count,
		iUserItem.CreatedAt,
		iUserItem.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserItemRepository) Update(
	ctx context.Context,
	tx *sql.Tx,
	iUserItem *IUserItem,
) (ok bool, err error) {
	query := "UPDATE i_user_item set Count = ?, UpdatedAt = ? WHERE user_id = ?, item_id = ?"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return false, err
	}
	txStmt := tx.Stmt(stmt)
	defer txStmt.Close()

	_, err = txStmt.ExecContext(
		ctx,
		iUserItem.Count,
		iUserItem.UpdatedAt,
		iUserItem.UserID,
		iUserItem.ItemID,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// NewUserItemRepository コンストラクタ [実装対象]
func NewUserItemRepository() *UserItemRepository {
	return &UserItemRepository{}
}
