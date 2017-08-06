BINARY=microstripe
BINFULLPATH=${GOPATH}/bin/${BINARY}

BUILD_VERSION=$(shell date '+%Y%m%d-%H%M')

default:
	build

build:
	go build -o ${BINFULLPATH}

buildprod:
	CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/${BINARY}

clean:
	if [ -f ${BINFULLPATH} ] ; then rm ${BINFULLPATH} ; fi

rundebug:
	${BINFULLPATH}

run:
	${BINFULLPATH}
