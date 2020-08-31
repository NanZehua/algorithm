package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
)

type RecordType map[string]map[string]int

type Record struct {
	userTags RecordType // 用户打过标签的次数
	tagsItem RecordType // 音乐被打过标签的次数，代表歌曲流行度
	tagsUser RecordType // 标签被用户标记次数
	itemUser RecordType // 音乐被不同用户标记次数
}

type tuple struct {
	value     string
	recommend float64
}

func newRecord() *Record {
	var data = &Record{
		userTags: make(map[string]map[string]int),
		tagsItem: make(map[string]map[string]int),
		tagsUser: make(map[string]map[string]int),
		itemUser: make(map[string]map[string]int),
	}

	return data
}

func (rd *Record) initStat(records []string) {
	var (
		fields []string
	)

	for _, record := range records {
		/*
			## ------------------------------------
			## | fields[0] | fields[1] | fields[2] |
			## ------------------------------------
			## | user      | item      | tag       |
			## ------------------------------------
			## | 南泽华     | 一曲相思   | 流行       |
			## ------------------------------------
		*/
		fields = strings.Fields(record)
		//if len(fields) != 3 {
		//	msg := fmt.Sprintf("data format anomaly: [%s]", fields)
		//	logrus.Errorf(msg)
		//	return fmt.Errorf(msg)
		//}

		if _, exist := rd.userTags[fields[0]]; exist {
			rd.userTags[fields[0]][fields[2]] = rd.userTags[fields[0]][fields[2]] + 1
		} else {
			rd.userTags[fields[0]] = map[string]int{
				fields[2]: rd.userTags[fields[0]][fields[2]] + 1,
			}
		}

		if _, exist := rd.tagsItem[fields[2]]; exist {
			rd.tagsItem[fields[2]][fields[1]] = rd.tagsItem[fields[2]][fields[1]] + 1
		} else {
			rd.tagsItem[fields[2]] = map[string]int{
				fields[1]: rd.tagsItem[fields[2]][fields[1]] + 1,
			}
		}

		if _, exist := rd.tagsUser[fields[2]]; exist {
			rd.tagsUser[fields[2]][fields[0]] = rd.tagsUser[fields[2]][fields[0]] + 1
		} else {
			rd.tagsUser[fields[2]] = map[string]int{
				fields[0]: rd.tagsUser[fields[2]][fields[0]] + 1,
			}
		}

		if _, exist := rd.itemUser[fields[1]]; exist {
			rd.itemUser[fields[1]][fields[0]] = rd.itemUser[fields[1]][fields[0]] + 1
		} else {
			rd.itemUser[fields[1]] = map[string]int{
				fields[0]: rd.itemUser[fields[1]][fields[0]] + 1,
			}
		}
	}

	return
}

func (rd *Record) recommend(user string, K int) string {
	var recommendItems = make(map[string]float64)
	for tag, wut := range rd.userTags[user] {
		for item, wti := range rd.tagsItem[tag] {
			if _, exist := recommendItems[item]; !exist {
				// 计算用户对物品兴趣度
				recommendItems[item] = (float64(wut) / math.Log(float64(1+len(rd.tagsUser[tag])))) *
					(float64(wti) / math.Log(float64(1+len(rd.itemUser[item]))))
			} else {
				recommendItems[item] += float64(wti) / math.Log(float64(1+len(rd.tagsUser[tag])))
			}
		}
	}

	rec := sorted(recommendItems, false)
	fmt.Println(">>>>>>", rec)

	var music []string
	for i := 0; i < K; i++ {
		music = append(music, rec[i].value)
	}

	return strings.Join(music, "/")
}

func loadData(filePath string) []string {
	fi, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	var records []string
	br := bufio.NewReader(fi)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}

		records = append(records, string(line))

	}

	return records
}

func sorted(mp map[string]float64, key bool) []tuple {
	var (
		newMpList []tuple
	)

	if key {
		var newMpKs = make([]string, 0)

		for key := range mp {
			newMpKs = append(newMpKs, key)
		}

		sort.Strings(newMpKs)
		for _, v := range newMpKs {
			newMpList = append(newMpList, tuple{
				value:     v,
				recommend: mp[v],
			})
		}
	} else {
		var (
			index = 0
			ps    = make(PairList, len(mp))
		)

		for k, v := range mp {
			ps[index] = Pair{k, v}
			index++
		}
		sort.Sort(ps)

		for _, p := range ps {
			newMpList = append(newMpList, tuple{
				value:     p.Key,
				recommend: p.Value,
			})
		}
	}

	return newMpList
}

type Pair struct {
	Key   string
	Value float64
}

type PairList []Pair

func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p PairList) Len() int { return len(p) }

func (p PairList) Less(i, j int) bool { return p[j].Value < p[i].Value }

func main() {
	var file = "/home/nanzehua/Demo/test/algorithm/tag/data/data.txt"
	records := loadData(file)
	fmt.Println("====================")
	rd := newRecord()
	rd.initStat(records)
	fmt.Println("用户打过标签的次数: ", rd.userTags)
	fmt.Println("====================")
	fmt.Println("音乐打过标签的次数: ", rd.tagsItem)
	fmt.Println("====================")
	fmt.Println("标签被用户使用次数: ", rd.tagsUser)
	fmt.Println("====================")
	fmt.Println("音乐被用户标记次数: ", rd.itemUser)
	fmt.Println("====================")
	rec := rd.recommend("南泽华", 2)
	fmt.Println("====================")
	fmt.Println("推荐歌曲: ", rec)
}
