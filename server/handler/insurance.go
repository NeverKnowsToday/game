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
	for i := 0; i < taskQueue.UnfinishedTaskCnt; i++ {
		go insuranceOp.RunTask(ctx)
	}

	insuranceOp.Wg.Wait()
	latestRes := make([]*InsurancesInfo, 0)
	for {
		if value, ok := <-insuranceOp.ResChan; ok {
			insurancesInfo := value.(InsurancesInfo)
			latestRes = append(latestRes,&insurancesInfo)
	    }else{
			break
		}
	}

	close(insuranceOp.ResChan)




	common.Response( c, nil, common.E_Success, latestRes)
	logger.Infof(fmt.Sprintf("end api InsuranceCompare"))
//
//	config := &InputConfig{
//		TaskID:            insuranceInfo.ProductID[0],
//		ExcelFilePath:     "/Users/zhanggong02/python/xjw/tmp/333.xlsx",
//		ExcelFileTempPath: "/Users/zhanggong02/python/xjw/history/333_123.xlsx",
//		Set: &Set{
//			Age: &SetInfo{
//				Value:    "39",
//				Location: "F12",
//				Sheet:    "资料输入",
//			},
//			Total: &SetInfo{
//				Value:    "180000",
//				Location: "H15",
//				Sheet:    "资料输入",
//			},
//		},
//		Get: &Get{
//			BenefitDiff: &BenefitDiff{
//				Location: "A8:Q114",
//				Sheet:    "利益演示",
//			},
//		},
//	}
//
//	configByte, err := json.Marshal(config)
//	if err != nil {
//		logger.Infof("json.Marshal(config) failed")
//	}
//
//	f, err := os.Create("./111.config")
//	if err != nil {
//		logger.Infof("create config failed")
//
//	}
//	defer f.Close()
//
//	n, err := f.WriteString(string(configByte))
//	if err != nil {
//		logger.Infof("f.WriteString(string(configByte)) failed")
//
//	}
//
//	fmt.Println("n = ", n)
//
//
//	command := "python D:/python/xjw/main.py D:/python/xjw/app.conf"
//	timeout := time.Duration(20)
//	execOut, execErr, err := linux.ExecWithShellTimeout(command, timeout)
//	if err != nil {
//		execErr := fmt.Errorf("execute [%s]\n execOut: %s\n execErr: %s \n error: %s", command, execOut, execErr, err)
//	//	logs.GoLogError(ctx, "error message: %s", execErr.Error())
//		logger.Infof("error message: %s", execErr.Error())
//
//		fmt.Printf("error message: %s", execErr.Error())
//		common.ErrorResponse(c, "1111", err.Error())
//
//		//return ghttp.NewJSONResponse(http.StatusExpectationFailed, execErr)
//		//return execErr
//	}
//
//	fmt.Printf("execOut ====================== = %+v\n", execOut)
//
//	 data := execOut
//
//
//	tmp := InsurancesResFromPython{}
//	err = json.Unmarshal([]byte(data), &tmp)
//	if err != nil {
//		fmt.Println("err ============================== ", err.Error())
//	}
//
//	res := make([]*ExcelInsurance, 0)
//
//
//	for _, v := range tmp.ExcelData.BenefitDisplay {
//		res_row := &ExcelInsurance{}
//
//		for _, value := range v {
//			fmt.Printf("for _, value := range v================== = %+v\n", value.Name)
//			fmt.Printf("for _, value := range value================== = %+v\n", value.Value)
//
//			if value.Name == "保单年度" {
//				fmt.Printf("**************==================************** = %+v\n", value.Name)
//			}
//			switch value.Name {
//			case "PolicyYear":
//				res_row.PolicyYear = value.Value
//			case "AgeAtTheEndOfPolicyYear":
//				res_row.AgeAtTheEndOfPolicyYear = value.Value
//			case "PremiumPaymentInTheYearBeforePolicyReduction":
//				res_row.PremiumPaymentInTheYearBeforePolicyReduction = value.Value
//			case "AccumulatedPremiumsBeforePolicyReduction":
//				res_row.AccumulatedPremiumsBeforePolicyReduction = value.Value
//			case "BasicSumInsuredForTheYearBeforePolicyReduction":
//				res_row.BasicSumInsuredForTheYearBeforePolicyReduction = value.Value
//			case "EffectiveSumInsuredForTheYearBeforePolicyReduction":
//				res_row.EffectiveSumInsuredForTheYearBeforePolicyReduction = value.Value
//			case "DeathOrTotalDisabilityInsuranceBenefitsAtTheEndOfThePolicyYearBeforePolicyReduction":
//				res_row.DeathOrTotalDisabilityInsuranceBenefitsAtTheEndOfThePolicyYearBeforePolicyReduction = value.Value
//			case "CashValueAtTheEndOfThePolicyYearBeforePolicyReduction":
//				res_row.CashValueAtTheEndOfThePolicyYearBeforePolicyReduction = value.Value
//			case "PolicyYearEndReductionAmountOptionalEntryColumn":
//				res_row.PolicyYearEndReductionAmountOptionalEntryColumn = value.Value
//			case "DecreaseAmountOfBasicSumInsuredAtTheEndOfThePolicyYear":
//				res_row.DecreaseAmountOfBasicSumInsuredAtTheEndOfThePolicyYear = value.Value
//			case "CumulativeReductionAmountAtTheEndOfThePolicyYear":
//				res_row.CumulativeReductionAmountAtTheEndOfThePolicyYear = value.Value
//			case "AnnualPremiumPaymentAfterPolicyReduction":
//				res_row.AnnualPremiumPaymentAfterPolicyReduction = value.Value
//			case "AccumulatedPremiumsAfterPolicyReduction":
//				res_row.AccumulatedPremiumsAfterPolicyReduction = value.Value
//			case "AnnualBasicInsuranceAmountAfterPolicyReduction":
//				res_row.AnnualBasicInsuranceAmountAfterPolicyReduction = value.Value
//			case "AnnualEffectiveInsuredAmountAfterPolicyReduction":
//				res_row.AnnualEffectiveInsuredAmountAfterPolicyReduction = value.Value
//			case "DeathOrTotalDisabilityInsuranceBenefittTheEndOfThePolicyYearAfterPolicyReduction":
//				res_row.DeathOrTotalDisabilityInsuranceBenefittTheEndOfThePolicyYearAfterPolicyReduction = value.Value
//			case "CashValueAtTheEndOfThePolicyYearAfterPolicyReduction":
//				res_row.CashValueAtTheEndOfThePolicyYearAfterPolicyReduction = value.Value
//			}
//
//		}
//		res = append(res, res_row)
//	}
//
//
//	insuranceInfo1 := &InsurancesInfo{
//		ID:          "1",
//		CompanyName: "平安",
//		ProductName: "放心投",
//		Info:        res,
//	}
//
//	insuranceInfo2 := &InsurancesInfo{
//		ID:          "2",
//		CompanyName: "人寿",
//		ProductName: "大胆投",
//		Info:        res,
//	}
//
//	insuranceInfo3 := &InsurancesInfo{
//		ID:          "3",
//		CompanyName: "水滴",
//		ProductName: "投投投",
//		Info:        res,
//	}
//
//
//	latestRes := make([]*InsurancesInfo, 0)
//
//	latestRes = append(latestRes, insuranceInfo1, insuranceInfo2, insuranceInfo3)
//
//	common.Response( c, nil, common.E_Success, latestRes)
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