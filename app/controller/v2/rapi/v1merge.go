package rapi

type assistance struct{}

// Breakdown
// @ID		AssistanceBreakdown
// @Router	/rider/v2/assistance/breakdown [GET]
// @Summary	R5001 获取救援原因
// @Tags	Assistance - 救援
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string		true	"骑手校验token"
// @Success	200				{object}	[]string	"请求成功"
func (*assistance) Breakdown() {}

// Create
// @ID		AssistanceCreate
// @Router	/rider/v2/assistance [POST]
// @Summary	R5002 发起救援
// @Tags	Assistance - 救援
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.AssistanceCreateReq	true	"救援参数"
// @Success	200				{object}	model.AssistanceCreateRes	"请求成功"
func (*assistance) Create() {}

// Cancel
// @ID		AssistanceCancel
// @Router	/rider/v2/assistance/cancel [POST]
// @Summary	R5003 取消救援
// @Tags	Assistance - 救援
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.AssistanceCancelReq	true	"取消请求"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*assistance) Cancel() {}

// Current
// @ID		AssistanceCurrent
// @Router	/rider/v2/assistance/current [GET]
// @Summary	R5004 当前救援
// @Tags	Assistance - 救援
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Success	200				{object}	model.AssistanceSocketMessage	"救援信息, 救援不存在的时候返回data为null"
func (*assistance) Current() {}

// List
// @ID		AssistanceList
// @Router	/rider/v2/assistance [GET]
// @Summary	R5005 救援列表
// @Tags	Assistance - 救援
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string														true	"骑手校验token"
// @Param	query			query		model.PaginationReq											true	"分页参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.AssistanceSimpleListRes}	"请求成功"
func (*assistance) List() {}

type battery struct{}

// Detail
// @ID		BatteryDetail
// @Router	/rider/v2/battery [GET]
// @Summary	RA001 获取电池详情
// @Tags	Battery - 电池
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string				true	"骑手校验token"
// @Success	200				{object}	model.BatteryDetail	"请求成功"
func (*battery) Detail() {}

type branch struct{}

// List
// @ID		BranchList
// @Router	/rider/v2/branch [GET]
// @Summary	R2001 列举网点
// @Tags	Branch - 网点
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	query			query		model.BranchWithDistanceReq		true	"根据距离获取网点请求参数"
// @Success	200				{object}	[]model.BranchWithDistanceRes	"请求成功"
func (*branch) List() {}

// Riding
// @ID		BranchRiding
// @Router	/rider/v2/branch/riding [GET]
// @Summary	R2002 网点骑行规划时间
// @Tags	Branch - 网点
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	query			query		model.BranchRidingReq	true	"desc"
// @Success	200				{object}	model.BranchRidingRes	"请求成功"
func (*branch) Riding() {}

// Facility
// @ID		BranchFacility
// @Router	/rider/v2/branch/facility/{fid} [GET]
// @Summary	R2004 设施详情
// @Tags	Branch - 网点
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	fid				path		string					true	"设置标识"
// @Param	lng				query		float64					true	"经度"
// @Param	lat				query		float64					true	"纬度"
// @Success	200				{object}	model.BranchFacilityRes	"请求成功"
func (*branch) Facility() {}

type business struct{}

// Active
// @ID		BusinessActive
// @Router	/rider/v2/business/active [POST]
// @Summary	R7001 激活骑士卡
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.BusinessCabinetReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatus	"请求成功"
func (*business) Active() {}

// Unsubscribe
// @ID		BusinessUnsubscribe
// @Router	/rider/v2/business/unsubscribe [POST]
// @Summary	R7002 退租
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.BusinessCabinetReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatus	"请求成功"
func (*business) Unsubscribe() {}

// Pause
// @ID		BusinessPause
// @Router	/rider/v2/business/pause [POST]
// @Summary	R7003 寄存
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.BusinessCabinetReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatus	"请求成功"
func (*business) Pause() {}

// Continue
// @ID		BusinessContinue
// @Router	/rider/v2/business/continue [POST]
// @Summary	R7004 取消寄存
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.BusinessCabinetReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatus	"请求成功"
func (*business) Continue() {}

// Status
// @ID		BusinessStatus
// @Router	/rider/v2/business/status [GET]
// @Summary	R7005 业务状态
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	query			query		model.BusinessCabinetStatusReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatusRes	"请求成功"
func (*business) Status() {}

// PauseInfo
// @ID		BusinessPauseInfo
// @Router	/rider/v2/business/pause/info [GET]
// @Summary	R7006 寄存信息
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	model.BusinessPauseInfoRes	"请求成功"
func (*business) PauseInfo() {}

