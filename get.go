package apiclient

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gilcrest/errors"
	"github.com/gilcrest/servertoken"
)

// ViaServerToken fetches a client from the database using the client's
// ServerToken set in the context
func ViaServerToken(ctx context.Context, tx *sql.Tx) (*Client, error) {
	const op errors.Op = "apiclient/ViaServerToken"

	token, err := servertoken.FromCtx(ctx)
	if err != nil {
		return nil, errors.E(op, err)
	}

	client := new(Client)
	err = tx.QueryRowContext(ctx, `SELECT client_num,
										  client_id,
										  client_name,
										  server_token,
										  homepage_url,
										  app_description,
										  redirect_uri,
										  client_secret,
										  primary_username,
										  create_client_num,
										  create_timestamp,
										  modify_client_num,
										  modify_timestamp
									 FROM auth.client
									WHERE server_token = $1`, token).
		Scan(&client.Number, &client.ID, &client.Name,
			&client.ServerToken, &client.HomeURL, &client.Description, &client.RedirectURI, &client.Secret,
			&client.PrimaryUserID, &client.CreateClientNumber, &client.CreateTimestamp,
			&client.UpdateClientNumber, &client.UpdateTimestamp)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.E(op, fmt.Errorf("No client with token %s", token))
	case err != nil:
		return nil, errors.E(op, err)
	default:
		return client, nil
	}
}
