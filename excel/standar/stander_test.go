package standar

import (
	"testing"
	"github.com/coderlindacheng/balabalago/time"
	"github.com/coderlindacheng/balabalago/special_string"
)

func checkUnitTypePolicySuccessed(testStr *string, wanted *int, t *testing.T) {
	result, err := timeUnitTypePolicy(*testStr)
	if err != nil {
		t.Fatal(err)
	}
	if *wanted != result {
		t.Fatalf("想要的结果是 %s 得到的结果是 %s", wanted, result)
	}
}

/*
	成功测试
 */
func TestTimeUnitTypePolicySuccessed(t *testing.T) {

	checkUnitTypePolicySuccessed(&"1", &time.MILLIS_PER_MINUTE, t)

	checkUnitTypePolicySuccessed(&"1\"", &time.MILLIS_PER_MINUTE, t)

	checkUnitTypePolicySuccessed(&"1\"30'", &(time.MILLIS_PER_MINUTE + time.MILLIS_PER_SECOND*30), t)

	checkUnitTypePolicySuccessed(&"30'3", &(time.MILLIS_PER_SECOND*30 + 3*100), t)

	checkUnitTypePolicySuccessed(&"30'", &(time.MILLIS_PER_SECOND * 30), t)

	checkUnitTypePolicySuccessed(&"1\"30'2", &(time.MILLIS_PER_MINUTE + time.MILLIS_PER_SECOND*30 + 2*100), t)
}

func checkUnitTypePolicyFailed(testStr *string, t *testing.T) {
	if _, err := timeUnitTypePolicy(*testStr); err == nil {
		t.Fatalf("%s 这个字符格式是必定不能通过测试的", *testStr)
	}
}

/*
	失败测试
 */
func TestTimeUnitTypePolicyFailed(t *testing.T) {
	checkUnitTypePolicyFailed(&special_string.SINGLE_QUOTE, t)
	checkUnitTypePolicyFailed(&special_string.QUOTE, t)
	checkUnitTypePolicyFailed(&(special_string.QUOTE + special_string.SINGLE_QUOTE), t)
	checkUnitTypePolicyFailed(&(special_string.SINGLE_QUOTE + special_string.QUOTE), t)
	checkUnitTypePolicyFailed(&"1'30'2\"", t)
}

func TestParseRowName(t *testing.T) {
	rowName := "长跑#男#时间"
	name, sex, policy, err := parseRowName(rowName)
	if err != nil {
		t.Fatal(err)
	}
	if name != "长跑" || sex != "男" {
		t.Fatal("组名解析错误")
	}

	testStr := "1\"30'2"
	wanted := time.MILLIS_PER_MINUTE + time.MILLIS_PER_SECOND*30 + 2*100
	result, err := policy(testStr);
	if err != nil {
		t.Fatal(err)
	}
	if wanted != result {
		t.Fatalf("想要的结果是 %s 得到的结果是 %s", wanted, result)
	}
}
