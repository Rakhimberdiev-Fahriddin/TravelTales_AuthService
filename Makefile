CURRENT_DIR=$(shell pwd)
DBURL := postgres://postgres:0412@localhost:5432/nt?sslmode=disable

proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}


mig-down:
	migrate -path databases/migrations -database '${DBURL}' -verbose down

.PHONY: mig-up mig-force

mig-up:
	migrate -path databases/migrations -database 'postgres://postgres:0412@localhost:5432/nt?sslmode=disable' -verbose up

mig-force:
	migrate -path databases/migrations -database 'postgres://postgres:0412@localhost:5432/nt?sslmode=disable' -verbose force $(version)


mig-create-table:
	migrate create -ext sql -dir databases/migrations -seq create_Followers_table

swag-init:
	swag init -g api/router.go --output api/docs