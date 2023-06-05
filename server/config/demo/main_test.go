<<<<<<<< HEAD:server/config/demo/main_test.go
package demo_test
========
package transactions_test
>>>>>>>> 7fb86839 (fix: move transactions to it's own module (#2664)):server/transactions/main_test.go

import (
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/testmock"
)

func TestMain(m *testing.M) {
	testmock.StartTestEnvironment()

	exitVal := m.Run()

	testmock.StopTestEnvironment()

	os.Exit(exitVal)
}
