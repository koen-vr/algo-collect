#!/bin/bash

date '+keyreg-teal-test start %Y%m%d_%H%M%S'

set -e
set -x
set -o pipefail
export SHELLOPTS

gcmd="goal -d ../../net1/Primary"
ACC1=$(${gcmd} account list | awk '{ print $3 }' | head -n 1)
ACC2=$(${gcmd} account list | awk '{ print $3 }' | tail -1)

APP_ID=1

# max index value 64511
INDEX=64511
UNITNAME="#64511"

# min index value 0
# INDEX=0
# UNITNAME="#00000"

# rand index value 123
# INDEX=123
# UNITNAME="#00123"

# bad combo value 123
# INDEX=123
# UNITNAME="#00321"

# create asset transactions

${gcmd} asset create --creator "$ACC1" \
  --total 1 \
  --unitname $UNITNAME \
  --asseturl "https://path/to/my/asset/details" \
  --decimals 0 \
  --out unsginedtransaction1.txn

# create reserve transactions
${gcmd} app call -f "$ACC1" \
  --app-id "$APP_ID" \
  --app-arg "str:reserve" \
  --app-arg "int:$INDEX" \
  --out unsginedtransaction2.txn

# combine both transactions
cat unsginedtransaction1.txn unsginedtransaction2.txn > combinedtransactions.txn

${gcmd} clerk group \
  -i combinedtransactions.txn \
  -o groupedtransactions.txn

${gcmd} clerk split \
  -i groupedtransactions.txn \
  -o splittransaction

${gcmd} clerk sign \
  -i splittransaction-0 \
  -o splittransaction-0.sig

${gcmd} clerk sign \
  -i splittransaction-1 \
  -o splittransaction-1.sig

cat splittransaction-0.sig splittransaction-1.sig > groupedtransactions.txs

${gcmd} clerk rawsend \
  -f groupedtransactions.txs

# read global state
${gcmd} app read --app-id "$APP_ID" --guess-format --global --from "$ACC1"
