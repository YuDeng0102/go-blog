package task

import (
	"encoding/json"
	"server/global"
	"server/utils/hotSearch"
	"time"

	"go.uber.org/zap"
)

func GetHotListSyncTask() error {
	sourceStrs := []string{"baidu", "zhihu", "kuaishou", "toutiao"}
	for _, sourceStr := range sourceStrs {
		source := hotSearch.NewSource(sourceStr)
		hotSearchData, err := source.GetHotSearchData(30)
		if err != nil {
			global.Log.Error("Failed to get hot search data:", zap.Error(err))
			return err
		}
		data, err := json.Marshal(hotSearchData)
		if err != nil {
			global.Log.Error("Failed to marshal hot search data:", zap.Error(err))
			return err
		}
		if err := global.Redis.Set(sourceStr, data, time.Hour).Err(); err != nil {
			global.Log.Error("Failed to set hot search data in Redis:", zap.Error(err))
			return err
		}
	}
	return nil

}
