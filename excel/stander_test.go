package excel

import (
	"testing"
	"testing/quick"
)

func TestDefaultUnitTypePolicy(t *testing.T){
	v:=defaultUnitTypePolicy("123")
	if err:=quick.CheckEqual(v,1234,nil);err !=nil {
		t.Error(err)
	}
}
