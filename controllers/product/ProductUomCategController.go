package product

import (
	"encoding/json"
	"goERP/controllers/base"
	md "goERP/models"
	"strconv"
	"strings"
)

type ProductUomCategController struct {
	base.BaseController
}

func (ctl *ProductUomCategController) Post() {
	action := ctl.Input().Get("action")
	switch action {
	case "validator":
		ctl.Validator()
	case "table": //bootstrap table的post请求
		ctl.PostList()
	case "create":
		ctl.PostCreate()
	default:
		ctl.PostList()
	}
}
func (ctl *ProductUomCategController) Get() {
	ctl.PageName = "单位类别管理"
	action := ctl.Input().Get("action")
	switch action {
	case "create":
		ctl.Create()
	case "edit":
		ctl.Edit()
	case "detail":
		ctl.Detail()
	default:
		ctl.GetList()
	}
	ctl.Data["PageName"] = ctl.PageName + "\\" + ctl.PageAction
	ctl.URL = "/product/uomcateg/"
	ctl.Data["URL"] = ctl.URL
	ctl.Data["MenuProductUomCategActive"] = "active"
}
func (ctl *ProductUomCategController) Put() {
	id := ctl.Ctx.Input.Param(":id")
	ctl.URL = "/product/uomcateg/"
	if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
		if uomCateg, err := md.GetProductUomCategById(idInt64); err == nil {
			if err := ctl.ParseForm(&uomCateg); err == nil {

				if err := md.UpdateProductUomCategById(uomCateg); err == nil {
					ctl.Redirect(ctl.URL+id+"?action=detail", 302)
				}
			}
		}
	}
	ctl.Redirect(ctl.URL+id+"?action=edit", 302)
}
func (ctl *ProductUomCategController) Validator() {
	name := ctl.GetString("name")
	recordID, _ := ctl.GetInt64("recordId")
	name = strings.TrimSpace(name)
	result := make(map[string]bool)
	obj, err := md.GetProductUomCategByName(name)
	if err != nil {
		result["valid"] = true
	} else {
		if obj.Name == name {
			if recordID == obj.Id {
				result["valid"] = true
			} else {
				result["valid"] = false
			}

		} else {
			result["valid"] = true
		}

	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}
func (ctl *ProductUomCategController) PostCreate() {
	uom := new(md.ProductUomCateg)
	if err := ctl.ParseForm(uom); err == nil {
		if id, err := md.AddProductUomCateg(uom); err == nil {
			ctl.Redirect("/product/uomcateg/"+strconv.FormatInt(id, 10)+"?action=detail", 302)
		} else {
			ctl.Get()
		}
	} else {
		ctl.Get()
	}
}
func (ctl *ProductUomCategController) productUomCategList(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {
	var arrs []md.ProductUomCateg
	paginator, arrs, err := md.GetAllProductUomCateg(query, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		// result["recordsFiltered"] = paginator.TotalCount
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["name"] = line.Name
			oneLine["Id"] = line.Id
			oneLine["id"] = line.Id
			uoms := line.Uoms
			mapValues := make(map[int64]string)
			for _, line := range uoms {
				mapValues[line.Id] = line.Name
			}
			oneLine["uoms"] = mapValues
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
func (ctl *ProductUomCategController) PostList() {
	query := make(map[string]string)
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	if result, err := ctl.productUomCategList(query, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()
}
func (ctl *ProductUomCategController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	categInfo := make(map[string]interface{})
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if categ, err := md.GetProductUomCategById(idInt64); err == nil {
				ctl.PageAction = categ.Name
				categInfo["name"] = categ.Name
			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordId"] = id
	ctl.Data["UomCateg"] = categInfo
	ctl.Layout = "base/base.html"

	ctl.TplName = "product/product_uom_categ_form.html"
}
func (ctl *ProductUomCategController) Detail() {
	//获取信息一样，直接调用Edit
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["Action"] = "detail"
}
func (ctl *ProductUomCategController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-product-uom-categ"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "product/product_uom_categ_list_search.html"
}
func (ctl *ProductUomCategController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.Layout = "base/base.html"
	ctl.PageAction = "创建"
	ctl.TplName = "product/product_uom_categ_form.html"
}