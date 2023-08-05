package service

import (
	"context"
	stdsql "database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionearnings"
	"github.com/auroraride/aurservd/internal/ent/promotionreferrals"
	"github.com/auroraride/aurservd/internal/ent/promotionwithdrawal"

	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type promotionStatisticsService struct {
	ctx context.Context
}

func NewPromotionStatisticsService() *promotionStatisticsService {
	return &promotionStatisticsService{
		ctx: context.Background(),
	}
}

// Earnings 我的收益统计
func (s *promotionStatisticsService) Earnings(mem *ent.PromotionMember, req *promotion.StatisticsReq) promotion.StatisticsEarningsRes {
	var v []promotion.StatisticsEarningsRes
	status := promotion.EarningsStatusCanceled

	// 我的收益
	q := ent.Database.PromotionEarnings.Query().Where(
		promotionearnings.MemberID(mem.ID),
	).Modify(func(s *sql.Selector) {
		s.Select(
			sql.As(fmt.Sprintf("COALESCE(SUM ( CASE WHEN status <> %d AND (commission_rule_key = 'firstLevelNewSubscribe' OR commission_rule_key = 'secondLevelNewSubscribe') THEN amount ELSE 0 END ),0)", status), "totalNewSignEarnings"),
			sql.As(fmt.Sprintf("COALESCE(SUM ( CASE WHEN status <> %d AND (commission_rule_key = 'firstLevelRenewalSubscribe' OR commission_rule_key = 'secondLevelRenewalSubscribe') THEN amount ELSE 0 END ),0)", status), "totalRenewalEarnings"),
			sql.As(fmt.Sprintf("COALESCE(SUM ( CASE WHEN status <> %d THEN amount ELSE 0 END),0)", status), "totalEarnings"),
		)
	})
	if req.Start != nil || req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(promotionearnings.CreatedAtGTE(start), promotionearnings.CreatedAtLT(end))
	}
	q.ScanX(s.ctx, &v)
	return v[0]
}

// Team 首页 团队统计
func (s *promotionStatisticsService) Team(mem *ent.PromotionMember, req *promotion.StatisticsReq) promotion.StatisticsTeamRes {
	v := promotion.StatisticsTeamRes{}
	sqls := `
WITH RECURSIVE member_hierarchy AS (
    SELECT pr.referred_member_id, 1 AS level, pr.rider_id, r.created_at, r.person_id
    FROM promotion_referrals pr
    INNER JOIN rider r ON pr.rider_id = r.id
    ` + fmt.Sprintf(" WHERE referring_member_id = %d", mem.ID) + `

    UNION ALL

    SELECT mr.referred_member_id, mh.level + 1 AS level, mr.rider_id, mr.created_at, r.person_id
    FROM member_hierarchy mh
    INNER JOIN promotion_referrals mr ON mh.referred_member_id = mr.referring_member_id
    INNER JOIN rider r ON mr.rider_id = r.id
    WHERE mh.level < 2
)
SELECT
    COUNT(DISTINCT mh.referred_member_id) AS totalTeam,
    COUNT(DISTINCT CASE WHEN r.num_new_sign = 1 THEN mh.referred_member_id END) AS totalNewSign,
    COUNT(DISTINCT CASE WHEN r.num_renewals > 0 THEN mh.referred_member_id END) AS totalRenewal
FROM member_hierarchy mh
LEFT JOIN (
    SELECT person_id,
           SUM(CASE WHEN num_subscriptions = 1 AND status <> 0 THEN 1 ELSE 0 END) AS num_new_sign,
           SUM(CASE WHEN num_subscriptions > 1 THEN 1 ELSE 0 END) AS num_renewals
    FROM (
        SELECT r.person_id,
               COUNT(*) AS num_subscriptions,
               MAX(CASE WHEN status <> 0 THEN 1 ELSE 0 END) AS status
        FROM rider r
        LEFT JOIN subscribe s ON r.id = s.rider_id
        GROUP BY r.person_id
    ) AS subscription_stats
    GROUP BY person_id
) AS r ON mh.person_id = r.person_id
			`

	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start).Format(carbon.DateTimeLayout)
		end := tools.NewTime().ParseNextDateStringX(*req.End).Format(carbon.DateTimeLayout)
		sqls += fmt.Sprintf(" WHERE mh.created_at >= '%s' AND mh.created_at < '%s'", start, end)
	}
	rows, err := ent.Database.QueryContext(s.ctx, sqls)
	if err != nil {
		snag.Panic(err)
	}
	defer func(rows *stdsql.Rows) {
		err = rows.Close()
		if err != nil {
			snag.Panic(err)
		}
	}(rows)
	for rows.Next() {
		err = rows.Scan(&v.TotalTeam, &v.TotalNewSign, &v.TotalRenewal)
		if err != nil {
			snag.Panic(err)
		}
	}
	return v
}

