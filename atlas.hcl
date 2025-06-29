env "local" {
  src = "file://migrations/schema.sql"
  url = "postgres://postgres:password@localhost:5432/sub_cos_counter?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}