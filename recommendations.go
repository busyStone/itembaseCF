package main

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

type similarFuncType func(m dataSetType, item1, item2 string) float32

// 皮尔逊相似度计算
// 使用构造的数据集计算 item1 item2 的相似度
// 返回 相似度
func pearson(m dataSetType, item1, item2 string) float32 {

	var sum1, sum2, sum1Sq, sum2Sq, pSum float64
	var cnt uint32

	for usr1, v1 := range m[item1] {
		v2, ok := m[item2][usr1]
		if ok {
			sum1 += float64(v1)
			sum2 += float64(v2)

			sum1Sq += math.Pow(float64(v1), 2)
			sum2Sq += math.Pow(float64(v2), 2)

			pSum += float64(v1 * v2)

			cnt++
		}
	}

	if cnt == 0 {
		return 0
	}

	num := pSum - (sum1 * sum2 / float64(cnt))
	den := math.Sqrt((sum1Sq - math.Pow(sum1, 2)/float64(cnt)) * (sum2Sq - math.Pow(sum2, 2)/float64(cnt)))
	if den == 0 {
		return 0
	}

	return float32(num / den)
}

func distance(m dataSetType, item1, item2 string) float32 {

	var sum float64
	var similar bool

	for item, v1 := range m[item1] {
		v2, ok := m[item2][item]
		if ok {
			similar = true

			sum += math.Pow(float64(v1)-float64(v2), 2)
		}
	}

	if similar {
		return float32(1 / (1 + math.Sqrt(sum)))
	}

	return 0.0
}

func en_distance(m dataSetType, item1, item2 string) float32 {

	var sum float64
	var similar bool

	for item, v1 := range m[item1] {
		v2, ok := m[item2][item]
		if ok {
			similar = true

			sum += math.Pow(float64(v1)-float64(v2), 2)
		}
	}

	if similar {
		return float32(1 / (1 + sum))
	}

	return 0.0
}

// 使用构造的数据集计算与 refItem 最为匹配的item
// n 指定返回匹配个数
// similarFunc 指定相似度计算函数
func topMatches(m dataSetType, refItem string, n uint8, similarFunc similarFuncType) (RankSlice, error) {

	var rs RankSlice
	var r Rank

	// 判断参考元素在数据集中是否存在
	_, ok := m[refItem]
	if !ok {
		return nil, errors.New("There is no " + refItem)
	}

	for item, _ := range m {
		if item != refItem {
			r.name = item
			r.similar = similarFunc(m, refItem, item)

			rs = append(rs, r)
		}
	}

	sort.Sort(rs)
	sort.Sort(sort.Reverse(rs))

	if int(n) > len(rs) {
		n = uint8(len(rs))
	}

	//return rs[0:n], nil
	return rs, nil
}

// 计算数据集中每个项目 n 个相似 item
func calcSimilarItems(m dataSetType, n uint8, similarFunc similarFuncType) (map[string]RankSlice, error) {

	if len(m) == 0 {
		return nil, errors.New("data set is empty.")
	}

	simItems := make(map[string]RankSlice)

	cnt := 0
	for item, _ := range m {
		items, err := topMatches(m, item, n, similarFunc)
		if err != nil {
			return nil, err
		}

		simItems[item] = items

		cnt += 1

		if cnt%100 == 0 {
			fmt.Printf("%d / %d \n", cnt, len(m))
		}
	}

	return simItems, nil
}

//
// userItemsMap 指定用户打过分的项目
// itemMatchMap 通过 calcSimilarItems 计算出的相似项集合
func getRecommendedItems(userItemsMap map[string]float32, itemMatchMap map[string]RankSlice) RankSlice {

	totalSimMap := make(map[string]float32)
	rateSimMap := make(map[string]float32)

	// 依次计算用户打过分的项目 item1
	for item1, rating := range userItemsMap {
		s := itemMatchMap[item1]

		for _, booksimilar := range s {
			_, ok := userItemsMap[booksimilar.name] // 已经看过的不做计算
			if ok {
				continue
			}

			totalSimMap[booksimilar.name] += booksimilar.similar
			rateSimMap[booksimilar.name] += rating * booksimilar.similar
		}
	}

	var rs RankSlice
	var r Rank

	for k, v := range totalSimMap {
		r.name = k
		if v == 0 {
			r.similar = 0
		} else {
			r.similar = rateSimMap[k] / v
		}

		rs = append(rs, r)
	}

	sort.Sort(rs)
	sort.Sort(sort.Reverse(rs))

	return rs
}

func RecommendBook2User(userId string, limit uint) ([]Book, error) {

	mItem, mUser, _ := loaddata()

	rsMap, err := calcSimilarItems(mItem, 10, en_distance) // pearson)
	if err != nil {
		return nil, err
	}

	_, ok := mUser[userId]
	if !ok {
		return nil, errors.New("userId is wrong!")
	}

	rs := getRecommendedItems(mUser[userId], rsMap)

	if limit > uint(len(rs)) {
		limit = uint(len(rs))
	}

	rs = rs[0:limit]

	var books []Book
	var sum float32
	var book Book
	for _, v := range rs {
		book.name = v.name

		sum = 0
		for _, score := range mItem[v.name] {
			sum += score
		}

		book.score = sum / float32(len(mItem[v.name]))

		books = append(books, book)
	}

	sort.Sort(BookSlice(books))
	sort.Sort(sort.Reverse(BookSlice(books)))

	return books, nil
}
