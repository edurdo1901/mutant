#!/bin/bash -i

covignore_path=".covignore"
coverage_report_name="coverage.out"
coverage_destination_path="/tmp"
ci_coverage_destination_path="/tmp/report"

COVERAGE_REPORT=$coverage_destination_path/$coverage_report_name
COVERAGE_DESTINATION=$coverage_destination_path

filter_code_coverage_ignores() {

  if [ -f $covignore_path ]; then
      while read LINE  || [ $LINE ]; do
        for FILE in $LINE; do
          # Use the given line as a valid regex pattern to exclude files and
          # folders from the coverage report.
          cat $COVERAGE_REPORT | grep -v $FILE > "$COVERAGE_DESTINATION/$coverage_report_name.tmp"
          mv "$COVERAGE_DESTINATION/$coverage_report_name.tmp" "$COVERAGE_REPORT"
        done
      done < $covignore_path
  fi
}

echo "=> Running test and generating report"
go test -v ./... -covermode=atomic -coverprofile="$COVERAGE_REPORT" -coverpkg=./... -count=1
exit_code=$?

if [ $exit_code -ne 0 ]; then 
    exit $exit_code;
else 
    filter_code_coverage_ignores
    go tool cover -func=$COVERAGE_REPORT 
fi