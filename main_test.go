package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err)
	defer db.Close()

	clientID := 1

	gotClient, err := selectClient(db, clientID)
	require.NoError(t, err)

	assert.Equal(t, clientID, gotClient.ID)

	assert.NotEmpty(t, gotClient.FIO)
	assert.NotEmpty(t, gotClient.Login)
	assert.NotEmpty(t, gotClient.Birthday)
	assert.NotEmpty(t, gotClient.Email)
}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err)
	defer db.Close()

	clientID := -1

	gotClient, err := selectClient(db, clientID)
	require.Equal(t, err, sql.ErrNoRows)

	assert.Empty(t, gotClient.ID)
	assert.Empty(t, gotClient.FIO)
	assert.Empty(t, gotClient.Login)
	assert.Empty(t, gotClient.Birthday)
	assert.Empty(t, gotClient.Email)
}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	cl.ID, err = insertClient(db, cl)
	require.NoError(t, err)
	require.NotEmpty(t, cl.ID)

	gotClient, err := selectClient(db, cl.ID)
	require.NoError(t, err)

	assert.Equal(t, cl, gotClient)
}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err)
	defer db.Close()

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	id, err := insertClient(db, cl)

	require.NoError(t, err)
	require.NotEmpty(t, id)

	gotClient, err := selectClient(db, id)
	require.NoError(t, err)
	require.Equal(t, cl.FIO, gotClient.FIO)
	require.Equal(t, cl.Login, gotClient.Login)
	require.Equal(t, cl.Email, gotClient.Email)
	require.Equal(t, cl.Birthday, gotClient.Birthday)

	err = deleteClient(db, cl.ID)
	require.NoError(t, err)

	got, err := selectClient(db, cl.ID)
	require.Equal(t, sql.ErrNoRows, err)
	require.Empty(t, got)
}
