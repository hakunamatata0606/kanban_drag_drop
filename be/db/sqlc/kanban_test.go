package db_test

import (
	"context"
	"example/kanban/appstate"
	db "example/kanban/db/sqlc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTasks(t *testing.T) {
	appState := appstate.GetAppState()
	ctx := context.Background()
	err := db.CreateTask(
		ctx,
		appState.Db,
		&db.CreateTaskRequest{
			Name:        "test1",
			Title:       "title1",
			Description: "desc1",
			Status:      "idea",
		},
	)
	assert.Nil(t, err)
	err = db.CreateTask(
		ctx,
		appState.Db,
		&db.CreateTaskRequest{
			Name:        "test2",
			Title:       "title2",
			Description: "desc2",
			Status:      "done",
		},
	)
	assert.Nil(t, err)

	listTasks, err := db.ListTasks(ctx, appState.Db)
	assert.Nil(t, err)
	listTasksExpected := []db.ListTasksRow{
		db.ListTasksRow{
			Name:        "test1",
			Title:       "title1",
			Description: "desc1",
			Status:      "idea",
		},
		db.ListTasksRow{
			Name:        "test2",
			Title:       "title2",
			Description: "desc2",
			Status:      "done",
		},
	}
	assert.Equal(t, listTasksExpected, listTasks)

	err = db.DeleteTask(ctx, appState.Db, &db.DeleteTaskRequest{Name: "test1"})
	assert.Nil(t, err)
	err = db.DeleteTask(ctx, appState.Db, &db.DeleteTaskRequest{Name: "test2"})
	assert.Nil(t, err)
}

func TestCreateTask(t *testing.T) {
	appState := appstate.GetAppState()
	ctx := context.Background()
	err := db.CreateTask(
		ctx,
		appState.Db,
		&db.CreateTaskRequest{
			Name:        "test1",
			Title:       "title1",
			Description: "desc1",
			Status:      "idea",
		},
	)
	assert.Nil(t, err)
	err = db.CreateTask(
		ctx,
		appState.Db,
		&db.CreateTaskRequest{
			Name:        "test1",
			Title:       "title1",
			Description: "desc1",
			Status:      "idea",
		},
	)
	assert.NotNil(t, err)

	err = db.CreateTask(
		ctx,
		appState.Db,
		&db.CreateTaskRequest{
			Name:        "test2",
			Title:       "title1",
			Description: "desc1",
			Status:      "lala",
		},
	)
	assert.NotNil(t, err)

	err = db.UpdateTask(ctx, appState.Db, &db.UpdateTaskRequest{Name: "test1", Status: "done"})
	assert.Nil(t, err)
	listTasks, err := db.ListTasks(ctx, appState.Db)
	assert.Nil(t, err)
	listTasksExpected := []db.ListTasksRow{
		db.ListTasksRow{
			Name:        "test1",
			Title:       "title1",
			Description: "desc1",
			Status:      "done",
		},
	}
	assert.Equal(t, listTasksExpected, listTasks)

	err = db.DeleteTask(ctx, appState.Db, &db.DeleteTaskRequest{Name: "test1"})
	assert.Nil(t, err)
}
