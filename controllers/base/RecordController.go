//访问用户登录日志信息
package base

import (
	"encoding/json"
	md "goERP/models"
)

//列表视图列数-1，第一列为checkbox

type RecordController struct {
	BaseController
}

func (ctl *RecordController) Get() {
	ctl.PageName = "登陆记录管理"
	ctl.URL = "/record/"
	ctl.Data["URL"] = ctl.URL
	ctl.Data["MenuRecordActive"] = "active"
	ctl.GetList()
	ctl.Data["PageName"] = ctl.PageName + "\\" + ctl.PageAction

}
func (ctl *RecordController) Post() {
	action := ctl.Input().Get("action")
	switch action {
	case "table":
		ctl.PostList()
	case "one":
		ctl.GetOneRecord()
	default:
		ctl.PostList()
	}
}
func (ctl *RecordController) GetOneRecord() {

}
func (ctl *RecordController) PostList() {
	query := make(map[string]string)
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	if result, err := ctl.recordList(query, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}
func (ctl *RecordController) recordList(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var records []md.Record
	paginator, records, err := md.GetAllRecord(query, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		tableLines := make([]interface{}, 0, 4)
		for _, record := range records {
			oneLine := make(map[string]interface{})
			oneLine["Email"] = record.User.Email
			oneLine["Mobile"] = record.User.Mobile
			oneLine["Name"] = record.User.Name
			oneLine["NameZh"] = record.User.NameZh
			oneLine["UserAgent"] = record.UserAgent
			oneLine["CreateDate"] = record.CreateDate.Format("2006-01-02 15:04:05")
			oneLine["Logout"] = record.Logout.Format("2006-01-02 15:04:05")
			oneLine["Ip"] = record.Ip
			oneLine["Id"] = record.Id
			tableLines = append(tableLines, oneLine)
		}
		result["data"] = tableLines

		if jsonResult, er := json.Marshal(&paginator); er == nil {
			result["paginator"] = string(jsonResult)
			result["total"] = paginator.TotalCount
		}
	}
	return result, err
}
func (ctl *RecordController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-record"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "user/record_list_search.html"
}