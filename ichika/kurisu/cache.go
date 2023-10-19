package kurisu

import (
	"encoding/json"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/yunyun"
	"golang.org/x/exp/slices"
)

type PageCache map[string]yunyun.Page

func BuildKey(p yunyun.Page) string {
	return string(p.File)
}

type CacheManager struct {
	Cache     PageCache
	CacheFile yunyun.FullPathFile
}

func (cm *CacheManager) Load() error {
	cm.Cache = PageCache{}

	file, err := os.ReadFile(string(cm.CacheFile))
	if err != err {
		return err
	}

	return json.Unmarshal(file, &cm.Cache)
}

func (cm *CacheManager) Save() error {
	cacheAsString, err := json.Marshal(cm.Cache)
	puck.Logger.Debug("Saving cache:", "size", len(cm.Cache), "length", len(cacheAsString))

	if err == nil {
		return os.WriteFile(string(cm.CacheFile), cacheAsString, os.ModeType)
	}
	return err
}

func (cm *CacheManager) GetPage(page yunyun.Page) yunyun.Page {
	retrieved := cm.Cache[BuildKey(page)]
	return retrieved
}

func (cm *CacheManager) MergePage(page yunyun.Page) error {
	key := BuildKey(page)

	if val, ok := cm.Cache[key]; ok {
		cm.Cache[key] = mergePageDeltas(val, page)
	} else {
		cm.Cache[key] = page
	}
	return nil
}

func mergePageDeltas(oldPage yunyun.Page, newPage yunyun.Page) yunyun.Page {
	mergedPage := yunyun.Page(newPage)
	if cmp.Equal(oldPage, newPage) {
		return mergedPage
	}

	mergedDeltas := append(oldPage.Changes, newPage.Changes...)

	slices.SortFunc[[]yunyun.Change](mergedDeltas, yunyun.Change.Compare)

	finalDeltas := make([]yunyun.Change, len(mergedDeltas))
	c := -1
	for i, o := range mergedDeltas {
		if c >= 0 {
			if finalDeltas[c] != o {
				finalDeltas[i] = o
			}
		} else {
			finalDeltas[i] = o
		}
		c = i
	}
	mergedPage.Changes = slices.Compact[[]yunyun.Change](finalDeltas)
	return mergedPage
}
