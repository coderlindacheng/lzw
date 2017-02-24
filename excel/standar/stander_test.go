package standar

import (
	"testing"
	. "github.com/coderlindacheng/balabalago/time"
	. "github.com/coderlindacheng/balabalago/special_string"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

/*
	成功测试
 */
func (s *MySuite) TestTimeUnitTypePolicySucceed(c *C) {
	processor := func(testStr string, wanted int) {
		result, err := timeUnitTypePolicy(testStr)
		c.Assert(err, IsNil)
		c.Assert(result, Equals, wanted)
	}
	processor("1", MILLIS_PER_MINUTE)

	processor("1\"", MILLIS_PER_MINUTE)

	processor("1\"30'", MILLIS_PER_MINUTE+MILLIS_PER_SECOND*30)

	processor("30'3", MILLIS_PER_SECOND*30+3*100)

	processor("30'", MILLIS_PER_SECOND*30)

	processor("1\"30'2", MILLIS_PER_MINUTE+MILLIS_PER_SECOND*30+2*100)
}

/*
	失败测试
 */
func (s *MySuite) TestTimeUnitTypePolicyFailed(c *C) {
	processor := func(s string) {
		_, err := timeUnitTypePolicy(s)
		c.Assert(err, NotNil)
		var e error
		c.Assert(err, Implements, &e)
	}
	processor(SINGLE_QUOTE)
	processor(QUOTE)
	processor(QUOTE + SINGLE_QUOTE)
	processor(SINGLE_QUOTE + QUOTE)
	processor("1'30'2\"")
}

func (s *MySuite) TestParseRowNameSucceed(c *C) {
	rowName := "长跑#男#时间"
	name, sex, policy, err := parseRowName(rowName)
	c.Assert(err, IsNil)
	c.Assert(name, Equals, "长跑")
	c.Assert(sex, Equals, "男")

	testStr := "1\"30'2"
	wanted := MILLIS_PER_MINUTE + MILLIS_PER_SECOND*30 + 2*100
	result, err := policy(testStr)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, wanted)

}

func (s *MySuite) TestParseRowNameFailed(c *C) {
	processor := func(s string) {
		_, _, _, err := parseRowName(s)
		c.Assert(err, NotNil)
	}
	processor("长跑##时间")
	processor("#长跑时间")
	processor("#长跑#时间")
	processor("#长跑#时间#")
	processor("长跑时间#")
}