// Allocated
// @ID			BusinessAllocated
// @Router		/rider/v2/business/allocated/{id} [GET]
// @Summary		R7009 长连接轮询是否已分配
// @Description	用以判定待激活骑士卡是否需要签约 (allocated = true)
// @Tags		Business - 业务
// @Accept		json
// @Produce		json
// @Param		X-Rider-Token	header		string					true	"骑手校验token"
// @Param		id				path		uint64					true	"订阅ID"
// @Success		200				{object}	model.AllocateRiderRes	"请求成功"
func (*business) Allocated() {}

// SubscribeSigned
// @ID		BusinessSubscribeSigned
// @Router	/rider/v2/business/subscribe/signed/{id} [GET]
// @Summary	R7010 长连接轮询是否已签约
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	id				path		uint64					true	"订阅ID"
// @Success	200				{object}	model.SubscribeSigned	"请求成功"
func (*business) SubscribeSigned() {}

// GetProcess
// @ID		CabinetGetProcess
// @Router	/rider/v2/cabinet/process/{serial} [GET]
// @Summary	R4001 获取换电信息
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	serial			path		string					true	"电柜二维码"
// @Success	200				{object}	model.RiderExchangeInfo	"请求成功"
func (*cabinet) GetProcess() {}

// Process
// @ID		CabinetProcess
// @Router	/rider/v2/cabinet/process [POST]
// @Summary	R4002 操作换电
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		model.RiderExchangeProcessReq	true	"desc"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*cabinet) Process() {}

// ProcessStatus
// @ID		CabinetProcessStatus
// @Router	/rider/v2/cabinet/process/status [GET]
// @Summary	R4003 换电状态
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string								true	"骑手校验token"
// @Param	query			query		model.RiderExchangeProcessStatusReq	true	"desc"
// @Success	200				{object}	model.RiderExchangeProcessRes		"请求成功"
func (*cabinet) ProcessStatus() {}

// Report
// @ID		CabinetReport
// @Router	/rider/v2/cabinet/report [POST]
// @Summary	R4004 电柜故障上报
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.CabinetFaultReportReq	true	"desc"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*cabinet) Report() {}

// Fault
// @ID		CabinetFault
// @Router	/rider/v2/cabinet/fault [GET]
// @Summary	R4008 电柜故障列表
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string		true	"骑手校验token"
// @Success	200				{object}	[]string	"请求成功"
func (*cabinet) Fault() {}

type city struct{}

// List
// @ID		CityList
// @Router	/rider/v2/city [GET]
// @Summary	R2003 获取已开通城市
// @Tags	City - 城市
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	[]model.CityWithLocation	"请求成功"
func (*city) List() {}

type contract struct{}

// Sign
// @ID		ContractSign
// @Router	/rider/v2/contract/sign [POST]
// @Summary	R3003 签署合同
// @Tags	Contract - 合同
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	body			body		model.ContractSignReq	true	"desc"
// @Success	200				{object}	model.ContractSignRes	"请求成功"
func (*contract) Sign() {}

// SignResult
// @ID		ContractSignResult
// @Router	/rider/v2/constract/{sn} [GET]
// @Summary	R3004 合同签署结果
// @Tags	Contract - 合同
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	sn				path		string					true	"合同编号"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*contract) SignResult() {}

type enterprise struct{}

// Battery
// @ID		EnterpriseBattery
// @Router	/rider/v2/enterprise/battery [GET]
// @Summary	R3010 企业骑手获取可用电池
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string		true	"骑手校验token"
// @Param	cityId			query		uint64		true	"城市ID"
// @Success	200				{object}	[]string	"请求成功"
func (*enterprise) Battery() {}

// Subscribe
// @ID		EnterpriseSubscribe
// @Router	/rider/v2/enterprise/subscribe [POST]
// @Summary	R3011 企业骑手选择电池
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string									true	"骑手校验token"
// @Param	body			body		model.EnterpriseRiderSubscribeChooseReq	true	"电池选择请求"
// @Success	200				{object}	model.EnterpriseRiderSubscribeChooseRes	"请求成功"
func (*enterprise) Subscribe() {}

// SubscribeStatus
// @ID		EnterpriseSubscribeStatus
// @Router	/rider/v2/enterprise/subscribe [GET]
// @Summary	R3012 企业骑手订阅激活状态
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string	true	"骑手校验token"
// @Param	id				query		uint64	true	"订阅ID"
// @Success	200				{object}	bool	"TRUE已激活, FALSE未激活"
func (*enterprise) SubscribeStatus() {}

// SubscribeAlter
// @ID		EnterpriseSubscribeAlter
// @Router	/rider/v2/enterprise/subscribe/alter [POST]
// @Summary	R3013 加时申请
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		model.SubscribeAlterRiderReq	true	"申请增加订阅时长请求"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*enterprise) SubscribeAlter() {}

