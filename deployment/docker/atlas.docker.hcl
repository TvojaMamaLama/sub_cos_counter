env "docker" {
  src = "file://migrations/schema.sql"
  url = env("ATLAS_DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}