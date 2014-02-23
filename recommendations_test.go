package main

import (
	"fmt"
	"testing"
)

func TestSimilarDistance(t *testing.T) {
	_, mUser, _ := loadExampleData()

	if sim := distance(mUser, "Lisa Rose", "Gene Seymour"); sim != 0.29429805 {
		t.Error("distance(Lisa Rose, Gene Seymour) calc error! sim = " + fmt.Sprintf("%f", sim))
	} else {
		t.Log("distance(Lisa Rose, Gene Seymour) pass.")
	}

	if sim := distance(mUser, "Lisa Rose", "Toby"); sim != 0.34833147 {
		t.Error("distance(Lisa Rose, Toby) calc error! sim = " + fmt.Sprintf("%f", sim))
	} else {
		t.Log("distance(Lisa Rose, Toby) pass.")
	}

	if sim := distance(mUser, "Lisa Rose", "Lisa Rose"); sim != 1 {
		t.Error("distance(Lisa Rose, Lisa Rose) calc error! sim = " + fmt.Sprintf("%f", sim))
	} else {
		t.Log("distance(Lisa Rose, Lisa Rose) pass.")
	}

	if sim := distance(mUser, "Lisa Rose", ""); sim != 0 {
		t.Error("distance(Lisa Rose, ) calc error sim = " + fmt.Sprintf("%f", sim))
	} else {
		t.Log("distance(Lisa Rose, ) pass.")
	}

	if sim := distance(mUser, "", ""); sim != 0 {
		t.Error("distance() calc error sim = " + fmt.Sprintf("%f", sim))
	} else {
		t.Log("distance() pass.")
	}
}

func TestSimilarPearson(t *testing.T) {
	_, mUser, _ := loadExampleData()
	if sim := pearson(mUser, "Lisa Rose", "Gene Seymour"); sim != 0.396059 {
		t.Error("pearson(Lisa Rose, Gene Seymour) calc error sim = ")
	} else {
		t.Log("pearson(Lisa Rose, Gene Seymour) pass.")
	}

	if sim := pearson(mUser, "Lisa Rose", "Lisa Rose"); sim != 1 {
		t.Error("pearson(Lisa Rose, Lisa Rose) calc error sim = " + fmt.Sprintf("%f", sim))
	} else {
		t.Log("pearson(Lisa Rose, Lisa Rose) pass.")
	}

	if sim := pearson(mUser, "Lisa Rose", ""); sim != 0 {
		t.Error("pearson(Lisa Rose, ) calc error sim = " + fmt.Sprintf("%f", sim))
	} else {
		t.Log("pearson(Lisa Rose, ) pass.")
	}

	if sim := pearson(mUser, "", ""); sim != 0 {
		t.Error("pearson() calc error sim = " + fmt.Sprintf("%f", sim))
	} else {
		t.Log("pearson() pass.")
	}
}

func TestTopMatches(t *testing.T) {

	var testFail bool

	resultUser := []Rank{
		{"Lisa Rose", 0.9912407},
		{"Mick LaSalle", 0.92447346},
		{"Claudia Puig", 0.89340514},
	}

	resultItem := []Rank{
		{"You, Me and Dupree", 0.6579517},
		{"Lady in the Water", 0.48795003},
		{"Snakes on a Plane", 0.1118034},
		{"The Night Listener", -0.1798472},
		{"Just My Luck", -0.42289004},
	}

	mItem, mUser, _ := loadExampleData()

	s, _ := topMatches(mUser, "Toby", 3, pearson)
	for i := 0; i < len(resultUser); i++ {
		if s[i] != resultUser[i] {
			testFail = true

			goto TestTopMatchesExit
		}
	}

	s, _ = topMatches(mItem, "Superman Returns", 5, pearson)
	for i := 0; i < len(resultItem); i++ {
		if s[i] != resultItem[i] {
			testFail = true

			goto TestTopMatchesExit
		}
	}

	s, _ = topMatches(mItem, "Superman Returns", 0, pearson)
	if len(s) != 0 {
		testFail = true
	}

	s, _ = topMatches(mItem, "", 1, pearson)
	if s != nil {
		testFail = true
	}

TestTopMatchesExit:

	if testFail {
		t.Error("topMatches use pearson test failed!")
	} else {
		t.Log("topMatches use pearson test pass.")
	}

}

func TestCalcSimilarItems(t *testing.T) {
	// Lady in the Water
	resultLW := []Rank{
		{"You, Me and Dupree", 0.4},
		{"The Night Listener", 0.2857143},
	}

	mItem, _, _ := loadExampleData()

	// 计算结果与书中不太一样  顺序是一样的
	// 直接使用 topMatches 计算所得结果也与书中此示例不同
	// 中文版 distance 最后结果为 1/(1+sqrt(sum_of_squares))
	// 英文版 distance 最后结果为 1/(1+sum_of_squares)  源码也是这样
	// 英文版应该是错误的 中文版做了修订 但是 中文版在 calcSimilarItems 是仍然引用了 英文版的结果。。。
	s, _ := calcSimilarItems(mItem, 2, en_distance)

	if len(s["Lady in the Water"]) < 2 {
		t.Error("calcSimilarItems result num error, test failed!")
		return
	}

	for i := 0; i < 2; i++ {
		if s["Lady in the Water"][i].name != resultLW[i].name ||
			s["Lady in the Water"][i].similar != resultLW[i].similar {
			t.Error("calcSimilarItems calc error, test failed")
			return
		}
	}

	t.Log("calcSimilarItems use en_distance test pass.")
}

func TestGetRecommendedItems(t *testing.T) {

	resultRank := []Rank{
		{"The Night Listener", 3.1826348},
		{"Just My Luck", 2.5983317},
		{"Lady in the Water", 2.4730878},
	}

	mItem, mUser, _ := loadExampleData()

	rsMap, err := calcSimilarItems(mItem, 10, en_distance)
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, ok := mUser["Toby"]
	if !ok {
		t.Error("There is no Toby")
		return
	}

	rs := getRecommendedItems(mUser["Toby"], rsMap)

	if len(rs) != len(resultRank) {
		t.Error("getRecommendedItems calc error, failed!")
		return
	}

	for i := 0; i < len(rs); i++ {
		if rs[i] != resultRank[i] {
			t.Error("getRecommendedItems calc error, failed!")
			return
		}

	}

	rs = getRecommendedItems(mUser["Lisa Rose"], rsMap)
	if len(rs) != 0 {
		t.Error("getRecommendedItems calc error, failed!")
		return
	}

	t.Log("getRecommendedItems pass.")
}
