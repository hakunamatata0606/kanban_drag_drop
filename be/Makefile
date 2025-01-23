migrate_up:
	migrate -path db/migrations/ -database "mysql://root:123@tcp(localhost:13306)/kanban_drag_drop" -verbose up
migrate_down:
	migrate -path db/migrations/ -database "mysql://root:123@tcp(localhost:13306)/kanban_drag_drop" -verbose down
sqlc:
	sqlc generate
.PHONY: migrate_up, migrate_down, sqlc