package task

import (
	"encoding/json"
	"fmt"
	"github.com/game/server/common"
	"github.com/game/server/db"
	"github.com/game/server/utils/linux"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"os"
	"strings"
	"time"
)

type InsuranceOp struct {
	*TaskQueue
	Ctx   *gin.Context
	InsuranceInfo InsuranceInfo
}

func NewInsuranceOp(t *TaskQueue, c *gin.Context, insuranceInfo InsuranceInfo) *InsuranceOp{
	return &InsuranceOp{
		TaskQueue: t,
		Ctx: c,
		InsuranceInfo:    insuranceInfo,
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


	setInfo := new(Set)
	err = json.Unmarshal([]byte(excel.SetInfo), &setInfo)
	if err != nil {
		common.ErrorResponse(op.Ctx, common.E_JsonUnmarshalError, err.Error())
		return err
	}

	logger.Infof(fmt.Sprintf("setInfo = %+v", setInfo))
	logger.Infof(fmt.Sprintf("setInfo.Age.Value = %s", setInfo.Age.Value))
	logger.Infof(fmt.Sprintf("setInfo.Age.Sheet = %s", setInfo.Age.Sheet))
	logger.Infof(fmt.Sprintf("setInfo.Age.Location = %s", setInfo.Age.Location))


	getInfo := new(Get)
	err = json.Unmarshal([]byte(excel.GetInfo), &getInfo)
	if err != nil {
		common.ErrorResponse(op.Ctx, common.E_JsonUnmarshalError, err.Error())
		return err
	}

	logger.Infof(fmt.Sprintf("getInfo = %+v", getInfo))
	logger.Infof(fmt.Sprintf("getInfo.BenefitDiff.Sheet = %s", getInfo.BenefitDiff.Sheet))
	logger.Infof(fmt.Sprintf("getInfo.BenefitDiff.Location = %s", getInfo.BenefitDiff.Location))


	uid, err := uuid.GenerateRandomBytes(8)

	tempPath := strings.TrimRight(excel.ExcelPath, ".xlsx")
	tempPath = tempPath + "_" + string(uid) + ".xlsx"
	logger.Infof(fmt.Sprintf("tempPath = %s", tempPath))

	config := &InputConfig{
		TaskID:            string(id),
		ExcelFilePath:     excel.ExcelPath,
		ExcelFileTempPath: tempPath,
		Set: &Set{
			Age: &SetInfo{
				Value:    setInfo.Age.Value,
				Location: setInfo.Age.Location,
				Sheet:    setInfo.Age.Sheet,
			},
			Total: &SetInfo{
				Value:    setInfo.Total.Location,
				Location: setInfo.Total.Location,
				Sheet:    setInfo.Total.Sheet,
			},
		},
		Get: &Get{
			BenefitDiff: &BenefitDiff{
				Location: getInfo.BenefitDiff.Location,
				Sheet:    getInfo.BenefitDiff.Sheet,
			},
		},
	}

	configByte, err := json.Marshal(config)
	if err != nil {
		logger.Infof(fmt.Sprintf("json.Marshal(config) failed = %s", err.Error()))

		//logger.Infof("json.Marshal(config) failed")
	}

	configName := string(uid) + ".conf"
	logger.Infof(fmt.Sprintf("configName = %s", configName))

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

	configPath := "D:/python/xjw/" + configName
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

	fmt.Printf("execOut ====================== = %+v\n", execOut)

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
			fmt.Printf("for _, value := range v================== = %+v\n", value.Name)
			fmt.Printf("for _, value := range value================== = %+v\n", value.Value)

			if value.Name == "保单年度" {
				fmt.Printf("**************==================************** = %+v\n", value.Name)
			}
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

	op.ResChan <- res

	return nil

}