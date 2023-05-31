package task

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"testing"
)

func TestInsurance(t *testing.T) {

	unmarshalSetInfo := `{"set": {"age": {"value": "39","location": "F12","sheet": "资料输入"},"total": {"value": "180000","location": "H15","sheet": "资料输入"}}}`
	logger.Infof(fmt.Sprintf("&unmarshalSetInfo = %s",unmarshalSetInfo))

	setInfo := new(SetInfo)
	err := json.Unmarshal([]byte(unmarshalSetInfo), &setInfo)
	if err != nil {
		logger.Infof(fmt.Sprintf("err = json.Unmarshal([]byte(unmarshalSetInfo), &setInfo)"))

		//common.ErrorResponse(op.Ctx, common.E_JsonUnmarshalError, err.Error())
		//return err
	}

	logger.Infof(fmt.Sprintf("setInfo = %+v", setInfo))
	logger.Infof(fmt.Sprintf("setInfo.Age.Value = %s", setInfo.Set.Age.Value))
	logger.Infof(fmt.Sprintf("setInfo.Age.Sheet = %s", setInfo.Set.Age.Sheet))
	logger.Infof(fmt.Sprintf("setInfo.Age.Location = %s", setInfo.Set.Age.Location))


	unmarshalGetInfo := `{"get": {"benefit_diff": {"location": "A8:Q114","sheet": "利益演示"}}}`
	logger.Infof(fmt.Sprintf("&unmarshalGetInfo = %s",unmarshalGetInfo))

	getInfo := GetInfo{
		 Get: &Get{
			BenefitDiff: &BenefitDiff{},
			},
	}
	err = json.Unmarshal([]byte(unmarshalGetInfo), &getInfo)
	if err != nil {
		logger.Infof(fmt.Sprintf("err = json.Unmarshal([]byte(unmarshalGetInfo), &getInfo) err = %s", err.Error()))

		//common.ErrorResponse(op.Ctx, common.E_JsonUnmarshalError, err.Error())
		//return err
	}

	logger.Infof(fmt.Sprintf("getInfo = %+v", getInfo))
	logger.Infof(fmt.Sprintf("getInfo.BenefitDiff.Location = %s", getInfo.Get.BenefitDiff.Location))
	logger.Infof(fmt.Sprintf("getInfo.BenefitDiff.Sheet = %s", getInfo.Get.BenefitDiff.Sheet))

	uid, err := uuid.GenerateRandomBytes(8)
	uuuid,err := uuid.GenerateUUID()
	logger.Infof(fmt.Sprintf("uid= %s", uid))
	logger.Infof(fmt.Sprintf("uid= %s", uuuid))


}