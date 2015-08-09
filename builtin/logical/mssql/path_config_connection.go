package mssql

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigConnection(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/connection",
		Fields: map[string]*framework.FieldSchema{
			"value": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "DB connection string",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathConnectionWrite,
		},

		HelpSynopsis:    pathConfigConnectionHelpSyn,
		HelpDescription: pathConfigConnectionHelpDesc,
	}
}

func (b *backend) pathConnectionWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	connString := data.Get("value").(string)

	// Verify the string
	db, err := sql.Open("mssql", connString)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error validating connection info: %s", err)), nil
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error validating connection info: %s", err)), nil
	}

	// Store it
	entry, err := logical.StorageEntryJSON("config/connection", connString)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	// Reset the DB connection
	b.ResetDB()
	return nil, nil
}

const pathConfigConnectionHelpSyn = `
Configure the connection string to talk to MsSQL.
`

const pathConfigConnectionHelpDesc = `
This path configures the connection string used to connect to MsSQL.
The value of the string is a Data Source Name (DSN). An example is
using "server=localhost,1433;user id=sa;password=password;initial catalog=db"

When configuring the connection string, the backend will verify its validity.
`
