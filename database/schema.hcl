table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  column "password" {
    null = false
    type = text
  }
  column "enabled" {
    null = false
    type = boolean
  }
  primary_key {
    columns = [column.id]
  }
}
table "sessions" {
  schema = schema.public
  column "id" {
    null = false
    type = text
  }
  column "username" {
    null = false
    type = text
  }
  column "started" {
    null = false
    type = integer
  }
  column "ends" {
    null = false
    type = integer
  }
  primary_key {
    columns = [column.id]
  }
}
schema "public" {
  comment = "Default public gomin schema"
}
