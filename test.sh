# 10274 - manual testing
#bazel run cli:cli -- user sets delete -n 10274-1
#bazel run cli:cli -- user sets delete -n 10274
#bazel run cli:cli -- user sets get
#bazel run cli:cli -- user sets set -n 10274
#bazel run cli:cli -- user sets get

bazel test --action_env=REBRICKABLE_USERNAME=$REBRICKABLE_USERNAME \
 --action_env=REBRICKABLE_PASSWORD=$REBRICKABLE_PASSWORD \
 --action_env=REBRICKABLE_API_KEY=$REBRICKABLE_API_KEY \
 //...