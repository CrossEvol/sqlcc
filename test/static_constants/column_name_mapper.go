package static_constants

import "fmt"

type todo_ struct {
	TABLE      string
	ALIAS      string
	ID         string
	TITLE      string
	STATUS     string
	CREATED_AT string
	PRIORITY   string
	CONTENT    string
	CREATED_BY string
}

func NewTodo_() todo_ {
	return todo_{
		TABLE:      "Todo",
		ID:         "id",
		TITLE:      "title",
		STATUS:     "status",
		CREATED_AT: "created_at",
		PRIORITY:   "priority",
		CONTENT:    "content",
		CREATED_BY: "created_by",
	}
}

func NewAliasTodo_(alias string) todo_ {
	withAlias := func(field string) string {
		return fmt.Sprintf("%s.%s", alias, field)
	}

	return todo_{
		TABLE:      "Todo",
		ALIAS:      alias,
		ID:         withAlias("id"),
		TITLE:      withAlias("title"),
		STATUS:     withAlias("status"),
		CREATED_AT: withAlias("created_at"),
		PRIORITY:   withAlias("priority"),
		CONTENT:    withAlias("content"),
		CREATED_BY: withAlias("created_by"),
	}
}

/*
CREATE TABLE "Todo" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "title" TEXT NOT NULL DEFAULT '',
  "status" TEXT NOT NULL DEFAULT 'pending',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "priority" TEXT NOT NULL DEFAULT 'high',
  "content" TEXT NOT NULL DEFAULT '',
  "created_by" TEXT NOT NULL DEFAULT 'admin',
);

*/
