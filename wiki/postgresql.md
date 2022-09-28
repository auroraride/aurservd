# PostgreSQL 数据库

## 预处理

```postgresql
SET TIMEZONE = 'Asia/Shanghai';
CREATE EXTENSION postgis;
CREATE EXTENSION pg_trgm;
CREATE EXTENSION btree_gin;
```

## 备份

```postgresql
pg_dump -U postgres -d auroraride -f auroraride.sql
```

## 还原

```postgresql
psql -d auroraride-prod -U liasica -f auroraride.sql
```

## 优化行数统计效率

- [Faster PostgreSQL Counting](https://dzone.com/articles/faster-postgresql-counting)
- [提升 PostgreSQL 中 Count 的速度](https://www.oschina.net/translate/faster-postgresql-counting)