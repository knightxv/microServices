package db

import (
	"context"
	"time"
)

const (
	accountTempCode = "ACCOUNT_TEMP_CODE"
)

func (d *DataBases) JudgeAccountEXISTS(account string) (bool, error) {
	key := accountTempCode + account
	n, err := d.RDB.Exists(context.Background(), key).Result()
	if n > 0 {
		return true, err
	} else {
		return false, err
	}
}
func (d *DataBases) SetAccountCode(account string, code, ttl int) (err error) {
	key := accountTempCode + account
	return d.RDB.Set(context.Background(), key, code, time.Duration(ttl)*time.Second).Err()
}
func (d *DataBases) GetAccountCode(account string) (string, error) {
	key := accountTempCode + account
	return d.RDB.Get(context.Background(), key).Result()
}
