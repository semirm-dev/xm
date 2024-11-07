COVERAGE_FILE = coveragefile

test-cover:
	go test -v ./companies ./gateway -coverprofile=${COVERAGE_FILE}
	go tool cover -html=${COVERAGE_FILE} && go tool cover -func ${COVERAGE_FILE} && unlink ${COVERAGE_FILE}
