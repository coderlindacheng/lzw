package source

import (
	. "gopkg.in/check.v1"
	"github.com/coderlindacheng/lzw/excel/common"
	"fmt"
	"github.com/coderlindacheng/lzw/excel/standar"
	"testing"
	"github.com/coderlindacheng/balabalago/types"
)

const TEST_INPUT_FILE_NAME = "../../原始表.xlsx"
const TEST_STANDAR_FILE_NAME = "../../评分标准.xlsx"
const TEST_OUTPUT_FILE_NAME = "../../分数表.xlsx"

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func init() {
	common.ReadSheet(TEST_STANDAR_FILE_NAME, standar.Read)
}

func (s *TestSuite) SetUpTest(c *C) {
	datasToOutput = make([]*types.Pair, 0, 20) //要重新初始化一下,要不然测试会出问题
}

func (s *TestSuite) TestReadSheet(c *C) {
	c.Assert(common.ReadSheet(TEST_INPUT_FILE_NAME, Read), IsNil)
	if len(datasToOutput) < 1 {
		c.FailNow()
	}
	for _, v := range datasToOutput {
		if v == nil {
			continue
		}
		fmt.Println(v)
	}
}

func (s *TestSuite) TestOutputSheet(c *C) {
	common.ReadSheet(TEST_INPUT_FILE_NAME, Read)
	c.Assert(Output(TEST_OUTPUT_FILE_NAME), IsNil)
}
