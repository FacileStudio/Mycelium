#!/usr/bin/env sh
#
# Counts what every suite database on this host still holds for an erased
# person. Read-only: every query runs in a read-only transaction and only
# counts come back, never a row, never a credential.
#
#   sh scripts/erasure-check.sh <registre subject> <email>
#
# Each app's Postgres container is found by its name pattern. Every text and
# uuid column of every table is compared with the subject and the email
# (case-insensitive), and Synapse also with the person's Matrix id. Bulle's
# per-bot crypto stores (<db>_mx_<user>_*) are counted by name, and Mycelium's
# JSON stores are grepped inside a read-only mount of its volume. Only
# non-zero counts are printed, then one total per app.

set -eu

subject="${1:-}"
email="${2:-}"
case "$subject" in
*[!0-9a-f-]* | "") echo "usage: $0 <subject uuid> <email>" >&2; exit 2 ;;
esac
case "$email" in
*[\'\"\\\ ]* | "" | *[!@]*@*@*) echo "usage: $0 <subject uuid> <email>" >&2; exit 2 ;;
esac
mxid="@${email%%@*}:${MATRIX_SERVER:-prisme.facile.studio}"

# app  container pattern  database
APPS='
registre registre-db registre
bulle compose-*-db-1:bulle bulle
courrier facile-courrier-*-db-1 courrier
nuage nuage-db nuage
agenda facile-agenda-*-db-1 agenda
opus opus-db opus
sablier facile-sablier-*-db-1 sablier
sonar sonar-db sonar
glouton glouton-postgres glouton
ardoise facile-charles-*-ardoise-postgres-1 ardoise
prisme facile-prisme-*-db-1 prisme
synapse facile-prisme-*-db-1 synapse
'

# The container whose name matches the pattern. A pattern with ":<db>" picks
# among several matches the one that holds that database.
find_container() {
  pattern="${1%%:*}"
  want="${1#*:}"
  [ "$want" = "$1" ] && want=""
  docker ps --format '{{.Names}}' | while read -r name; do
    case "$name" in
    $pattern)
      if [ -z "$want" ] || docker exec "$name" sh -c 'psql -U "$POSTGRES_USER" -At -c "select 1 from pg_database where datname = '"'$want'"'"' 2>/dev/null | grep -q 1; then
        echo "$name"
        return
      fi
      ;;
    esac
  done
}

# One statement lists every text or uuid column of every base table with the
# number of rows equal to one of the needles. The needles travel as psql
# variables and are quoted with %L inside format(), so nothing in them is
# ever parsed as SQL. The statement goes in on stdin: psql interpolates
# variables there and not in a -c command.
SCAN="
SELECT table_name || '.' || column_name, n FROM (
  SELECT c.table_name, c.column_name,
    (xpath('/row/n/text()', query_to_xml(format(
      'SELECT count(*) AS n FROM %I.%I WHERE %I::text = %L OR lower(%I::text) = lower(%L) OR %I::text = %L',
      c.table_schema, c.table_name, c.column_name, :'subject', c.column_name, :'email', c.column_name, :'mxid'),
      false, true, '')))[1]::text::int AS n
  FROM information_schema.columns c
  JOIN information_schema.tables t USING (table_schema, table_name)
  WHERE t.table_type = 'BASE TABLE'
    AND c.table_schema NOT IN ('pg_catalog', 'information_schema')
    AND c.data_type IN ('text', 'character varying', 'uuid', 'citext')
) counts WHERE n > 0 ORDER BY 1"

scan_db() {
  docker exec -e PGOPTIONS='-c default_transaction_read_only=on' "$1" sh -c \
    'printf "%s" "$5" | psql -U "$POSTGRES_USER" -d "$1" -At -F " " -q -v ON_ERROR_STOP=1 -v subject="$2" -v email="$3" -v mxid="$4"' \
    sh "$2" "$subject" "$email" "$mxid" "$SCAN"
}

# Bulle keeps one database per bot device, named after the user id.
BULLE_STORES="
SELECT 'stores.of_the_person', count(*) FROM pg_database d
WHERE EXISTS (SELECT 1 FROM users u LEFT JOIN porte_identities i ON i.user_id = u.id
              WHERE (i.subject = :'subject' OR lower(u.email) = lower(:'email'))
                AND starts_with(d.datname, 'bulle_mx_' || u.id || '_'))
UNION ALL
SELECT 'stores.without_a_users_row', count(*) FROM pg_database d
WHERE starts_with(d.datname, 'bulle_mx_')
  AND NOT EXISTS (SELECT 1 FROM users u WHERE starts_with(d.datname, 'bulle_mx_' || u.id || '_'))"

bulle_stores() {
  docker exec -e PGOPTIONS='-c default_transaction_read_only=on' "$1" sh -c \
    'printf "%s" "$3" | psql -U "$POSTGRES_USER" -d bulle -At -F " " -q -v ON_ERROR_STOP=1 -v subject="$1" -v email="$2"' \
    sh "$subject" "$email" "$BULLE_STORES" | awk '$NF > 0'
}

# Mycelium stores accounts, spaces, tokens and handled notices as JSON files.
mycelium_files() {
  volume="$(docker volume ls --format '{{.Name}}' | grep -- '_mycelium-data$' | head -1)"
  [ -n "$volume" ] || return 0
  docker run --rm --platform linux/amd64 -v "$volume:/data:ro" alpine:3.20 sh -c '
    for f in .users.json .spaces.json tokens.json devices.json .deletions.json; do
      [ -f "/data/$f" ] || continue
      n=$(grep -c -i -F -e "$1" -e "$2" "/data/$f" || true)
      [ "$n" -gt 0 ] && echo "$f $n"
    done; true' sh "$subject" "$email"
}

report() {
  app="$1"
  shift
  rows="$("$@" 2>&1)" || { echo "$app: query failed: $(echo "$rows" | tail -1)"; return; }
  count=0
  if [ -n "$rows" ]; then
    echo "$rows" | while read -r where n; do printf '  %-10s %-45s %s\n' "$app" "$where" "$n"; done
    count=$(echo "$rows" | awk '{ s += $NF } END { print s + 0 }')
  fi
  echo "$app total: $count"
}

printf 'leftovers for subject %s and email %s\n\n' "$subject" "$email"
echo "$APPS" | while read -r app pattern db; do
  [ -n "$app" ] || continue
  container="$(find_container "$pattern")"
  if [ -z "$container" ]; then
    echo "$app: no running container matches $pattern"
    continue
  fi
  report "$app" scan_db "$container" "$db"
  [ "$app" = bulle ] && report "bulle-mx" bulle_stores "$container"
done
report mycelium mycelium_files
