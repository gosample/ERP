package models

import (
	"errors"
	"fmt"
	"goERP/utils"
	"strings"

	"github.com/astaxie/beego/orm"
)

// 公司
type Company struct {
	Base
	Name       string           `orm:"unique" json:"name"`               //公司名称
	Children   []*Company       `orm:"reverse(many)" json:"childs"`      //子公司
	Parent     *Company         `orm:"rel(fk);null" json:"parent"`       //上级公司
	Department []*Department    `orm:"reverse(many)" json:"departments"` //部门
	Country    *AddressCountry  `orm:"rel(fk);null" json:"country"`      //国家
	Province   *AddressProvince `orm:"rel(fk);null" json:"province"`     //身份
	City       *AddressCity     `orm:"rel(fk);null" json:"city"`         //城市
	District   *AddressDistrict `orm:"rel(fk);null" json:"district"`     //区县
	Street     string           `orm:"default(\"\")" json:"street"`      //街道
}

func init() {
	orm.RegisterModel(new(Company))
}

// AddCompany insert a new Company into database and returns
// last inserted Id on success.
func AddCompany(obj *Company) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(obj)
	return id, err
}

// GetCompanyById retrieves Company by Id. Returns error if
// Id doesn't exist
func GetCompanyById(id int64) (obj *Company, err error) {
	o := orm.NewOrm()
	obj = &Company{Base: Base{Id: id}}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetCompanyByName retrieves Company by Name. Returns error if
// Name doesn't exist
func GetCompanyByName(name string) (obj *Company, err error) {
	o := orm.NewOrm()
	obj = &Company{Name: name}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetAllCompany retrieves all Company matches certain condition. Returns empty list if
// no records exist
func GetAllCompany(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (utils.Paginator, []Company, error) {
	var (
		objArrs   []Company
		paginator utils.Paginator
		num       int64
		err       error
	)
	o := orm.NewOrm()
	qs := o.QueryTable(new(Company))
	qs = qs.RelatedSel()

	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return paginator, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return paginator, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return paginator, nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return paginator, nil, errors.New("Error: unused 'order' fields")
		}
	}

	qs = qs.OrderBy(sortFields...)
	if cnt, err := qs.Count(); err == nil {
		paginator = utils.GenPaginator(limit, offset, cnt)
	}
	if num, err = qs.Limit(limit, offset).All(&objArrs, fields...); err == nil {
		paginator.CurrentPageSize = num
	}
	return paginator, objArrs, err
}

// UpdateCompany updates Company by Id and returns error if
// the record to be updated doesn't exist
func UpdateCompanyById(m *Company) (err error) {
	o := orm.NewOrm()
	v := Company{Base: Base{Id: m.Id}}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCompany deletes Company by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCompany(id int64) (err error) {
	o := orm.NewOrm()
	v := Company{Base: Base{Id: id}}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Company{Base: Base{Id: id}}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}