payload=$(mktemp $TMPDIR/script-request.XXXXXX)

cat > $payload <&0

uri=$(jq -r '.source.locker_uri // ""' < $payload)
username="$(jq -r '.source.username // ""' < $payload)"
password=$(jq -r '.source.password // ""' < $payload)
ca_cert=$(jq -r '.source.ca_cert // ""' < $payload)
skip_ssl_validation=$(jq -r '.source.skip_ssl_validation // ""' < $payload)

lock="$(jq -r '.source.lock_name // ""' < $payload)"
key="$(jq -r '.params.key // ""' < $payload)"
locked_by="$(jq -r '.params.locked_by // ""' < $payload)"
operation="$(jq -r '.params.lock_op // ""' < $payload)"

if [[ -z "$uri" ]]; then
  echo >&2 "invalid payload (missing locker_uri):"
  cat $payload >&2
  exit 1
fi

if [[ -z "${lock}" ]]; then
  echo >&2 "invalid payload (missing lock_name)"
  cat $payload >&2
  exit 1
fi

if [[ -n ${ca_cert} ]]; then
  ca_cert_file=$(mktemp)
  cat <<EOF > $ca_cert_file
$ca_cert
EOF
  ca_cert_flag="--cacert ${ca_cert_file}"
fi


if [[ $skip_ssl_validation == "true" ]]; then
	skip_ssl_flag="-k"
fi
if [[ -n $username && -n $password ]]; then
	auth_flag="-u ${username}:${password}"
fi

calc_reference() {
  http_req /locks | sha1sum | cut -d " " -f1
}

http_req() {
	resource=$1
	shift
	curl -fsS $auth_flag $ca_cert_flag $skip_ssl_flag "${uri}${resource}" $@ 2>&1
}
