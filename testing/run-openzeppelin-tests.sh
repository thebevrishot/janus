#!/bin/sh
cleanupDocker () {
  docker-compose -f docker-compose-openzeppelin.yml -p ci kill
  docker-compose -f docker-compose-openzeppelin.yml -p ci rm -f
}
trap 'cleanupDocker ; echo "Tests Failed For Unexpected Reasons"' HUP INT QUIT PIPE TERM
docker-compose -p ci -f docker-compose-openzeppelin.yml build && docker-compose -p ci -f docker-compose-openzeppelin.yml up -d
if [ $? -ne 0 ] ; then
  echo "Docker Compose Failed"
  exit -1
fi
docker logs ci_openzeppelin_1 -f&
EXIT_CODE=`docker wait ci_openzeppelin_1`
if [ -z ${EXIT_CODE+z} ] || [ -z ${EXIT_CODE} ] || ([ "0" != "$EXIT_CODE" ] && [ "" != "$EXIT_CODE" ]) ; then
  docker logs qtum_seeded_testchain
  docker logs ci_janus_1
  docker logs ci_openzeppelin_1
  echo "Tests Failed - Exit Code: $EXIT_CODE"
else
  echo "Tests Passed"
fi
cleanupDocker
exit $EXIT_CODE
