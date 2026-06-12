# Generate go files from sql queries
gen_sql:
   sqlc generate

# Generate DB docs
db_docs:
   tbls doc --rm-dist