// Record 首页 战绩统计
func (s *promotionStatisticsService) Record(mem *ent.PromotionMember) promotion.StatisticsRecordRes {
	var res []promotion.StatisticsRecordRes

	if mem == nil {
		return res[0]
	}

	// 累计收益
	ent.Database.PromotionEarnings.Query().Where(
		promotionearnings.MemberID(mem.ID),
		promotionearnings.StatusNEQ(promotion.EarningsStatusCanceled.Value()),
	).Aggregate(ent.As(ent.Sum(promotionearnings.FieldAmount), "totalEarnings")).ScanX(s.ctx, &res)

	// 累计邀请人数
	totalInvitation := ent.Database.PromotionReferrals.Query().Where(promotionreferrals.ReferringMemberID(mem.ID)).CountX(s.ctx)
	res[0].TotalInvitation = uint64(int64(totalInvitation))

	return res[0]
}

// Wallet 我的钱包统计
func (s *promotionStatisticsService) Wallet(mem *ent.PromotionMember) promotion.StatisticsWalletRes {
	var v []promotion.StatisticsWalletRes
	start := carbon.Now().StartOfDay()
	end := start.AddDay()
	// 昨天日期
	yesterdayStart := start.SubDay()
	yesterdayEnd := end.SubDay()
	if mem == nil {
		return v[0]
	}

	ent.Database.PromotionEarnings.Query().
		Where(
			promotionearnings.MemberID(mem.ID),
			promotionearnings.StatusNEQ(promotion.EarningsStatusCanceled.Value()),
		).Modify(func(s *sql.Selector) {
		s.Select(
			sql.As("COALESCE(SUM(amount),0)", "totalEarnings"),
			sql.As(fmt.Sprintf("COALESCE(SUM(CASE WHEN created_at >= '%s' AND created_at < '%s' THEN amount ELSE 0 END),0)", start.String(), end.String()), "todayEarnings"),
			sql.As(fmt.Sprintf("COALESCE(SUM(CASE WHEN created_at >= '%s' AND created_at < '%s' THEN amount ELSE 0 END),0)", yesterdayStart.String(), yesterdayEnd.String()), "yesterdayEarnings"),
		)
	}).ScanX(s.ctx, &v)

	v[0].TotalBalance = mem.Balance + mem.Frozen
	v[0].Balance = mem.Balance
	v[0].FreezeBalance = mem.Frozen
	v[0].StatisticsWithdrawalRes = s.Withdrawal(mem)

	return v[0]
}

// Withdrawal 我的钱包提现统计
func (s *promotionStatisticsService) Withdrawal(mem *ent.PromotionMember) promotion.StatisticsWithdrawalRes {
	var v []promotion.StatisticsWithdrawalRes

	ent.Database.PromotionWithdrawal.Query().Where(
		promotionwithdrawal.MemberID(mem.ID),
		promotionwithdrawal.Status(promotion.WithdrawalStatusSuccess.Value()),
	).Aggregate(ent.As(ent.Sum(promotionwithdrawal.FieldAmount), "totalWithdrawal")).ScanX(s.ctx, &v)
	return v[0]
}

// EarningsDetail 我的钱包收益明细统计
func (s *promotionStatisticsService) EarningsDetail(mem *ent.PromotionMember, req *promotion.StatisticsReq) promotion.StatisticsEarningsDetailRes {
	res := promotion.StatisticsEarningsDetailRes{}
	res.FirstLevelEarnings = s.FirstLevelEarnings(mem, req)
	res.SecondLevelEarnings = s.SecondLevelEarnings(mem, req)
	return res
}

