## 2022-09-28

- 备份数据库
- 删除所有外键
```postgresql
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
```