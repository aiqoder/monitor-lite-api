package result

import (
	"fmt"
	"github.com/duke-git/lancet/v2/strutil"
	"gorm.io/gorm"
	"slices"
	"strconv"
	"strings"
)

// 用于解析通用的map参数，通常用在分页参数解析的地方

type Pager struct {
	Current int `json:"current"`
	Size    int `json:"size"`
}

func Map2Page(m map[string]string) Pager {
	pager := Pager{}
	current, ok := m["current"]

	if ok {
		pager.Current, _ = strconv.Atoi(current)
	}

	size, ok := m["size"]

	if ok {
		pager.Size, _ = strconv.Atoi(size)
	}

	return pager
}

type Map2ScopeIgnore struct {
	ImportantKeys []string
	IgnoreVals    []string
	IgnoreKeys    []string
	NoLikeKeys    []string          // 所有的值除，id之外默认都进行like匹配，这里填入需要精确匹配的列
	FieldZero     map[string]string // 数据类型维护，主要用户零值判断，默认认为0和""是零值
}

// Map2ScopeIgnoreOption t 表示key / val， val 表示需要忽略的值
type Map2ScopeIgnoreOption func(t *Map2ScopeIgnore)

func NewMapScopeIgnore(opts ...Map2ScopeIgnoreOption) *Map2ScopeIgnore {
	ignore := Map2ScopeIgnore{
		ImportantKeys: []string{},                  // 不需要判断是否为0值
		IgnoreKeys:    []string{"current", "size"}, // 丢弃这些key，不在sql当中查询
		NoLikeKeys:    []string{"id"},
		FieldZero: map[string]string{
			"id": "0",
		},
	}

	for _, opt := range opts {
		opt(&ignore)
	}

	return &ignore
}

// WithImportantKeys 某个键不进行默认值判断
func WithImportantKeys(keys ...string) Map2ScopeIgnoreOption {
	return func(t *Map2ScopeIgnore) {
		// 将key 添加到重要key当中
		for _, key := range keys {
			if !slices.Contains(t.ImportantKeys, key) {
				t.ImportantKeys = append(t.ImportantKeys, key)
			}
		}
	}
}

func WithFieldZero(m map[string]string) Map2ScopeIgnoreOption {
	return func(t *Map2ScopeIgnore) {
		// 将key 添加到重要key当中
		for key, val := range m {
			t.FieldZero[key] = val
		}
	}
}
func WithNoLikeKeysKeys(keys ...string) Map2ScopeIgnoreOption {
	return func(t *Map2ScopeIgnore) {
		// 将key 添加到重要key当中
		for _, key := range keys {
			if !slices.Contains(t.NoLikeKeys, key) {
				t.NoLikeKeys = append(t.NoLikeKeys, key)
			}
		}
	}
}

func WithIgnoreKeys(keys ...string) Map2ScopeIgnoreOption {
	return func(t *Map2ScopeIgnore) {
		for _, key := range keys {
			if !slices.Contains(t.IgnoreKeys, key) {
				t.IgnoreKeys = append(t.IgnoreKeys, key)
			}
		}
	}
}

func Map2ScopePager(m map[string]string, opts ...Map2ScopeIgnoreOption) func(db *gorm.DB) *gorm.DB {
	pager := Map2Page(m)
	return func(db *gorm.DB) *gorm.DB {
		db = db.Scopes(Map2ScopeWhere(m, opts...))
		return db.Offset((pager.Current - 1) * pager.Size).Limit(pager.Size)
	}
}

func Map2ScopeWhere(m map[string]string, opts ...Map2ScopeIgnoreOption) func(db *gorm.DB) *gorm.DB {
	ignore := NewMapScopeIgnore(opts...)
	return func(db *gorm.DB) *gorm.DB {
		for key, val := range m {
			if slices.Contains(ignore.IgnoreKeys, key) {
				continue
			}

			// 未设置0值，使用默认0值

			if _, ok := ignore.FieldZero[key]; ok {
				if val == "" {
					continue
				}

				if val == "0" {
					continue
				}
			} else if ignore.FieldZero[key] == val {
				continue
			}

			snakeCaseKey := strutil.SnakeCase(key)

			if val == " " {
				db = db.Where(fmt.Sprintf("`%s` = ? OR `%s` IS NULL", snakeCaseKey, snakeCaseKey), "")
				continue
			}

			// 是否是全字段匹配
			fullFields := slices.Contains(ignore.NoLikeKeys, key)

			// 逗号分割的数据，视为数组
			for index, v := range strings.Split(val, ",") {
				if fullFields {
					if index == 0 {
						db = db.Where(fmt.Sprintf("`%s` = ?", snakeCaseKey), v)
					} else {
						db = db.Or(fmt.Sprintf("`%s` = ?", snakeCaseKey), v)
					}
				} else {
					if index == 0 {
						db = db.Where(fmt.Sprintf("`%s` LIKE ?", snakeCaseKey), fmt.Sprintf("%%%s%%", v))
					} else {
						db = db.Or(fmt.Sprintf("`%s` LIKE ?", snakeCaseKey), fmt.Sprintf("%%%s%%", v))
					}
				}
			}
		}
		return db
	}
}
