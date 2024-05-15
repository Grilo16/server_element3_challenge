set -e

host="$1"
port="$2"
shift 2
cmd="$@"

until nc -z "$host" "$port"; do
  >&2 echo "SQL Server is unavailable - sleeping"
  sleep 1
done

>&2 echo "SQL Server is up - executing command"
if [ -z "$cmd" ]; then
  echo "No command specified to execute. Exiting."
  exit 1
else
  exec $cmd
fi