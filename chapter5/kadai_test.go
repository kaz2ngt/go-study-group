package chapter5

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	expectedFlagOption := FlagOption{
		fields:    1,
		delimiter: ",",
	}
	t.Run("正常：フラグで取得した情報が正常", func(t *testing.T) {
		assert.Equal(t, nil, Validation(1, expectedFlagOption))
	})
	t.Run("異常：ファイルが指定されていない", func(t *testing.T) {
		assert.Equal(t, fmt.Errorf("ファイルパスを指定してください。"), Validation(0, expectedFlagOption))
	})
	t.Run("異常：-f は1以上である必要があります", func(t *testing.T) {
		invalidFlagOption := FlagOption{
			fields:    0,
			delimiter: ",",
		}
		assert.Equal(t, fmt.Errorf("-f は1以上である必要があります"), Validation(1, invalidFlagOption))
	})
}

func TestCut(t *testing.T) {
	t.Run("正常：正常に1列目をカットできる", func(t *testing.T) {
		flagOption := FlagOption{
			fields:    1,
			delimiter: ",",
		}
		reader := bytes.NewBufferString("a,b,c\nd,e,f")
		writer := new(bytes.Buffer)
		assert.Equal(t, nil, Cut(reader, writer, flagOption))
		assert.Equal(t, "a\nd\n", writer.String())
	})
	t.Run("正常：正常に指定した列をカットできる", func(t *testing.T) {
		flagOption := FlagOption{
			fields:    2,
			delimiter: ",",
		}
		reader := bytes.NewBufferString("a,b,c\nd,e,f")
		writer := new(bytes.Buffer)
		assert.Equal(t, nil, Cut(reader, writer, flagOption))
		assert.Equal(t, "b\ne\n", writer.String())
	})
	t.Run("正常：正常に指定したdelimiterでカットできる", func(t *testing.T) {
		flagOption := FlagOption{
			fields:    1,
			delimiter: " ",
		}
		reader := bytes.NewBufferString("a b c\nd e f")
		writer := new(bytes.Buffer)
		assert.Equal(t, nil, Cut(reader, writer, flagOption))
		assert.Equal(t, "a\nd\n", writer.String())
	})
	t.Run("異常：指定した列数が実際より多い場合エラー", func(t *testing.T) {
		flagOption := FlagOption{
			fields:    4,
			delimiter: ",",
		}
		reader := bytes.NewBufferString("a,b,c\nd,e,f")
		writer := new(bytes.Buffer)
		assert.Equal(t, fmt.Errorf("-fの値に該当するデータがありません"), Cut(reader, writer, flagOption))
	})
}
