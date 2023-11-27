# PostgreSQL 数据库

## 初始化

`sudo -u postgresql psql`


```postgresql
CREATE USER auroraride WITH PASSWORD 'QGyof&AYf8rSEg*#^zzjiXRN1ykWn*4#';
CREATE DATABASE auroraride OWNER auroraride;
GRANT ALL PRIVILEGES ON DATABASE auroraride TO auroraride;
```

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

```
str, args := query.Modify(func(s *sql.Selector) {
    s.SelectExpr(sql.Raw("COUNT(1) AS count"))
}).sqlQuery(context.Background()).Query()
fmt.Println(str, args)
```
