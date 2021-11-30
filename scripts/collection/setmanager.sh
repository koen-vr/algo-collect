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

# opt in transactions
${gcmd} app optin --app-id "$APP_ID" \
  --from "$ACC2" \
  --out unsginedtransaction1.txn

# set manager transactions
${gcmd} app call -f "$ACC1" \
  --app-id "$APP_ID" \
  --app-arg "str:manager" \
  --app-account "$ACC2" \
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
${gcmd} app read --app-id "$APP_ID" --guess-format --local --from "$ACC1"
${gcmd} app read --app-id "$APP_ID" --guess-format --local --from "$ACC2"
