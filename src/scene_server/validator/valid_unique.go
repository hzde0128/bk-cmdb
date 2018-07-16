/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package validator

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// validCreateUnique  valid create unique
func (valid *ValidMap) validCreateUnique(valData map[string]interface{}) error {
	if 0 >= len(valid.isOnly) {
		blog.V(3).Infof("is only array is zero for %s", valid.objID)
		return nil
	}

	objID := valid.objID
	searchCond := make(map[string]interface{})
	// only search data not in diable status
	searchCond[common.BKDataStatusField] = map[string]interface{}{common.BKDBNE: common.DataStatusDisabled}

	// retrive isonly value
	for key, val := range valData {
		if valid.isOnly[key] {
			searchCond[key] = val
		}
	}

	if innerObject[valid.objID] {
		objID = common.BKINnerObjIDObject
		searchCond[common.BKObjIDField] = valid.objID
	}

	if 0 == len(searchCond) {
		return nil
	}

	result, err := valid.CoreAPI.ObjectController().Instance().SearchObjects(valid.ctx, objID, valid.pheader, &metadata.QueryInput{Condition: searchCond})
	if nil != err {
		return err
	}
	if !result.Result {
		return valid.errif.Error(result.Code)
	}

	if 0 >= result.Data.Count {
		blog.Error("duplicate data ")
		return valid.errif.Error(common.CCErrCommDuplicateItem)
	}

	return nil
}

// validUpdateUnique valid update unique
func (valid *ValidMap) validUpdateUnique(valData map[string]interface{}, instID int64) error {
	if 0 >= len(valid.isOnly) {
		return nil
	}

	objID := valid.objID
	searchCond := make(map[string]interface{})
	mapData, err := valid.getInstDataByID(instID)
	if nil != err {
		return err
	}

	// retrive isonly value
	for key, val := range mapData {
		if valid.isOnly[key] {
			searchCond[key] = val
		}
	}
	for key, val := range valData {
		if valid.isOnly[key] {
			searchCond[key] = val
		}
	}

	objIDName := util.GetObjIDByType(objID)
	searchCond[objIDName] = map[string]interface{}{common.BKDBNE: instID}
	// only search data not in diable status
	searchCond[common.BKDataStatusField] = map[string]interface{}{common.BKDBNE: common.DataStatusDisabled}

	if !innerObject[valid.objID] {
		searchCond[common.BKObjIDField] = valid.objID
		objID = common.BKINnerObjIDObject
	}

	result, err := valid.CoreAPI.ObjectController().Instance().SearchObjects(valid.ctx, objID, valid.pheader, &metadata.QueryInput{Condition: searchCond})
	if nil != err {
		return err
	}
	if !result.Result {
		return valid.errif.Error(result.Code)
	}

	if 0 >= result.Data.Count {
		blog.Error("duplicate data ")
		return valid.errif.Error(common.CCErrCommDuplicateItem)
	}
	return nil
}

// getInstDataByID get inst data by id
func (valid *ValidMap) getInstDataByID(instID int64) (map[string]interface{}, error) {
	objID := valid.objID
	searchCond := make(map[string]interface{})

	if innerObject[valid.objID] {
		objID = common.BKINnerObjIDObject
		searchCond[common.BKObjIDField] = objID
		searchCond[common.BKInstIDField] = instID
	} else {
		objIDName := util.GetObjIDByType(objID)
		searchCond[objIDName] = instID
	}

	result, err := valid.CoreAPI.ObjectController().Instance().SearchObjects(valid.ctx, objID, valid.pheader, &metadata.QueryInput{Condition: searchCond})
	if nil != err {
		return nil, err
	}
	if !result.Result {
		return nil, valid.errif.Error(result.Code)
	}
	if len(result.Data.Info) == 0 {
		return nil, nil
	}

	if len(result.Data.Info[0]) > 0 {
		return result.Data.Info[0], nil
	}
	return nil, valid.errif.Error(common.CCErrCommNotFound)
}
