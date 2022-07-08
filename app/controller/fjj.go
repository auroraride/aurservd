// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-08
// Based on aurservd by liasica, magicrolan@qq.com.

package controller

// FJJ01
// @ID           FJJ01
// @Router       /crm/customer [POST]
// @Summary      F01 添加客户
// @Tags         [F]房金聚接口
// @Accept       json
// @Produce      json
// @Param        body  body     CrmCustomer  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func FJJ01() {}

// FJJ02
// @ID           ManagerFjjFJJ02
// @Router       /manager/crm/customer [GET]
// @Summary      F02 获取客户列表「需要认证」
// @Tags         [F]房金聚接口
// @Accept       json
// @Produce      json
// @Param        body  body     CrmCustomerListReq  false  "筛选条件"
// @Success      200  {object}  model.StatusResponse{items=[]CrmCustomerWithID}  "请求成功"
func FJJ02() {}

// FJJ03
// @ID           ManagerFjjFJJ03
// @Router       /manager/crm/customer/{id} [PUT]
// @Summary      F03 修改客户
// @Tags         [F]房金聚接口
// @Accept       json
// @Produce      json
// @Param        body  body     CrmCustomerWithID  true  "修改数据"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func FJJ03() {}

type CrmCustomerWithID struct {
    ID uint64 `json:"id" param:"id" validate:"required" trans:"客户ID"`
    CrmCustomer
}

type CrmCustomerListReq struct {
    Page int `json:"page" query:"page"` // 分页数据, 从0开始

    Name     string `json:"name" query:"name"`         // 姓名
    Phone    string `json:"phone" query:"phone"`       // 手机号
    Status   string `json:"status" query:"status"`     // [选择]状态: 交件 / 签约 / 审批 / 放款 / 被拒
    Sn       string `json:"sn" query:"sn"`             // 客户编号
    Type     string `json:"type" query:"type"`         // [选择]贷款类型: 房车贷 / 房产保贷 / 抵押贷 / 优客贷 / 随心智贷 / 中银消费 / 法拍按揭
    Channel  string `json:"channel" query:"channel"`   // 渠道
    Employee string `json:"employee" query:"employee"` // 业务员
    Branch   string `json:"branch" query:"branch"`     // 办理分行
    Manager  string `json:"manager" query:"manager"`   // 客户经理
}

type CrmCustomer struct {
    Name           *string  `json:"name" validate:"required"`            // 姓名
    Gender         *uint8   `json:"gender" validate:"required"`          // 性别 1:男 2:女
    Phone          *string  `json:"phone" validate:"required"`           // 手机号
    IDCardNumber   *string  `json:"id_card_number" validate:"required"`  // 身份证号
    DistrictID     *uint64  `json:"district_id" validate:"required"`     // 区域ID
    Address        *string  `json:"address" validate:"required"`         // 常住地址
    Type           *string  `json:"type" validate:"required"`            // 贷款类型: 房车贷 / 房产保贷 / 抵押贷 / 优客贷 / 随心智贷 / 中银消费 / 法拍按揭
    Employee       *string  `json:"employee" validate:"required"`        // 业务员
    Amount         *float64 `json:"amount" validate:"required"`          // 贷款需求(万元)
    Channel        *string  `json:"channel" validate:"required"`         // 渠道
    ChannelPhone   *string  `json:"channel_phone" validate:"required"`   // 渠道手机号
    ChannelContact *string  `json:"channel_contact" validate:"required"` // 渠道联系人

    DeliveryDate      *string   `json:"delivery_date,omitempty"`      // 交件日期
    Branch            *string   `json:"branch,omitempty"`             // 办理分行
    FaceDate          *string   `json:"face_date,omitempty"`          // 面签日期
    HouseInfo         *string   `json:"house_info,omitempty"`         // 房屋信息
    HousePrice        *string   `json:"house_price,omitempty"`        // 房屋估价(万元)
    Credit            *string   `json:"credit,omitempty"`             // 征信情况: 正常 / 连三累六 / 当月逾期 / 五类观察 / 其他
    SecondaryMortgage *bool     `json:"secondary_mortgage,omitempty"` // 是否二抵
    Income            *float64  `json:"income,omitempty"`             // 年收入
    IDCardPhotos      *[]string `json:"id_card_photos,omitempty"`     // 身份证照片
    IndustryPhotos    *[]string `json:"industry_photos,omitempty"`    // 产调表
    DrivingPhotos     *[]string `json:"driving_photos,omitempty"`     // 行驶证

    Status     *string  `json:"status,omitempty"`      // 状态: 交件 / 签约 / 审批 / 放款 / 被拒
    Manager    *string  `json:"manager,omitempty"`     // 客户经理
    LoanAmount *float64 `json:"loan_amount,omitempty"` // 审批额度(万元)
    LoanYears  *int     `json:"loan_years,omitempty"`  // 贷款年限
    LoanDate   *string  `json:"loan_date,omitempty"`   // 放款日期
    RealAmount *float64 `json:"real_amount,omitempty"` // 实际放款(万元)
}
