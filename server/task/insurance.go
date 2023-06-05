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
	tempPath := "D:\\app\\gin-backend\\history\\" + excel.ProductId + "_" + excel.ExcelName + ".xlsx"
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
				Value:    setInfo.Set.Age.Value,
				Location: setInfo.Set.Age.Location,
				Sheet:    setInfo.Set.Age.Sheet,
			},
			Total: &SetBase{
				Value:    setInfo.Set.Total.Value,
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
			if value.Name == "保单年度" {
				fmt.Printf("**************==================************** = %s\n", value.Name)
			}
			switch value.Name {
			case "PY":
				res_row.PY = value.Value
			case "AATEOPY":
				res_row.AATEOPY = value.Value
			case "PPITYBPR":
				res_row.PPITYBPR = value.Value
			case "APBPR":
				res_row.APBPR = value.Value
			case "BSIFTYBPR":
				res_row.BSIFTYBPR = value.Value
			case "ESIFTYBPR":
				res_row.ESIFTYBPR = value.Value
			case "DOTDIBATEOTPYBPR":
				res_row.DOTDIBATEOTPYBPR = value.Value
			case "CVATEOTPYBPR":
				res_row.CVATEOTPYBPR = value.Value
			case "PYERAOEC":
				res_row.PYERAOEC = value.Value
			case "DAOBSIATEOTPY":
				res_row.DAOBSIATEOTPY = value.Value
			case "CRAATEOTPY":
				res_row.CRAATEOTPY = value.Value
			case "APPAPR":
				res_row.APPAPR = value.Value
			case "APAPR":
				res_row.APAPR = value.Value
			case "ABIAAPR":
				res_row.ABIAAPR = value.Value
			case "AEIAAPR":
				res_row.AEIAAPR = value.Value
			case "DOTDIBTEOTPYAPR":
				res_row.DOTDIBTEOTPYAPR = value.Value
			case "CVATEOTPYAPR":
				res_row.CVATEOTPYAPR = value.Value
			}

		}
		res = append(res, res_row)
	}


	tmpInsurancesInfo := &InsurancesInfo{
		ID: excel.ProductId,
		Info: res,
		CompanyName: excel.ExcelName,
		ProductName: excel.ExcelName,
	}

	op.ResChan <- tmpInsurancesInfo

	return nil

}



func InsuranceSerial(c * gin.Context, id, age, total string) (*InsurancesInfo, error){
	logger.Infof(fmt.Sprintf("Insurance serrial start"))

	//id := op.Get().(string)
	logger.Infof(fmt.Sprintf("InsuranceOp id = %s ", id))

	excel,err := db.GetProductExcelParserByID(id)
	if err != nil {
		common.ErrorResponse(c, common.E_GetFailed, err.Error())
		return nil, err
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

		common.ErrorResponse(c, common.E_JsonUnmarshalError, err.Error())
		return nil, err
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
		common.ErrorResponse(c, common.E_JsonUnmarshalError, err.Error())
		return nil, err
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
				Value:    age,
				Location: setInfo.Set.Age.Location,
				Sheet:    setInfo.Set.Age.Sheet,
			},
			Total: &SetBase{
				Value:    total,
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
		common.ErrorResponse(c, "1111", err.Error())

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
			fmt.Printf("for _, value := range value================== = %+v\n", value.Value)

			if value.Name == "保单年度" {
				fmt.Printf("**************==================************** = %s\n", value.Name)
			}
			switch value.Name {
			case "PY":
				res_row.PY = value.Value
			case "AATEOPY":
				res_row.AATEOPY = value.Value
			case "PPITYBPR":
				res_row.PPITYBPR = value.Value
			case "APBPR":
				res_row.APBPR = value.Value
			case "BSIFTYBPR":
				res_row.BSIFTYBPR = value.Value
			case "ESIFTYBPR":
				res_row.ESIFTYBPR = value.Value
			case "DOTDIBATEOTPYBPR":
				res_row.DOTDIBATEOTPYBPR = value.Value
			case "CVATEOTPYBPR":
				res_row.CVATEOTPYBPR = value.Value
			case "PYERAOEC":
				res_row.PYERAOEC = value.Value
			case "DAOBSIATEOTPY":
				res_row.DAOBSIATEOTPY = value.Value
			case "CRAATEOTPY":
				res_row.CRAATEOTPY = value.Value
			case "APPAPR":
				res_row.APPAPR = value.Value
			case "APAPR":
				res_row.APAPR = value.Value
			case "ABIAAPR":
				res_row.ABIAAPR = value.Value
			case "AEIAAPR":
				res_row.AEIAAPR = value.Value
			case "DOTDIBTEOTPYAPR":
				res_row.DOTDIBTEOTPYAPR = value.Value
			case "CVATEOTPYAPR":
				res_row.CVATEOTPYAPR = value.Value
			}

		}
		res = append(res, res_row)
	}

for _, v := range res {
	//logger.Infof(fmt.Sprintf("res.v.PolicyYear: %f", v.PolicyYear))
	logger.Infof(fmt.Sprintf("res.v : %+v", v))
}


	tmpInsurancesInfo := &InsurancesInfo{
		ID: excel.ProductId,
		Info: res,
		CompanyName: excel.ExcelName,
		ProductName: excel.ExcelName,
	}

	//op.ResChan <- tmpInsurancesInfo

	return tmpInsurancesInfo, nil

}