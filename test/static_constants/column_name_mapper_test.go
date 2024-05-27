package static_constants

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapper(t *testing.T) {
	t1 := NewTodo_()
	require.Equal(t, "Todo", t1.TABLE)
	require.Equal(t, "id", t1.ID)
	require.Equal(t, "title", t1.TITLE)
	require.Equal(t, "status", t1.STATUS)
	require.Equal(t, "created_at", t1.CREATED_AT)
	require.Equal(t, "priority", t1.PRIORITY)
	require.Equal(t, "content", t1.CONTENT)
	require.Equal(t, "created_by", t1.CREATED_BY)

	t2 := NewAliasTodo_("t")
	require.Equal(t, "Todo", t2.TABLE)
	require.Equal(t, "t", t2.ALIAS)
	require.Equal(t, "t.id", t2.ID)
	require.Equal(t, "t.title", t2.TITLE)
	require.Equal(t, "t.status", t2.STATUS)
	require.Equal(t, "t.created_at", t2.CREATED_AT)
	require.Equal(t, "t.priority", t2.PRIORITY)
	require.Equal(t, "t.content", t2.CONTENT)
	require.Equal(t, "t.created_by", t2.CREATED_BY)
}
