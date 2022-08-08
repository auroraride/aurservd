SELECT model,
    SUM(num)                                                                                                             AS num,

    -- 平台调拨
    SUM(CASE WHEN type = 0 AND cabinet_id IS NOT NULL THEN num ELSE 0 END)                                               AS c_transfer,
    SUM(CASE WHEN type = 0 AND store_id IS NOT NULL OR (store_id IS NULL AND cabinet_id IS NUll) THEN num ELSE 0 END)    AS s_transfer,

    -- 调拨数据
    SUM(CASE WHEN num < 0 AND cabinet_id IS NOT NULL THEN -num ELSE 0 END)                                               AS c_outboundNum,
    SUM(CASE WHEN num > 0 AND cabinet_id IS NOT NULL THEN num ELSE 0 END)                                                AS c_inboundNum,
    SUM(CASE WHEN num < 0 AND (store_id IS NOT NULL OR (store_id IS NULL AND cabinet_id IS NUll)) THEN -num ELSE 0 END)  AS s_outboundNum,
    SUM(CASE WHEN num > 0 AND (store_id IS NOT NULL OR (store_id IS NULL AND cabinet_id IS NUll)) THEN num ELSE 0 END)   AS s_inboundNum,

    -- 业务数据
    SUM(CASE WHEN type = 1 AND cabinet_id IS NOT NULL THEN -num ELSE 0 END)                                              AS c_active,
    SUM(CASE WHEN type = 2 AND cabinet_id IS NOT NULL THEN num ELSE 0 END)                                               AS c_pause,
    SUM(CASE WHEN type = 3 AND cabinet_id IS NOT NULL THEN -num ELSE 0 END)                                              AS c_continue,
    SUM(CASE WHEN type = 4 AND cabinet_id IS NOT NULL THEN num ELSE 0 END)                                               AS c_unsubscribe,
    SUM(CASE WHEN type = 1 AND (store_id IS NOT NULL OR (store_id IS NULL AND cabinet_id IS NUll)) THEN -num ELSE 0 END) AS s_active,
    SUM(CASE WHEN type = 2 AND (store_id IS NOT NULL OR (store_id IS NULL AND cabinet_id IS NUll)) THEN num ELSE 0 END)  AS s_pause,
    SUM(CASE WHEN type = 3 AND (store_id IS NOT NULL OR (store_id IS NULL AND cabinet_id IS NUll)) THEN -num ELSE 0 END) AS s_continue,
    SUM(CASE WHEN type = 4 AND (store_id IS NOT NULL OR (store_id IS NULL AND cabinet_id IS NUll)) THEN num ELSE 0 END)  AS s_unsubscribe
FROM stock
WHERE model IS NOT NULL
        AND (store_id IS NOT NULL OR cabinet_id IS NOT NULL OR rider_id IS NOT NULL)
GROUP BY model;