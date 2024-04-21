package model

type AgreementUserType uint8

const (
	AgreementUserTypePersonal   AgreementUserType = iota + 1 // 个签协议
	AgreementUserTypeEnterprise                              // 团签协议
)

func (s AgreementUserType) Value() uint8 {
	return uint8(s)
}
