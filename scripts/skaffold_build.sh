#!/usr/bin/env bash
set -e

cmd=$1
imageName=$2

go vet ./...

BUILDARGS=$(
	cat <<-EOM
		{
		  "operator" : {
		    "name" : "Operator",
		    "exec" : "./cmd/manager",
		    "bin_out" : "redhat-marketplace-operator"
		  },
		  "reporter" : {
		    "name" : "Reporter",
		    "exec" : "./cmd/reporter",
		    "bin_out" : "redhat-marketplace-reporter"
		  },
		  "metric-state" : {
		    "name" : "Metric State",
		    "exec" : "./cmd/metrics",
		    "bin_out" : "redhat-marketplace-metric-state"
		  },
		  "authcheck" : {
		    "name" : "AuthCheck",
		    "exec" : "./cmd/authvalid",
		    "bin_out" : "redhat-marketplace-authcheck"
		  }
		}
	EOM
)

name=$(echo $BUILDARGS | jq -r --arg i $imageName '.[$i].name')
exec=$(echo $BUILDARGS | jq -r --arg i $imageName '.[$i].exec')
bin=$(echo $BUILDARGS | jq -r --arg i $imageName '.[$i].bin_out')

if [[ "$cmd" == "dependencies" ]]; then
	OUTPUT=$(go list -f '{{ join .Deps "\n" }}' $exec 2>/dev/null | grep github.com/redhat-marketplace/redhat-marketplace-operator | sed -e 's/github.com\/redhat-marketplace\/redhat-marketplace-operator\///' | xargs find . | uniq | xargs jq -nc '$ARGS.positional' --args)
	execFiles=$(find . $exec | xargs jq -nc '$ARGS.positional' --args)
	OUTPUT=$(echo $OUTPUT | jq --argjson e $execFiles '. + $e')
	OUTPUT=$(echo $OUTPUT | jq '. + ["build/Dockerfile", "build/bin/entrypoint", "build/bin/user_setup"]')
	echo $OUTPUT
	exit 0
fi

QUAY_EXPIRATION=${QUAY_EXPIRATION:-never}
VERSION=${VERSION:-latest}

if [ "${DOCKER_EXEC}" == "" ]; then
	DOCKER_EXEC=$(command -v docker)
fi

${DOCKER_EXEC} buildx &>/dev/null
if [ $? -eq 0 ]; then
	DOCKER="${DOCKER_EXEC} buildx"
	if $PUSH_IMAGE; then
		ARGS="--push"
	else
		ARGS="--load"
	fi
else
	DOCKER="${DOCKER}"
fi

unset PUSH

echo $args $PUSH

export DOCKER_BUILDKIT=1

${DOCKER} build -f ./build/Dockerfile $ARGS --tag $IMAGE --build-arg name="$name" --build-arg exec=$exec --build-arg bin=$bin --build-arg bin_out=$bin --build-arg BUILDKIT_INLINE_CACHE=1 --build-arg app_version=\"$VERSION\" --build-arg quay_expiration=\"$QUAY_EXPIRATION\" .

if [ $PUSH_IMAGE ] && [ "$ARGS" != "--push" ]; then
	${DOCKER} push $IMAGE
fi