// SubscribeAlterList
// @ID		EnterpriseSubscribeAlterList
// @Router	/rider/v2/enterprise/subscribe/alter/list [GET]
// @Summary	R3014 加时申请列表
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string															true	"骑手校验token"
// @Param	query			query		model.SubscribeListRiderReq										true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.SubscribeAlterApplyListRes}	"请求成功"
func (*enterprise) SubscribeAlterList() {}

// JoinEnterprise
// @ID		EnterpriseJoinEnterprise
// @Router	/rider/v2/enterprise/join [POST]
// @Summary	R3015 企业骑手加入团签
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	body	body		model.EnterproseInfoReq	true	"加入团签请求"
// @Success	200		{object}	bool					"请求成功"
func (s *enterprise) JoinEnterprise() {}

// RiderEnterpriseInfo
// @ID		EnterpriseRiderEnterpriseInfo
// @Router	/rider/v2/enterprise/info [GET]
// @Summary	R3016 骑手团签信息
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	query			query		model.EnterproseInfoReq	true	"团签信息请求"
// @Success	200				{object}	model.EnterproseInfoRsp	"请求成功"
func (*enterprise) RiderEnterpriseInfo() {}

// ExitEnterprise
// @ID		EnterpriseExitEnterprise
// @Router	/rider/v2/enterprise/exit [POST]
// @Summary	R3017 退出团签
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*enterprise) ExitEnterprise() {}

type exchange struct{}

// Store
// @ID		ExchangeStore
// @Router	/rider/v2/exchange/store [POST]
// @Summary	R4005 门店换电
// @Tags	Exchange - 换电
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	body			body		model.ExchangeStoreReq	true	"desc"
// @Success	200				{object}	model.ExchangeStoreRes	"请求成功"
func (*exchange) Store() {}

// Overview
// @ID		ExchangeOverview
// @Router	/rider/v2/exchange/overview [GET]
// @Summary	R4006 换电概览
// @Tags	Exchange - 换电
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	model.ExchangeOverview	"请求成功"
func (*exchange) Overview() {}

// Log
// @ID		ExchangeLog
// @Router	/rider/v2/exchange/log [GET]
// @Summary	R4007 换电记录
// @Tags	Exchange - 换电
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string													true	"骑手校验token"
// @Param	query			query		model.PaginationReq										true	"分页请求参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.ExchangeRiderListRes}	"请求成功"
func (*exchange) Log() {}

// Refund
// @ID		OrderRefund
// @Router	/rider/v2/order/refund [POST]
// @Summary	R3006 申请退款
// @Tags	Order - 订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string			true	"骑手校验token"
// @Param	body			body		model.RefundReq	true	"desc"
// @Success	200				{object}	model.RefundRes	"请求成功"
func (*order) Refund() {}

// List
// @ID		OrderList
// @Router	/rider/v2/order [GET]
// @Summary	R3007 骑士卡购买历史
// @Tags	Order - 订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string										true	"骑手校验token"
// @Param	query			query		model.PaginationReq							true	"desc"
// @Success	200				{object}	model.StatusResponse						"请求成功"
// @Success	200				{object}	model.PaginationRes{items=[]model.Order}	"请求成功"
func (*order) List() {}

// Detail
// @ID		OrderDetail
// @Router	/rider/v2/order/{id} [GET]
// @Summary	R3008 订单详情
// @Tags	Order - 订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string		true	"骑手校验token"
// @Param	id				path		int			true	"订单ID"
// @Success	200				{object}	model.Order	"请求成功"
func (*order) Detail() {}

// Status
// @ID		OrderStatus
// @Router	/rider/v2/order/status [GET]
// @Summary	R3009 订单支付状态
// @Tags	Order - 订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	outTradeNo		query		string					true	"订单编号"
// @Success	200				{object}	model.OrderStatusRes	"请求成功"
func (*order) Status() {}

type plan struct{}

// List
// @ID		PlanList
// @Router	/rider/v2/plan [GET]
// @Summary	R3001 新购骑士卡
// @Tags	Plan - 骑士卡
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	query			query		model.PlanListRiderReq	true	"骑士卡列表请求参数"
// @Success	200				{object}	model.PlanNewlyRes		"请求成功"
func (*plan) List() {}

// Renewly
// @ID		PlanRenewly
// @Router	/rider/v2/plan/renewly [GET]
// @Summary	R3002 续费骑士卡
// @Tags	Plan - 骑士卡
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	model.RiderPlanRenewalRes	"请求成功"
func (*plan) Renewly() {}

type reserve struct{}