// FirstLevelEarnings 一级收益
func (s *promotionStatisticsService) FirstLevelEarnings(mem *ent.PromotionMember, req *promotion.StatisticsReq) promotion.StatisticsEarningsDetail {
	var v []promotion.StatisticsEarningsDetail
	status := promotion.EarningsStatusCanceled

	q := ent.Database.PromotionEarnings.Query().Where(
		promotionearnings.MemberID(mem.ID),
	).Modify(func(s *sql.Selector) {
		s.Select(
			sql.As(fmt.Sprintf("COALESCE(SUM ( CASE WHEN status <> %d AND commission_rule_key = 'firstLevelNewSubscribe' THEN amount ELSE 0 END ),0)", status), "totalNewSignEarnings"),
			sql.As(fmt.Sprintf("COALESCE(SUM ( CASE WHEN status <> %d AND commission_rule_key = 'firstLevelRenewalSubscribe' THEN amount ELSE 0 END ),0)", status), "totalRenewalEarnings"),
			sql.As(fmt.Sprintf("COALESCE(SUM ( CASE WHEN status <> %d AND (commission_rule_key = 'firstLevelNewSubscribe' OR commission_rule_key = 'firstLevelRenewalSubscribe') THEN amount ELSE 0 END ),0)", status), "totalEarnings"),
		)
	})
	if req.Start != nil || req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(promotionearnings.CreatedAtGTE(start), promotionearnings.CreatedAtLT(end))
	}
	q.ScanX(s.ctx, &v)
	return v[0]
}

// SecondLevelEarnings 二级收益
func (s *promotionStatisticsService) SecondLevelEarnings(mem *ent.PromotionMember, req *promotion.StatisticsReq) promotion.StatisticsEarningsDetail {
	var v []promotion.StatisticsEarningsDetail
	status := promotion.EarningsStatusCanceled

	q := ent.Database.PromotionEarnings.Query().Where(
		promotionearnings.MemberID(mem.ID),
	).Modify(func(s *sql.Selector) {
		s.Select(
			sql.As(fmt.Sprintf("COALESCE(SUM ( CASE WHEN status <> %d AND commission_rule_key = 'secondLevelNewSubscribe' THEN amount ELSE 0 END ),0)", status), "totalNewSignEarnings"),
			sql.As(fmt.Sprintf("COALESCE(SUM ( CASE WHEN status <> %d AND commission_rule_key = 'secondLevelRenewalSubscribe' THEN amount ELSE 0 END ),0)", status), "totalRenewalEarnings"),
			sql.As(fmt.Sprintf("COALESCE(SUM ( CASE WHEN status <> %d AND (commission_rule_key = 'secondLevelNewSubscribe' OR commission_rule_key = 'secondLevelRenewalSubscribe') THEN amount ELSE 0 END ),0)", status), "totalEarnings"),
		)
	})
	if req.Start != nil || req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(promotionearnings.CreatedAtGTE(start), promotionearnings.CreatedAtLT(end))
	}
	q.ScanX(s.ctx, &v)
	return v[0]
}

// MyTeamStatistics 我的-团队统计
func (s *promotionStatisticsService) MyTeamStatistics(mem *ent.PromotionMember) []promotion.StatisticsTeamDetailRes {
	var v []promotion.StatisticsTeamDetailRes
	sqls := `
WITH RECURSIVE member_hierarchy AS (
    SELECT referred_member_id, 1 AS level, rider_id, pr.created_at, person_id
    FROM promotion_referrals pr
    INNER JOIN rider r ON pr.rider_id = r.id
    ` + fmt.Sprintf(" WHERE referring_member_id = %d", mem.ID) + `

    UNION ALL

    SELECT mr.referred_member_id, mh.level + 1 AS level, mr.rider_id, mr.created_at, r.person_id
    FROM member_hierarchy mh
    INNER JOIN promotion_referrals mr ON mh.referred_member_id = mr.referring_member_id
    INNER JOIN rider r ON mr.rider_id = r.id
    WHERE mh.level < 2
)

SELECT
    level,
    COUNT(DISTINCT mh.referred_member_id) AS TotalTeam,
    COUNT(DISTINCT CASE WHEN r.num_new_sign = 1 THEN mh.referred_member_id END) AS TotalNewSign,
    COUNT(DISTINCT CASE WHEN r.num_renewals > 0 THEN mh.referred_member_id END) AS TotalRenewal
FROM member_hierarchy mh
LEFT JOIN (
    SELECT person_id,
           SUM(CASE WHEN num_subscriptions = 1 AND status <> 0 THEN 1 ELSE 0 END) AS num_new_sign,
           SUM(CASE WHEN num_subscriptions >= 2 THEN 1 ELSE 0 END) AS num_renewals
    FROM (
        SELECT r.person_id,
               COUNT(*) AS num_subscriptions,
               MAX(CASE WHEN status <> 0 THEN 1 ELSE 0 END) AS status
        FROM rider r
        LEFT JOIN subscribe s ON r.id = s.rider_id
        GROUP BY r.person_id
    ) AS subscription_stats
    GROUP BY person_id
) AS r ON mh.person_id = r.person_id
GROUP BY level;
			`
	rows, err := ent.Database.QueryContext(s.ctx, sqls)
	if err != nil {
		snag.Panic(err)
	}
	defer func(rows *stdsql.Rows) {
		err = rows.Close()
		if err != nil {
			snag.Panic(err)
		}
	}(rows)
	for rows.Next() {
		row := promotion.StatisticsTeamDetailRes{}
		err = rows.Scan(&row.Level, &row.TotalTeam, &row.TotalNewSign, &row.TotalRenewal)
		if err != nil {
			snag.Panic(err)
		}
		v = append(v, row)
	}
	return v
}

