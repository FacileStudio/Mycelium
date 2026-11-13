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
# (case-insensitive), and Synapse also with the person's Matrix id. Every json
# and jsonb column is searched for them as a JSON string value. A text or json
# value that only contains the email is a mention, counted apart from the
# total. Each column has ERASURE_COLUMN_TIMEOUT (default 60s); a column that
# runs out is printed as "timeout" and counted as not scanned. Bulle's
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
timeout="${ERASURE_COLUMN_TIMEOUT:-60s}"
case "$timeout" in
*[!0-9a-z]* | [!0-9]* | "") echo "ERASURE_COLUMN_TIMEOUT must look like 60s or 500ms" >&2; exit 2 ;;
esac

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
vision facile-vision-*-vision-db-1 vision
plume facile-plume-*-plume-db-1 plume
echo compose-*-postgres-1:echo echo
sonde sonde-*-db-1 sonde
capsule capsule-db capsule
casier casier-db casier
journal journal-db journal
antenne facile-antenne-*-postgres-1 antenne
place facile-place-*-db-1 place
toile toile-db toile
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

# The first statement writes a psql script with one counting statement per
# text, uuid, json or jsonb column of every base table, so a slow column costs
# its own statement_timeout and nothing more. Text and uuid columns count rows
# equal to a needle. Json and jsonb columns count rows that hold a needle as a
# JSON string value anywhere. Rows that only contain the email come back apart
# as mentions. The needles travel as psql variables, are quoted with %L inside
# format() and have their LIKE wildcards escaped, so nothing in them is ever
# parsed as SQL or as a pattern. Statements go in on stdin: psql interpolates
# variables there and not in a -c command. The second psql runs the script and
# prints "<table.column> <rows> <mentions>", or the SQLSTATE of a failed
# column, 57014 being the timeout.
SCAN="
WITH needle(v) AS (VALUES (:'subject'), (:'email'), (:'mxid')),
esc AS (
  SELECT v, replace(replace(replace(v, '\\', '\\\\'), '%', '\\%'), '_', '\\_') AS e FROM needle
),
pat AS (
  SELECT array_agg(v) AS exact,
    array_agg('%\"' || e || '\"%') AS json,
    max('%' || e || '%') FILTER (WHERE v = :'email') AS mention
  FROM esc
),
col AS (
  SELECT c.table_schema, c.table_name, c.column_name, c.data_type IN ('json', 'jsonb') AS is_json,
    c.data_type = 'uuid' AS is_uuid
  FROM information_schema.columns c
  JOIN information_schema.tables t USING (table_schema, table_name)
  WHERE t.table_type = 'BASE TABLE'
    AND c.table_schema NOT IN ('pg_catalog', 'information_schema')
    AND c.data_type IN ('text', 'character varying', 'uuid', 'citext', 'json', 'jsonb')
),
pred AS (
  SELECT col.*, CASE WHEN is_json
      THEN format('%I::text ILIKE ANY (%L::text[])', column_name, pat.json)
      ELSE format('%I::text = ANY (%L::text[]) OR lower(%I::text) = lower(%L)',
        column_name, pat.exact, column_name, :'email')
    END AS hit,
    CASE WHEN is_uuid THEN 'false' ELSE format('%I::text ILIKE %L', column_name, pat.mention) END AS mention
  FROM col, pat
)
SELECT format('SET statement_timeout = %L;', :'timeout')
UNION ALL
SELECT line FROM (
  SELECT table_name, column_name, unnest(ARRAY[
    format('SELECT count(*) FILTER (WHERE %s) AS hit, count(*) FILTER (WHERE (%s) AND NOT (%s)) AS mention FROM %I.%I \\gset',
      hit, mention, hit, table_schema, table_name),
    '\\if :ERROR',
    format('\\echo %s.%s :SQLSTATE', table_name, column_name),
    '\\else',
    format('\\echo %s.%s :hit :mention', table_name, column_name),
    '\\endif']) AS line, generate_series(1, 6) AS step
  FROM pred ORDER BY 1, 2, step
) lines"

scan_db() {
  rows="$(docker exec -e PGOPTIONS='-c default_transaction_read_only=on' "$1" sh -c '
    script="$(printf "%s" "$6" | psql -U "$POSTGRES_USER" -d "$1" -At -q -v ON_ERROR_STOP=1 \
      -v subject="$2" -v email="$3" -v mxid="$4" -v timeout="$5")" || exit $?
    printf "%s\n" "$script" | psql -U "$POSTGRES_USER" -d "$1" -At -q 2>/dev/null' \
    sh "$2" "$subject" "$email" "$mxid" "$timeout" "$SCAN")" || return $?
  printf '%s\n' "$rows" | awk '
    NF == 2 { print $1, ($2 == "57014" ? "timeout" : "error " $2) }
    NF == 3 && $2 > 0 { print $1, $2 }
    NF == 3 && $3 > 0 { print $1, "mentions", $3 }'
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
  [ -n "$rows" ] && echo "$rows" | while read -r where what; do printf '  %-10s %-45s %s\n' "$app" "$where" "$what"; done
  echo "$rows" | awk -v app="$app" '
    NF == 2 && $2 ~ /^[0-9]+$/ { total += $2 }
    $2 == "mentions" { mentions += $3 }
    $2 == "timeout" || $2 == "error" { failed++ }
    END {
      line = app " total: " total + 0
      if (mentions) line = line ", mentions: " mentions
      if (failed) line = line ", columns not scanned: " failed
      print line
    }'
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
