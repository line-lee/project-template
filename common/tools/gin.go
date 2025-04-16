package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/project-template/common/models/variety/bo"
	"reflect"
)

func GinParamBind(ctx *gin.Context, param interface{}) error {
	err := ctx.ShouldBind(param)
	if err != nil && err.Error() != "EOF" {
		return err
	}

	value, exists := ctx.Get(bo.HolderInformation)
	if exists {
		vo := reflect.ValueOf(param)
		if field := vo.Elem().FieldByName(bo.HolderInformationX); field.CanSet() {
			field.SetString(value.(string))
		}
		if field := vo.Elem().FieldByName(bo.HolderInformation); field.CanSet() {
			field.SetString(value.(string))
		}
	}
	return nil
}