// getYearMonthSlice 获取年月切片
func (s *promotionStatisticsService) getYearMonthSlice(startDate, endDate string) (res []promotion.StatisticsTeamGrowthTrendRes) {
	layout := "2006-01-02"
	start, _ := time.Parse(layout, startDate)
	end, _ := time.Parse(layout, endDate)

	current := start
	for current.Before(end) || current.Equal(end) {
		row := promotion.StatisticsTeamGrowthTrendRes{}
		row.Month = current.Format("2006-01")
		res = append(res, row)
		current = current.AddDate(0, 1, 0)
	}
	return res
}

func (s *promotionStatisticsService) TeamGrowth(mem *ent.PromotionMember, req *promotion.StatisticsReq) []promotion.StatisticsTeamGrowthTrendRes {

	if req.Start == nil || req.End == nil {
		snag.Panic("请查询6个月内的数据")
	}
	startMonth := carbon.ParseByLayout(*req.Start, carbon.DateLayout).StartOfMonth()
	endMonth := carbon.ParseByLayout(*req.End, carbon.DateLayout).EndOfMonth()
	var res []promotion.StatisticsTeamGrowthTrendRes

	if startMonth.DiffAbsInMonths(endMonth) > 6 {
		snag.Panic("请查询6个月内的数据")
	}

	res = s.getYearMonthSlice(*req.Start, *req.End)

	sqls := `WITH RECURSIVE member_hierarchy AS (
			  SELECT referred_member_id, 1 AS level,rider_id,created_at
			  FROM promotion_referrals
		  ` + fmt.Sprintf(" WHERE referring_member_id = %d", mem.ID) + `
			  UNION ALL
			  SELECT mr.referred_member_id, mh.level + 1 AS level,mr.rider_id,mr.created_at
			  FROM member_hierarchy mh
			  INNER JOIN promotion_referrals as mr ON mh.referred_member_id = mr.referring_member_id
			  WHERE mh.level < 2
			)
			SELECT
			  TO_CHAR(DATE_TRUNC('month', created_at),'YYYY-MM') AS month,
			  SUM(CASE WHEN level = 1 THEN 1 ELSE 0 END) AS first_level_referrals,
			  SUM(CASE WHEN level = 2 THEN 1 ELSE 0 END) AS second_level_referrals
			FROM member_hierarchy
			WHERE
			    DATE_TRUNC('month', created_at) >= '` + startMonth.Format("Y-m-d") + `' AND
			    DATE_TRUNC('month', created_at) <= '` + endMonth.Format("Y-m-d") + `'
			GROUP BY DATE_TRUNC('month', created_at)
			ORDER BY DATE_TRUNC('month', created_at);
`
	rows, err := ent.Database.QueryContext(s.ctx, sqls)
	if err != nil {
		snag.Panic(err)
	}
	defer func(rows *stdsql.Rows) {
		err = rows.Close()
		if err != nil {
			snag.Panic(err)
		}
	}(rows)
	for rows.Next() {
		row := promotion.StatisticsTeamGrowthTrendRes{}
		err = rows.Scan(&row.Month, &row.FirstLevelNum, &row.SecondLevelNum)
		if err != nil {
			snag.Panic(err)
		}
		for k, v := range res {
			if row.Month == v.Month {
				res[k].FirstLevelNum = row.FirstLevelNum
				res[k].SecondLevelNum = row.SecondLevelNum
				break
			}
		}
	}
	return res
}
