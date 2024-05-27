This repo is a helper tool for sqlc.

# Notes

## how to get column name from database

### mysql
```mysql
SELECT DISTINCT COLUMN_NAME, COLUMN_TYPE FROM information_schema.`COLUMNS` WHERE TABLE_NAME = 'activity' AND
TABLE_SCHEMA = 'answer';
SELECT DISTINCT COLUMN_NAME, COLUMN_TYPE FROM information_schema.`COLUMNS` WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ?;
```

### postgres
```postgresql
SELECT DISTINCT "column_name", data_type FROM information_schema.COLUMNS WHERE table_schema = 'public' AND "
table_name" = 'Post';
SELECT DISTINCT "column_name", data_type FROM information_schema.COLUMNS WHERE table_schema = 'public' AND "
table_name" = ?;
```

### sqlite
```sqlite
PRAGMA table_info('Todo');
PRAGMA table_info(?);
```