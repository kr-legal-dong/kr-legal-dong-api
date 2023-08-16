package dong

import (
	"kr-legal-dong-api/db"
	"kr-legal-dong-api/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDong godoc
// @Summary 동 조회
// @Description 동 조회
// @Param request query dong.GetDong.request true "query params"
// @Success 200 {object} model.DataResponse{data=[]dong.GetDong.dong} "success"
// @Router /v1/dong [get]
func GetDong(c *gin.Context) {
	type request struct {
		SiCode  *string `form:"siCode"`
		GuCode  *string `form:"guCode"`
		Keyword *string `form:"keyword"`
		Active  *string `form:"active"`
	}

	type dong struct {
		Code   string `json:"code" binding:"required"`
		SiCode string `json:"siCode" binding:"required"`
		SiName string `json:"siName" binding:"required"`
		GuCode string `json:"guCode" binding:"required"`
		GuName string `json:"guName" binding:"required"`
		Name   string `json:"name" binding:"required"`
		Active bool   `json:"active" binding:"required"`
	}

	req := &request{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.DefaultResponse{
			Message: "bad_request",
		})
		c.Abort()
		return
	}

	where := []string{}
	value := []any{}

	if req.Keyword != nil {
		where = append(where, "do.name LIKE ?")
		value = append(value, "%"+*req.Keyword+"%")
	}

	if req.Active != nil {
		active, err := strconv.ParseBool(*req.Active)

		if err != nil {
			c.JSON(http.StatusBadRequest, &model.DefaultResponse{
				Message: "bad_request",
			})
			c.Abort()
			return
		}

		where = append(where, "do.active = ?")
		value = append(value, active)
	}

	if req.GuCode != nil {
		where = append(where, "gu.code = ?")
		value = append(value, req.GuCode)
	}

	if req.SiCode != nil {
		where = append(where, "si.code = ?")
		value = append(value, req.SiCode)
	}

	whereString := ""

	for index, data := range where {
		if index == 0 {
			whereString += " WHERE "
		} else {
			whereString += " AND "
		}
		whereString += data
	}

	rows, err := db.DB.Query("SELECT do.code, si.code, si.name, gu.code, gu.name, do.name, do.active FROM dong AS do INNER JOIN gu ON do.gu_code = gu.code INNER JOIN si ON gu.si_code = si.code"+whereString, value...)

	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.DefaultResponse{
			Message: "failed_select",
		})
		c.Abort()
		return
	}

	data := []dong{}

	for rows.Next() {
		var result dong
		rows.Scan(&result.Code, &result.SiCode, &result.SiName, &result.GuCode, &result.GuName, &result.Name, &result.Active)
		data = append(data, result)
	}

	c.JSON(http.StatusOK, &model.DataResponse{
		Message: "success",
		Data:    data,
	})
}
