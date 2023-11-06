package handler

import (
	"fmt"
	"github.com/game/server/common"
	"github.com/game/server/db"
	"github.com/game/server/task"
	"github.com/gin-gonic/gin"
)

func InsuranceCompare(c *gin.Context) {
	//Operation := "保险对比"
	//token := common.Header(c, "token")
	//tokenInfo, err := public_db.GetToken(token)
	//if common.CheckError("", Operation, c, err, common.E_UserTokenInvalid) {
	//	return
	//}

	logger.Infof(fmt.Sprintf("start api InsuranceCompare"))


	var insuranceRes InsuranceRes
	if common.CheckError(c, common.UnmarshalBody(c, &insuranceRes), common.E_ParseBodyFailed) {
		return
	}

	insuranceInfo := task.InsuranceInfo{
		Age: insuranceRes.Age,
		Total: insuranceRes.Total,
		Sex: insuranceRes.Sex,
		Period: insuranceRes.Period,
		ProductID: insuranceRes.ProductID,
	}

	latestRes := make([]*task.InsurancesInfo, 0)
    for _, id:= range insuranceInfo.ProductID {
		//tmpInsurancesInfo := &task.InsurancesInfo{}
		taskInsurancesInfo, err := task.InsuranceSerial(c,id, insuranceRes.Age, insuranceRes.Total, insuranceInfo.Sex,insuranceInfo.Period)
		if err != nil {
			logger.Infof(fmt.Sprintf("err := InsuranceSerial %s", err.Error()))
			common.ErrorResponse(c, "1111", fmt.Sprintf("eInsuranceSerial err = %s", err.Error()) )
		}

		//tmpInsurancesInfo.ProductName = taskInsurancesInfo.ProductName
		//				tmpInsurancesInfo.ID = taskInsurancesInfo.ID
		//				tmpInsurancesInfo.CompanyName = taskInsurancesInfo.CompanyName
		//				//tmpInsurancesInfo.Info = taskInsurancesInfo.Info
		//				for _,v := range taskInsurancesInfo.Info {
		//					tmpInsurancesInfo.Info = append(tmpInsurancesInfo.Info, v)
		//				}

		latestRes = append(latestRes, taskInsurancesInfo)
	}


	common.Response( c, nil, common.E_Success, latestRes)

}


func GetExcelList(c *gin.Context){
	logger.Infof("start api GetExcelList")

	excelList, err := db.GetProductExcelParsers()
	if err != nil {
		common.ErrorResponse(c, common.E_GetFailed, err.Error())
		return
	}

	for k, v := range excelList {
		logger.Infof(fmt.Sprintf("start api GetExcelList[%d] = %+v", k, v))
		logger.Infof(fmt.Sprintf("start api GetExcelList[%d] = %s", k, v.ExcelPath))
	}

	common.Response( c, nil, common.E_Success, excelList)

}




