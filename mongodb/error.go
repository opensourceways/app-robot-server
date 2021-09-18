package mongodb

import (
	"fmt"

	"github.com/opensourceways/app-robot-server/dbmodels"
)

var (
	errNoDBRecord = dbError{code: dbmodels.ErrNoDBRecord, err: fmt.Errorf("no record")}
)

type dbError struct {
	code dbmodels.DBErrCode
	err  error
}

func (dbe dbError) Error() string {
	if dbe.err == nil {
		return ""
	}
	return dbe.err.Error()
}

func (dbe dbError) IsErrorOf(code dbmodels.DBErrCode) bool {
	return dbe.code == code
}

func (dbe dbError) ErrCode() dbmodels.DBErrCode {
	return dbe.code
}

func newDBError(code dbmodels.DBErrCode, err error) dbmodels.IDBError {
	return dbError{code: code, err: err}
}

func newSystemError(err error) dbmodels.IDBError {
	return dbError{code: dbmodels.ErrSystemError, err: err}
}
