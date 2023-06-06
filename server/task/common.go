package task


import (
"github.com/game/server/logger/logging"
)

var logger = logging.GetLogger("handler", logging.DEFAULT_LEVEL)


type InsuranceInfo struct {
	Age       string   `json:"age"`
	Sex       string   `json:"sex"`
	Period    string   `json:"period"`
	Total     string   `json:"total"`
	ProductID []string `json:"product_id"`
}


//type InsuranceRes struct {
//	Age       string   `json:"age"`
//	Total     string   `json:"total"`
//	ProductID []string `json:"product_id"`
//}


type ExcelInsurance struct {
	PY float32 // 保单年度
	AATEOPY float32 // 保单年度末年龄
	PPITYBPR float32 // 减保前年交保险费
	APBPR float32 // 减保前累计保险费
	BSIFTYBPR float32 // 减保前年度基本保额
	ESIFTYBPR float32 // 减保前年度有效保额
	DOTDIBATEOTPYBPR float32 // 减保前保单年度末身故或全残保险金
	CVATEOTPYBPR float32 // 减保前保单年度末现金价值
	PYERAOEC float32 // 保单年度末减保金额（可选输入列）
	DAOBSIATEOTPY float32 // 保单年度末基本保额减少额度
	CRAATEOTPY float32 // 保单年度末累计减保金额
	APPAPR float32 // 减保后年交保险费
	APAPR float32 // 减保后累计保险费
	ABIAAPR float32 // 减保后年度基本保额
	AEIAAPR float32 // 减保后年度有效保额
	DOTDIBTEOTPYAPR float32 // 减保后保单年度末身故或全残保险金
	CVATEOTPYAPR float32 // 减保后保单年度末现金价值
}
//
//"保单年度": "PY",
//"保单年度末年龄": "AATEOPY",
//"减保前年交保险费": "PPITYBPR",
//"减保前累计保险费": " APBPR",
//"减保前年度基本保额": "BSIFTYBPR",
//"减保前年度有效保额": "ESIFTYBPR",
//"减保前保单年度末身故或全残保险金": "DOTDIBATEOTPYBPR",
//"减保前保单年度末现金价值": "CVATEOTPYBPR",
//"保单年度末减保金额（可选输入列）": "PYERAOEC",
//"保单年度末基本保额减少额度": "DAOBSIATEOTPY",
//"保单年度末累计减保金额": "CRAATEOTPY",
//"减保后年交保险费": "APPAPR",
//"减保后累计保险费": "APAPR",
//"减保后年度基本保额": "ABIAAPR",
//"减保后年度有效保额": "AEIAAPR",
//"减保后保单年度末身故或全残保险金": "DOTDIBTEOTPYAPR",
//"减保后保单年度末现金价值": "CVATEOTPYAPR",

// 保单年度
// 保单年度末年龄
// 减保前年交保险费
// 减保前累计保险费
// 减保前年度基本保额
// 减保前年度有效保额
// 减保前保单年度末身故或全残保险金
// 减保前保单年度末现金价值
// 保单年度末减保金额（可选输入列）
// 保单年度末基本保额减少额度
// 保单年度末累计减保金额
// 减保后年交保险费
// 减保后累计保险费
// 减保后年度基本保额
// 减保后年度有效保额
// 减保后保单年度末身故或全残保险金
// 减保后保单年度末现金价值

type InsurancesInfo struct {
	ID          string	`json:"id"`
	Info        []*BenefitDisplay `json:"info"`
	CompanyName string  `json:"company_name"`
	ProductName string	`json:"product_name"`
	Age         string  `json:"age"`
	Sex         string  `json:"sex"`
	Total       string  `json:"total"`
	Period      string  `json:"period"`
}

type InputConfig struct {
	TaskID            string `json:"task_id"`
	ExcelFilePath     string `json:"excel_file_path"`
	ExcelFileTempPath string `json:"excel_file_temp_path"`
	Set               *Set   `json:"set"`
	Get               *Get   `json:"get"`
}


type SetInfo struct {
	Set *Set `json:"set"`
}

type Set struct {
	Age   *SetBase `json:"age"`
	Total *SetBase `json:"total"`
	Sex   *SetBase `json:"sex"`
	Period *SetBase `json:"period"`
}

type SetBase struct {
	Value    string `json:"value"`
	Location string `json:"location"`
	Sheet    string `json:"sheet"`
}



type GetInfo struct {
	Get *Get `json:"get"`
}

type Get struct {
	BenefitDiff []*BenefitDiff `json:"benefit_diff"`
}

type BenefitDiff struct {
	Name string `json:"name"`
	Location string `json:"location"`
	Sheet    string `json:"sheet"`
}

// {
//     "task_id": "123",
//     "excel_file_path": "/Users/zhanggong02/python/xjw/tmp/333.xlsx",
//     "excel_file_temp_path": "/Users/zhanggong02/python/xjw/history/333_123.xlsx",
//     "set": {
//         "age": {
//             "value": "39",
//             "location": "F12",
//             "sheet": "资料输入"
//         },
//         "total": {
//             "value": "180000",
//             "location": "H15",
//             "sheet": "资料输入"
//         }
//     },
//     "get": {
//         "benefit_diff": {
//             "location": "A8:Q114",
//             "sheet": "利益演示"
//         }
//     }
// }

type InsurancesResFromPython struct {
	ErrorMessage string                `json:"error_message"`
	ExcelData    *InsurancesFromPython `json:"excel_data"`
}

type InsurancesFromPython struct {
	TaskID         string              `json:"task_id"`
	SetInfo        *SetConfig          `json:"set_config"`
	TaskInfo       *TaskInfo           `json:"task_info"`
	BenefitDisplay []*BenefitDisplay `json:"benefit_display"`
}

type SetConfig struct {
	FilePath string `json:"file_path"`
	Age      int    `json:"age"`
	Total    int    `json:"total"`
}

type TaskInfo struct {
	ExecuteTime     string `json:"execute_time"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	HistoryFilePath string `json:"history_file_path"`
}

type BenefitDisplay struct {
	Name  string  `json:"name"`
	Values []float32 `json:"values"`
	LableAge []string `json:"lable_age"`
	LableMutiple []string `json:"lable_mutiple"`
}