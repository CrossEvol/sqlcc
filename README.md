This repo is a helper tool for sqlc.

# Todo
## when generate insert or update stmt, should remove the id column based on generated approarch

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

## how to verify a column is pk && auto_increment 
### sqlite
```sqlite
PRAGMA table_info(downloaded_image);
``` 

### mysql
`auto_increment`
```mysql
SELECT COLUMN_NAME
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'answer'
AND TABLE_NAME = 'activity'
AND EXTRA = 'auto_increment';
```

`primary key`
```mysql
SELECT COLUMN_NAME
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = 'answer'
AND TABLE_NAME = 'activity'
AND CONSTRAINT_NAME = 'PRIMARY';
```

### postgres
`auto_increment`
```postgresql
SELECT column_name, is_identity, column_default
FROM information_schema.columns
WHERE table_schema = 'public'
AND table_name = 'Post' AND column_default LIKE 'nextval(%';
```