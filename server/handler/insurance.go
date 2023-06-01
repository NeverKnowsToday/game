package handler

import (
	"context"
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
		ProductID: insuranceRes.ProductID,
	}
	ctx := context.Background()
	// 创建并发队列
	taskQueue := task.NewTaskQueue()
	// 把并发任务ID加入队列
	for _, id := range insuranceRes.ProductID {
		taskQueue.Append(id)
	}
	// 创建要并发的任务
	insuranceOp := task.NewInsuranceOp(taskQueue, c, insuranceInfo)
	taskQueue.Init(insuranceOp, 5)
	// 并发执行
	//ResChan := make(chan *InsurancesInfo, 9)

	for i := 0; i < taskQueue.UnfinishedTaskCnt; i++ {
		go insuranceOp.RunTask(ctx)
	}

	logger.Infof(fmt.Sprintf("start insuranceOp.Wg.Wait()"))

	insuranceOp.Wg.Wait()
	logger.Infof(fmt.Sprintf("end insuranceOp.Wg.Wait()"))


	latestRes := make([]*task.InsurancesInfo, 0)
	for {
		//tmpInsurancesInfo := &task.InsurancesInfo{}
		if value, ok := <-insuranceOp.ResChan; ok {
			switch  value.(type) {
				case *task.InsurancesInfo :
					taskInsurancesInfo := value.(*task.InsurancesInfo)

					//tmpInsurancesInfo.ProductName = taskInsurancesInfo.ProductName
					//tmpInsurancesInfo.ID = taskInsurancesInfo.ID
					//tmpInsurancesInfo.CompanyName = taskInsurancesInfo.CompanyName
					//
					//for _,v := range taskInsurancesInfo.Info {
					//	tmpInsurancesInfo.Info = append(tmpInsurancesInfo.Info, v)
					//}
					latestRes = append(latestRes,taskInsurancesInfo)

			default:
				logger.Infof(fmt.Sprintf("switch  value.(type)= value.(type) not support"))
				common.ErrorResponse(c, "1111", fmt.Sprintf("switch  value.(type)= value.(type) not support") )
			}

	    }else{
			logger.Infof(fmt.Sprintf("value, ok := <-insuranceOp.ResChan; is not ok"))
			break
		}
	}

	close(insuranceOp.ResChan)

	logger.Infof(fmt.Sprintf("len(latestRes) =  %d",len(latestRes)))
	logger.Infof(fmt.Sprintf("latestRes =  %+v",latestRes))

	for _, v := range latestRes{
		logger.Infof(fmt.Sprintf("latestRes v =  %+v",   v))
	}


	common.Response( c, nil, common.E_Success, latestRes)
	logger.Infof(fmt.Sprintf("end api InsuranceCompare"))

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