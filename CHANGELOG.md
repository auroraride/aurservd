## 2022-09-28

- 备份数据库
- 删除所有外键
```postgresql
ALTER TABLE plan_pms RENAME TO plan_models;
ALTER TABLE cabinet_bms RENAME TO cabinet_models;

DO
$$
    DECLARE
        r RECORD;
    BEGIN
        FOR r IN (SELECT 'ALTER TABLE ' || QUOTE_IDENT(ns.nspname) || '.' || QUOTE_IDENT(tb.relname) ||
                         ' DROP CONSTRAINT ' || QUOTE_IDENT(conname) || ';' AS sql
                  FROM pg_constraint c
                           JOIN pg_class tb ON tb.oid = c.conrelid
                           JOIN pg_namespace ns ON ns.oid = tb.relnamespace
                  WHERE ns.nspname IN ('public') AND c.contype = 'f')
            LOOP
                EXECUTE r.sql;
            END LOOP;
    END;
$$;

UPDATE rider r
SET name = p.name, id_card_number = p.id_card_number
FROM person p
WHERE p.id = r.person_id;

```