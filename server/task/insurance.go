package task

import (
	"encoding/json"
	"fmt"
	"github.com/game/server/common"
	"github.com/game/server/db"
	"github.com/game/server/utils/linux"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

//var ResChan  chan *InsurancesInfo


type InsuranceOp struct {
	*TaskQueue
	Ctx   *gin.Context
	InsuranceInfo InsuranceInfo
	//ResChan       *InsurancesInfo
}

func NewInsuranceOp(t *TaskQueue, c *gin.Context, insuranceInfo InsuranceInfo) *InsuranceOp{
	return &InsuranceOp{
		TaskQueue: t,
		Ctx: c,
		InsuranceInfo:    insuranceInfo,
		//ResChan:  make(chan *InsurancesInfo, 9),
	}
}

func (op *InsuranceOp)DbpaasRun() error {
	logger.Infof(fmt.Sprintf("InsuranceOp start"))

	id := op.Get().(string)
	logger.Infof(fmt.Sprintf("InsuranceOp id = %s ", id))

	excel,err := db.GetProductExcelParserByID(id)
	if err != nil {
		common.ErrorResponse(op.Ctx, common.E_GetFailed, err.Error())
		return err
	}

	logger.Infof(fmt.Sprintf("InsuranceOp db excel = %+v ", excel))
	logger.Infof(fmt.Sprintf("InsuranceOp db  excel.ProductId = %s ", excel.ProductId))

	logger.Infof(fmt.Sprintf("InsuranceOp db excel.SetInfo = %s ", excel.SetInfo))
	logger.Infof(fmt.Sprintf("InsuranceOp db excel.getInfo = %s ", excel.GetInfo))

	unmarshalSetInfo := "{" + excel.SetInfo + "}"
	logger.Infof(fmt.Sprintf("InsuranceOp db unmarshalSetInfo = %s ", unmarshalSetInfo))

	setInfo := new(SetInfo)
	err = json.Unmarshal([]byte(unmarshalSetInfo), &setInfo)
	if err != nil {
		logger.Errorf(fmt.Sprintf("err = json.Unmarshal([]byte(unmarshalSetInfo), &setInfo) %s",err.Error()))

		common.ErrorResponse(op.Ctx, common.E_JsonUnmarshalError, err.Error())
		return err
	}

	logger.Infof(fmt.Sprintf("setInfo = %+v", setInfo))
	logger.Infof(fmt.Sprintf("setInfo.Set.Age.Value = %s", setInfo.Set.Age.Value))
	logger.Infof(fmt.Sprintf("setInfo.Set.Age.Sheet = %s", setInfo.Set.Age.Sheet))
	logger.Infof(fmt.Sprintf("setInfo.Set.Age.Location = %s", setInfo.Set.Age.Location))

	unmarshalGetInfo := "{" + excel.GetInfo + "}"
	logger.Infof(fmt.Sprintf("InsuranceOp db unmarshalGetInfo = %s ", unmarshalGetInfo))

	getInfo := GetInfo{
		Get: &Get{
			BenefitDiff: &BenefitDiff{},
		},
	}
	err = json.Unmarshal([]byte(unmarshalGetInfo), &getInfo)
	if err != nil {
		logger.Errorf(fmt.Sprintf("err = json.Unmarshal([]byte(unmarshalGetInfo), &getInfo) = %s ", err.Error()))
		common.ErrorResponse(op.Ctx, common.E_JsonUnmarshalError, err.Error())
		return err
	}

	logger.Infof(fmt.Sprintf("getInfo = %+v", getInfo))
	logger.Infof(fmt.Sprintf("getInfo.Get.BenefitDiff.Sheet = %s", getInfo.Get.BenefitDiff.Sheet))
	logger.Infof(fmt.Sprintf("getInfo.Get.BenefitDiff.Location = %s", getInfo.Get.BenefitDiff.Location))


	//uid, err := uuid.GenerateUUID()

	//tempPath := strings.TrimRight(excel.ExcelPath, ".xlsx")
	tempPath := "D:\\app\\gin-backend\\history\\" + excel.ProductId + "_" + excel.ExcelName
	logger.Infof(fmt.Sprintf("tempPath = %s", tempPath))

	logger.Infof(fmt.Sprintf("excel.ExcelPath = %s", excel.ExcelPath))

	//excelPath := strings.Replace(excel.ExcelPath, "\\", "\\\\", -1)
	logger.Infof(fmt.Sprintf("excelPath = %s", excel.ExcelPath))
	config := &InputConfig{
		TaskID:            string(id),
		ExcelFilePath:     excel.ExcelPath,
		ExcelFileTempPath: tempPath,
		Set: &Set{
			Age: &SetBase{
				Value:    op.InsuranceInfo.Age,
				Location: setInfo.Set.Age.Location,
				Sheet:    setInfo.Set.Age.Sheet,
			},
			Total: &SetBase{
				Value:    op.InsuranceInfo.Total,
				Location: setInfo.Set.Total.Location,
				Sheet:    setInfo.Set.Total.Sheet,
			},
		},
		Get: &Get{
			BenefitDiff: &BenefitDiff{
				Location: getInfo.Get.BenefitDiff.Location,
				Sheet:    getInfo.Get.BenefitDiff.Sheet,
			},
		},
	}

	configByte, err := json.Marshal(config)
	if err != nil {
		logger.Infof(fmt.Sprintf("json.Marshal(config) failed = %s", err.Error()))

		//logger.Infof("json.Marshal(config) failed")
	}

	configName := excel.ProductId + ".conf"
	logger.Infof(fmt.Sprintf("configName = %s", configName))

	configPath := "D:/app/gin-backend/" + configName
	logger.Infof(fmt.Sprintf("configPath = %s", configPath))

	f, err := os.Create(configName)
	if err != nil {
		logger.Infof(fmt.Sprintf("create config failed = %s", err.Error()))

	//	logger.Infof("create config failed")

	}
	defer f.Close()

	n, err := f.WriteString(string(configByte))
	if err != nil {
		logger.Infof(fmt.Sprintf("f.WriteString(string(configByte)) failed err = %s", err.Error()))
	}

	fmt.Println("n = ", n)

	//configPath := "D:/python/xjw/" + configName
	logger.Infof(fmt.Sprintf("configPath = %s", configPath))

	command := "python D:/python/xjw/main.py " + configPath
	timeout := time.Duration(20)
	execOut, execErr, err := linux.ExecWithShellTimeout(command, timeout)
	if err != nil {
		execErr := fmt.Errorf("execute [%s]\n execOut: %s\n execErr: %s \n error: %s", command, execOut, execErr, err)
		//	logs.GoLogError(ctx, "error message: %s", execErr.Error())
		logger.Infof("error message: %s", execErr.Error())

		//fmt.Printf("error message: %s", execErr.Error())
		common.ErrorResponse(op.Ctx, "1111", err.Error())

		//return ghttp.NewJSONResponse(http.StatusExpectationFailed, execErr)
		//return execErr
	}

	//fmt.Printf("execOut ====================== = %+v\n", execOut)

	data := execOut


	tmp := InsurancesResFromPython{}
	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		fmt.Println("err ============================== ", err.Error())
	}

	res := make([]*ExcelInsurance, 0)


	for _, v := range tmp.ExcelData.BenefitDisplay {
		res_row := &ExcelInsurance{}

		for _, value := range v {
			//fmt.Printf("for _, value := range v================== = %+v\n", value.Name)
			//fmt.Printf("for _, value := range value================== = %+v\n", value.Value)
			//
			//if value.Name == "保单年度" {
			//	fmt.Printf("**************==================************** = %+v\n", value.Name)
			//}
			switch value.Name {
			case "PolicyYear":
				res_row.PolicyYear = value.Value
			case "AgeAtTheEndOfPolicyYear":
				res_row.AgeAtTheEndOfPolicyYear = value.Value
			case "PremiumPaymentInTheYearBeforePolicyReduction":
				res_row.PremiumPaymentInTheYearBeforePolicyReduction = value.Value
			case "AccumulatedPremiumsBeforePolicyReduction":
				res_row.AccumulatedPremiumsBeforePolicyReduction = value.Value
			case "BasicSumInsuredForTheYearBeforePolicyReduction":
				res_row.BasicSumInsuredForTheYearBeforePolicyReduction = value.Value
			case "EffectiveSumInsuredForTheYearBeforePolicyReduction":
				res_row.EffectiveSumInsuredForTheYearBeforePolicyReduction = value.Value
			case "DeathOrTotalDisabilityInsuranceBenefitsAtTheEndOfThePolicyYearBeforePolicyReduction":
				res_row.DeathOrTotalDisabilityInsuranceBenefitsAtTheEndOfThePolicyYearBeforePolicyReduction = value.Value
			case "CashValueAtTheEndOfThePolicyYearBeforePolicyReduction":
				res_row.CashValueAtTheEndOfThePolicyYearBeforePolicyReduction = value.Value
			case "PolicyYearEndReductionAmountOptionalEntryColumn":
				res_row.PolicyYearEndReductionAmountOptionalEntryColumn = value.Value
			case "DecreaseAmountOfBasicSumInsuredAtTheEndOfThePolicyYear":
				res_row.DecreaseAmountOfBasicSumInsuredAtTheEndOfThePolicyYear = value.Value
			case "CumulativeReductionAmountAtTheEndOfThePolicyYear":
				res_row.CumulativeReductionAmountAtTheEndOfThePolicyYear = value.Value
			case "AnnualPremiumPaymentAfterPolicyReduction":
				res_row.AnnualPremiumPaymentAfterPolicyReduction = value.Value
			case "AccumulatedPremiumsAfterPolicyReduction":
				res_row.AccumulatedPremiumsAfterPolicyReduction = value.Value
			case "AnnualBasicInsuranceAmountAfterPolicyReduction":
				res_row.AnnualBasicInsuranceAmountAfterPolicyReduction = value.Value
			case "AnnualEffectiveInsuredAmountAfterPolicyReduction":
				res_row.AnnualEffectiveInsuredAmountAfterPolicyReduction = value.Value
			case "DeathOrTotalDisabilityInsuranceBenefittTheEndOfThePolicyYearAfterPolicyReduction":
				res_row.DeathOrTotalDisabilityInsuranceBenefittTheEndOfThePolicyYearAfterPolicyReduction = value.Value
			case "CashValueAtTheEndOfThePolicyYearAfterPolicyReduction":
				res_row.CashValueAtTheEndOfThePolicyYearAfterPolicyReduction = value.Value
			}

		}
		res = append(res, res_row)
	}


	tmpInsurancesInfo := &InsurancesInfo{
		ID: excel.ProductId,
		Info: res,
		CompanyName: excel.Company,
		ProductName: excel.Product,
	}

	op.ResChan <- tmpInsurancesInfo

	return nil

}