package misc

import (
	"fmt"
	"github.com/crossevol/sqlcc/internal/util"
	"testing"
)

func TestSprintf(t *testing.T) {
	s := fmt.Sprintf(`SELECT column_name, is_identity, column_default FROM information_schema.columns WHERE table_schema = 'public' AND table_name = %s AND column_default LIKE 'nextval(%%'`, util.Quote2("Account"))
	fmt.Println(s)
}