// Unfinished
// @ID		ReserveUnfinished
// @Router	/rider/v2/reserve [GET]
// @Summary	R8001 获取未完成预约
// @Tags	Reserve - 预约
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	model.ReserveUnfinishedRes	"请求成功, 预约不存在时为`null`"
func (*reserve) Unfinished() {}

// Create
// @ID		ReserveCreate
// @Router	/rider/v2/reserve [POST]
// @Summary	R8002 创建预约
// @Tags	Reserve - 预约
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.ReserveCreateReq		true	"预约信息"
// @Success	200				{object}	model.ReserveUnfinishedRes	"请求成功"
func (*reserve) Create() {}

// Cancel
// @ID		ManagerReserveCancel
// @Router	/rider/v2/reserve/{id} [DELETE]
// @Summary	R8003 取消预约
// @Tags	Reserve - 预约
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	id				path		uint64					true	"预约ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*reserve) Cancel() {}

type rider struct{}

// Signin
// @ID		Signin
// @Router	/rider/v2/signin [POST]
// @Summary	R1001 登录或注册
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	body	body		model.RiderSignupReq	true	"desc"
// @Success	200		{object}	model.RiderSigninRes	"请求成功"
func (*rider) Signin() {}

// Contact
// @ID		Contact
// @Router	/rider/v2/contact [POST]
// @Summary	R1002 添加紧急联系人
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	body			body		model.RiderContact		true	"desc"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (r *rider) Contact() {}

// Authenticator
// @ID		Authenticator
// @Router	/rider/v2/authenticator [POST]
// @Summary	R1003 实名认证
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.RiderContact			true	"desc"
// @Success	200				{object}	model.FaceAuthUrlResponse	"请求成功"
func (*rider) Authenticator() {}

// AuthResult
// @ID		AuthResult
// @Router	/rider/v2/authenticator/{token} [GET]
// @Summary	R1004 实名认证结果
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	token			path		string					true	"实名认证token"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (r *rider) AuthResult() {}

// FaceResult
// @ID		FaceResult
// @Router	/rider/v2/face/{token} [GET]
// @Summary	R1005 获取人脸校验结果
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	token			path		string					true	"人脸校验token"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (r *rider) FaceResult() {}

func (r *rider) Demo() {}

// Profile
// @ID		Profile
// @Router	/rider/v2/profile [GET]
// @Summary	R1006 获取个人信息
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	model.RiderSigninRes	"请求成功"
func (r *rider) Profile() {}

// Deposit
// @ID		Deposit
// @Router	/rider/v2/deposit [GET]
// @Summary	R1007 获取已缴押金
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	model.RiderDepositRes	"请求成功"
func (*rider) Deposit() {}

// Deregister
// @ID		Deregister
// @Router	/rider/v2/deregister [DELETE]
// @Summary	R1008 注销账户
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*rider) Deregister() {}

type setting struct{}

// App
// @ID		SettingApp
// @Router	/rider/v2/setting/app [GET]
// @Summary	R6001 获取APP设置
// @Tags	Setting - 设置
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	model.SettingRiderApp	"请求成功"
func (*setting) App() {}

// Question
// @ID		SettingQuestion
// @Router	/rider/v2/setting/question [GET]
// @Summary	R6002 获取常见问题
// @Tags	Setting - 设置
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	[]model.SettingQuestion	"请求成功"
func (*setting) Question() {}

type wallet struct{}

// Overview
// @ID		WalletOverview
// @Router	/rider/v2/wallet/overview [GET]
// @Summary	R9001 钱包概览
// @Tags	Wallet - 钱包
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	model.WalletOverview	"请求成功"
func (*wallet) Overview() {}

// PointLog
// @ID		WalletPointLog
// @Router	/rider/v2/wallet/pointlog [GET]
// @Summary	R9002 积分日志
// @Tags	Wallet - 钱包
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string												true	"骑手校验token"
// @Param	query			query		model.PaginationReq									false	"分页选项"
// @Success	200				{object}	model.PaginationRes{items=[]model.PointLogListRes}	"请求成功"
func (*wallet) PointLog() {}

// Points
// @ID		WalletPoints
// @Router	/rider/v2/wallet/points [GET]
// @Summary	R9003 积分详情
// @Tags	Wallet - 钱包
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string			true	"骑手校验token"
// @Success	200				{object}	model.PointRes	"请求成功"
func (*wallet) Points() {}

// Coupons
// @ID		WalletCoupons
// @Router	/rider/v2/wallet/coupons [GET]
// @Summary	R9004 优惠券列表
// @Tags	Wallet - 钱包
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string				true	"骑手校验token"
// @Param	type			query		int					false	"查询类别 0:可使用 1:已使用 2:已过期"
// @Success	200				{object}	[]model.CouponRider	"请求成功"
func (*wallet) Coupons() {}
