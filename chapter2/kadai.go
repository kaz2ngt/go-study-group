package chapter2

import "fmt"

// 引数のスライスsliceの要素数が
// 0の場合、0とエラー
// 2以下の場合、要素を掛け算
// 3以上の場合、要素を足し算
// を返却。正常終了時、errorはnilでよい
func Calc(slice []int) (int, error) {
	// TODO Q1
	// ヒント：エラーにも色々な生成方法があるが、ここではシンプルにfmtパッケージの
	// fmt.Errorf(“invalid op=%s”, op) などでエラー内容を返却するのがよい
	// https://golang.org/pkg/fmt/#Errorf

	if length := len(slice); length == 0 {
		return 0, fmt.Errorf("invalid len(slice)=%d", length)
	} else if length <= 2 {
		// 積和(初期値1)
		sumOfProducts := 1
		for _, value := range slice {
			sumOfProducts *= value
		}
		return sumOfProducts, nil
	}

	// 和(初期値0)
	sum := 0
	for _, value := range slice {
		sum += value
	}

	return sum, nil
}

type Number struct {
	index int
}

// 構造体Numberを3つの要素数から成るスライスにして返却
// 3つの要素の中身は[{1} {2} {3}]とし、append関数を使用すること
func Numbers() []Number {
	// TODO Q2

	result := []Number{{1}, {2}, {3}}

	return result
}

// 引数mをforで回し、「値」部分だけの和を返却
// キーに「yon」が含まれる場合は、キー「yon」に関連する値は除外すること
// キー「yon」に関しては完全一致すること
func CalcMap(m map[string]int) int {
	// TODO Q3

	// for文内で毎回ifチェックするよりは早いはず
	sum := -m["yon"]
	for _, value := range m {
		sum += value
	}

	return sum
}

type Model struct {
	Value int
}

// 与えられたスライスのModel全てのValueに5を足す破壊的な関数を作成
func Add(models []Model) {
	// TODO  Q4

	for index := range models {
		models[index].Value += 5
	}
}

// 引数のスライスには重複な値が格納されているのでユニークな値のスライスに加工して返却
// 順序はスライスに格納されている順番のまま返却すること
// ex) 引数:[]slice{21,21,4,5} 戻り値:[]int{21,4,5}
func Unique(slice []int) []int {
	// TODO Q5

	// capを余裕をもって確保
	uniqueSlice := make([]int, 0, len(slice))
	uniqueMap := make(map[int]bool)
	for _, value := range slice {
		if uniqueMap[value] {
			continue
		}
		uniqueMap[value] = true
		uniqueSlice = append(uniqueSlice, value)
	}
	// cap削減
	uniqueSlice = append([]int{}, uniqueSlice...)

	return uniqueSlice
}

// 連続するフィボナッチ数(0, 1, 1, 2, 3, 5, ...)を返す関数(クロージャ)を返却
func Fibonacci() func() int {
	// TODO Q6 オプション

	current := 0
	next := 1
	return func() int {
		result := current
		current = next
		next = next + result
		return result
	}
}
