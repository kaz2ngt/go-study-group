package chapter6

import (
	"context"
	"database/sql"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// このテストはgomockを利用するサンプルです。
func TestSample(t *testing.T) {
	t.Run("サンプル1", func(t *testing.T) {
		// サブテストごとにMockのControllerを作成してください。
		ctrl := gomock.NewController(t)

		// 自動生成されたMockをNewする
		mock := NewMockIFUserItemRepository(ctrl)

		// ここからは意味がないテスト
		ctx := context.Background()
		userItem := IUserItem{
			UserID: 1,
			ItemID: 1,
			Count:  100,
		}

		mock.EXPECT().
			// ここに渡された変数は、値が一致しない場合はテストが成功しない（ポインタの場合はポインタの一致）
			// gomock.Any()は全ての値が許容される
			// *sql.TxやiUserItemはService内で生成されるため、ポインタの一致をチェックすることは難しい
			Update(ctx, gomock.Any(), gomock.Any()).
			// そのためDoAndReturnの関数内で値をチェックしてあげるとよい
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userItem1 *IUserItem) (ok bool, err error) {
					assert.Equal(t, userItem, *userItem1)
					// この関数の戻り値がMockを実行した時の戻り値になる
					return true, nil
				},
			)

		ok, err := mock.Update(ctx, nil, &userItem)
		assert.True(t, ok)
		assert.NoError(t, err)
	})
}

func TestUserItemService_Provide(t *testing.T) {
	t.Run("正常:正常に付与できる", func(t *testing.T) {
		// サブテストごとにMockのControllerを作成してください。
		ctrl := gomock.NewController(t)

		// 自動生成されたMockをNewする
		mock := NewMockIFUserItemRepository(ctrl)

		// ここからは意味がないテスト
		ctx := context.Background()
		var userID int64 = 1
		rewards := []Reward{
			{ItemID: 1, Count: 1},
			{ItemID: 2, Count: 1},
			{ItemID: 1, Count: 10},
		}

		mock.EXPECT().FindByUserIdAndItemIDs(ctx, gomock.Any(), userID, []int64{1, 2}).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userID1 int64, itemIDs1 []int64) (iUserItems []*IUserItem, err error) {
					iUserItems = []*IUserItem{
						{
							UserID:    1,
							ItemID:    1,
							Count:     1,
							CreatedAt: time.Unix(1, 0),
							UpdatedAt: time.Unix(1, 0),
						},
					}
					return iUserItems, nil
				},
			)
		mock.EXPECT().Insert(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, iUserItem1 *IUserItem) (ok bool, err error) {
					// itemID=2が正常にインサートされる
					assert.Equal(t, userID, iUserItem1.UserID)
					assert.Equal(t, int64(2), iUserItem1.ItemID)
					assert.Equal(t, int64(1), iUserItem1.Count)
					return true, nil
				},
			)

		mock.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, iUserItem1 *IUserItem) (ok bool, err error) {
					// itemID=1が正常にアップデートされる
					assert.Equal(t, userID, iUserItem1.UserID)
					assert.Equal(t, int64(1), iUserItem1.ItemID)
					assert.Equal(t, int64(12), iUserItem1.Count)
					return true, nil
				},
			)

		// 接続先設定
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:5446)/chapter6?parseTime=true")
		if err != nil {
			t.Fail()
		}
		// アプリケーションが終了するときにCloseするように
		defer db.Close()

		// こーゆーことかと思ったら違った
		userItemService := NewUserItemService(db, mock)
		userItemService.Provide(ctx, userID, rewards...)
	})
}
