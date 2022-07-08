// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-08
// Based on aurservd by liasica, magicrolan@qq.com.

package controller

// FJJ01
// @ID           FJJ01
// @Router       /crm/customer [POST]
// @Summary      F01添加客户
// @Tags         [F]房金聚接口
// @Accept       json
// @Produce      json
// @Param        body  body     CrmCustomerCreate  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func FJJ01() {

}

type CrmCustomerCreate struct {
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
