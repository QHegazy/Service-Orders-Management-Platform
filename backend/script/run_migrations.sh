#!/bin/bash

# Load environment variables from the .env file
if [ -f ".env" ]; then
  export $(grep -v '^#' .env | xargs)

  if [ -n "${DB_URL:-}" ]; then
    echo "✅ DB_URL found: $DB_URL"
  else
    echo "❌ DB_URL (or DATABASE_URL) not found in .env"
  fi
else
  echo "Error: .env file not found."
  exit 1
fi


MIGRATION_PATH="./internal/db/migrations" 

# Apply migrations (run the up migration)
apply_migrations() {
  echo "Applying migrations..."
  goose -dir $MIGRATION_PATH postgres $DB_URL up
  if [ $? -eq 0 ]; then
    echo "Migrations applied successfully."
  else
    echo "Error applying migrations."
    exit 1
  fi
}

# Rollback migrations (run the down migration)
rollback_migrations() {
  echo "Rolling back migrations..."
  goose -dir $MIGRATION_PATH postgres $DB_URL down
  if [ $? -eq 0 ]; then
    echo "Migrations rolled back successfully."
  else
    echo "Error rolling back migrations."
    exit 1
  fi
}

reset_migrations() {
  echo "Resetting migrations..."
  goose -dir "$MIGRATION_PATH" postgres "$DB_URL" reset
  if [ $? -eq 0 ]; then
    echo "Migrations reset successfully."
  else
    echo "Error resetting migrations."
    exit 1
  fi
}

# Force a specific migration version
force_version() {
  echo "Forcing migration version..."
  goose -dir $MIGRATION_PATH postgres $DB_URL force "$2"
  if [ $? -eq 0 ]; then
    echo "Migration version forced successfully."
  else
    echo "Error forcing migration version."
    exit 1
  fi
}

# Show current migration version
show_version() {
  echo "Current migration version:"
  goose -dir $MIGRATION_PATH postgres $DB_URL version
}

# Show usage information
usage() {
  echo "Usage: $0 {apply|rollback|version|force <version>}"
  echo "  apply   - Apply the up migrations"
  echo "  rollback- Rollback the down migrations"
  echo "  version - Show the current migration version"
  echo "  force <version> - Force the migration to a specific version"
  exit 0
}

# Main execution: Check arguments and run the appropriate function
if [ $# -eq 0 ]; then
  usage
  exit 1
fi

case "$1" in
  up)
    apply_migrations
    ;;
  down)
    rollback_migrations
    ;;
  version)
    show_version
    ;;
  reset)
    reset_migrations
    ;;
  force)
    if [ $# -eq 2 ]; then
      force_version "$1" "$2"
    else
      echo "Error: force command requires a version argument."
      exit 1
    fi
    ;;
  *)
    usage
    exit 1
    ;;
esac
