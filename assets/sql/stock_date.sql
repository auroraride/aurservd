SELECT date, store_id, cabinet_id, SUM(cnt) OVER (PARTITION BY store_id, cabinet_id ORDER BY date ASC)
FROM (SELECT date(created_at) AS date, SUM(num) AS cnt, store_id, cabinet_id
      FROM stock
      WHERE (store_id IS NOT NULL OR cabinet_id IS NOT NULL) AND date(created_at) <= '2022-08-14'
      GROUP BY date(created_at), store_id, cabinet_id) t
GROUP BY store_id, date, cnt, cabinet_id;