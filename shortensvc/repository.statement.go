package shortensvc

import "database/sql"

type (
	StmtID uint

	Statements map[StmtID]*sql.Stmt
)

const (
	// Skipt the first one, this makes it easier to sort the statements
	// without having to arrange the iota back to the first line.
	_ StmtID = iota
	createStmt
	withIDStmt
)

var rawStmts = map[StmtID]string{
	createStmt: `
		INSERT INTO url (long_url)
		VALUES (?)
	`,
	withIDStmt: `
		SELECT 	* 
		FROM 	url 
		WHERE 	long_url = ?
	`,
}
