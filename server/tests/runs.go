package tests

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"strings"

	"github.com/kubeshop/tracetest/server/pkg/id"
)

const (
	createSequeceQuery = `CREATE SEQUENCE IF NOT EXISTS "` + runSequenceName + `";`
	dropSequeceQuery   = `DROP SEQUENCE IF EXISTS "` + runSequenceName + `";`
	runSequenceName    = "%sequence_name%"
)

func dropSequece(ctx context.Context, tx *sql.Tx, testID id.ID) error {
	_, err := tx.ExecContext(
		ctx,
		replaceRunSequenceName(dropSequeceQuery, testID),
	)

	return err
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func replaceRunSequenceName(sql string, testID id.ID) string {
	// postgres doesn't like uppercase chars in sequence names.
	// testID might contain uppercase chars, and we cannot lowercase them
	// because they might lose their uniqueness.
	// md5 creates a unique, lowercase hash.
	seqName := "runs_test_" + md5Hash(testID.String()) + "_seq"
	return strings.ReplaceAll(sql, runSequenceName, seqName)
}
