package repositories

import (
	"architecture_template/common_dtos/response"
	"architecture_template/constants/notis"
	"architecture_template/helper"
	user_notis "architecture_template/services/user/constants/notis"
	redis_key "architecture_template/services/user/constants/redisKey"
	"architecture_template/services/user/entities"
	"context"
	"errors"
	"fmt"
)

func (tr *repo) GetUsersByRole(id string, c context.Context) (*[]entities.User, error) {
	var key string = fmt.Sprintf(redis_key.GetByRoleKey, id)

	// Retrieve data from redis cache
	if res, err, isValid := helper.GetDataFromRedis[[]entities.User](tr.redisCache, key, c); isValid {
		return res, err
	} else {
		tr.logger.Print(user_notis.RedisMsg + err.Error()) // Fetching data from cache meets problem
	}
	//------------------------------------------

	// Retrieve database
	var errLogMsg string = user_notis.UserRepoMsg + "GetUsersByRole - "
	var query string = "Select id, email, password, roleId, activeStatus from " + entities.GetTable() + " where roleId = ?"
	var internalErr error = errors.New(notis.InternalErr)
	var res *[]entities.User

	rows, err := tr.db.Query(query, id)
	if err != nil {
		go func() {
			if helper.SaveDataToRedis(tr.redisCache, key, response.DataStorage{
				ErrMsg: internalErr,
			}, c) != nil {
				tr.logger.Print(notis.RedisMsg)
			}
		}()

		tr.db.Close()
		tr.logger.Print(errLogMsg, err)
		return nil, internalErr
	}
	defer rows.Close()

	for rows.Next() {
		var x entities.User

		if err := rows.Scan(&x.UserId, &x.Email, &x.Pasword, &x.RoleId, &x.ActiveStatus); err != nil {
			go func() {
				if helper.SaveDataToRedis(tr.redisCache, key, response.DataStorage{
					ErrMsg: internalErr,
				}, c) != nil {
					tr.logger.Print(notis.RedisMsg)
				}
			}()

			tr.db.Close()
			tr.logger.Print(errLogMsg, err)
			return nil, internalErr
		}
		*res = append(*res, x)
	}

	if helper.SaveDataToRedis(tr.redisCache, fmt.Sprintf(redis_key.GetByRoleKey, id), response.DataStorage{
		Data: res,
	}, c) != nil {
		tr.logger.Print(notis.RedisMsg + helper.ConvertModelToString(res))
	}

	return res, nil
}