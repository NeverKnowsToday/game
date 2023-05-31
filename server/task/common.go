package task


import (
"github.com/game/server/logger/logging"
)

var logger = logging.GetLogger("handler", logging.DEFAULT_LEVEL)


type InsuranceInfo struct {
	Age       string   `json:"age"`
	Total     string   `json:"total"`
	ProductID []string `json:"product_id"`
}


type InsuranceRes struct {
	Age       string   `json:"age"`
	Total     string   `json:"total"`
	ProductID []string `json:"product_id"`
}


type ExcelInsurance struct {
	PolicyYear                                                                          float32 // 保单年度
	AgeAtTheEndOfPolicyYear                                                             float32 // 保单年度末年龄
	PremiumPaymentInTheYearBeforePolicyReduction                                        float32 // 减保前年交保险费
	AccumulatedPremiumsBeforePolicyReduction                                            float32 // 减保前累计保险费
	BasicSumInsuredForTheYearBeforePolicyReduction                                      float32 // 减保前年度基本保额
	EffectiveSumInsuredForTheYearBeforePolicyReduction                                  float32 // 减保前年度有效保额
	DeathOrTotalDisabilityInsuranceBenefitsAtTheEndOfThePolicyYearBeforePolicyReduction float32 // 减保前保单年度末身故或全残保险金
	CashValueAtTheEndOfThePolicyYearBeforePolicyReduction                               float32 // 减保前保单年度末现金价值
	PolicyYearEndReductionAmountOptionalEntryColumn                                     float32 // 保单年度末减保金额（可选输入列）
	DecreaseAmountOfBasicSumInsuredAtTheEndOfThePolicyYear                              float32 // 保单年度末基本保额减少额度
	CumulativeReductionAmountAtTheEndOfThePolicyYear                                    float32 // 保单年度末累计减保金额
	AnnualPremiumPaymentAfterPolicyReduction                                            float32 // 减保后年交保险费
	AccumulatedPremiumsAfterPolicyReduction                                             float32 // 减保后累计保险费
	AnnualBasicInsuranceAmountAfterPolicyReduction                                      float32 // 减保后年度基本保额
	AnnualEffectiveInsuredAmountAfterPolicyReduction                                    float32 // 减保后年度有效保额
	DeathOrTotalDisabilityInsuranceBenefittTheEndOfThePolicyYearAfterPolicyReduction    float32 // 减保后保单年度末身故或全残保险金
	CashValueAtTheEndOfThePolicyYearAfterPolicyReduction                                float32 // 减保后保单年度末现金价值
}

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
	ID          string
	Info        []*ExcelInsurance
	CompanyName string
	ProductName string
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
	BenefitDiff *BenefitDiff `json:"benefit_diff"`
}

type BenefitDiff struct {
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
	BenefitDisplay [][]*BenefitDisplay `json:"benefit_display"`
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
	Value float32 `json:"value"`
}