// user MovieLens data to implement item-base Collaborative Filtering.
//
// 1 load MovieLens data file u.data
// 2 to make the loaded data as map like {item: {usr1:3, usr2:4}}
// 3 implement item-base Collaborative Filtering
package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type dataSetType map[string]map[string]float32
type dataCellType map[string]float32

// load MovieLens data file u.data
// and to make the loaded data as
//     itemMap like {itemID: {usrID1:3, usrID2:4}}
// and userMap like {userID: {itemID1:3, itemID2:4}}
func loaddata() (itemMap, userMap dataSetType, err error) {

	// get item name list
	var items []string

	fitem, err := os.Open("u.item")
	if err != nil {
		return nil, nil, err
	}

	defer fitem.Close()

	scan := bufio.NewScanner(fitem)
	for scan.Scan() {
		name := strings.Split(scan.Text(), "|")
		items = append(items, name[1])
	}

	// parse data
	var mUser, mItem dataCellType
	var itemName string
	var itemId int

	fdata, err := os.Open("u.data")
	if err != nil {
		return nil, nil, err
	}

	defer fdata.Close()

	itemMap = make(dataSetType)
	userMap = make(dataSetType)

	scan = bufio.NewScanner(fdata)
	for scan.Scan() {
		s := strings.Fields(scan.Text())
		score, err := strconv.ParseFloat(s[2], 32)
		if err != nil {
			break
		}

		itemId, err = strconv.Atoi(s[1])
		if err != nil {
			break
		}

		if itemId > len(items) || itemId == 0 {
			err = errors.New("itemId overflow.")
			break
		} else {
			itemName = items[itemId-1]
		}

		// construct itemMap
		mUser = itemMap[itemName]
		if mUser == nil {
			mUser = make(dataCellType)
		}

		mUser[s[0]] = float32(score)

		itemMap[itemName] = mUser

		// construct userMap
		mItem = userMap[s[0]]
		if mItem == nil {
			mItem = make(dataCellType)
		}

		mItem[itemName] = float32(score)

		userMap[s[0]] = mItem
	}

	return itemMap, userMap, err
}

func loadExampleData() (itemMap, userMap dataSetType, err error) {
	// Lisa Rose
	m1 := dataCellType{
		"Lady in the Water":  2.5,
		"Snakes on a Plane":  3.5,
		"Just My Luck":       3.0,
		"Superman Returns":   3.5,
		"You, Me and Dupree": 2.5,
		"The Night Listener": 3.0,
	}
	// Gene Seymour
	m2 := dataCellType{
		"Lady in the Water":  3.0,
		"Snakes on a Plane":  3.5,
		"Just My Luck":       1.5,
		"Superman Returns":   5.0,
		"You, Me and Dupree": 3.5,
		"The Night Listener": 3.0,
	}
	// Michael Phillips
	m3 := dataCellType{
		"Lady in the Water":  2.5,
		"Snakes on a Plane":  3.0,
		"Superman Returns":   3.5,
		"The Night Listener": 4.0,
	}
	// Claudia Puig
	m4 := dataCellType{
		"Snakes on a Plane":  3.5,
		"Just My Luck":       3.0,
		"Superman Returns":   4.0,
		"You, Me and Dupree": 2.5,
		"The Night Listener": 4.5,
	}
	// Mick LaSalle
	m5 := dataCellType{
		"Lady in the Water":  3.0,
		"Snakes on a Plane":  4.0,
		"Just My Luck":       2.0,
		"Superman Returns":   3.0,
		"You, Me and Dupree": 2.0,
		"The Night Listener": 3.0,
	}
	// Jack Matthews
	m6 := dataCellType{
		"Lady in the Water":  3.0,
		"Snakes on a Plane":  4.0,
		"Superman Returns":   5.0,
		"You, Me and Dupree": 3.5,
		"The Night Listener": 3.0,
	}
	// Toby
	m7 := dataCellType{
		"Snakes on a Plane":  4.5,
		"Superman Returns":   4.0,
		"You, Me and Dupree": 1.0,
	}

	userMap = dataSetType{
		"Lisa Rose":        m1,
		"Gene Seymour":     m2,
		"Michael Phillips": m3,
		"Claudia Puig":     m4,
		"Mick LaSalle":     m5,
		"Jack Matthews":    m6,
		"Toby":             m7,
	}

	itemMap = make(dataSetType)

	for usr, m := range userMap {
		for k, v := range m {
			mTemp := itemMap[k]
			if mTemp == nil {
				mTemp = make(dataCellType)
			}

			mTemp[usr] = v
			itemMap[k] = mTemp
		}
	}

	return itemMap, userMap, nil
}
